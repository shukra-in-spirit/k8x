package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/constants"
	"github.com/shukra-in-spirit/k8x/internal/controllers"
	"github.com/shukra-in-spirit/k8x/internal/helpers"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

func (listener *K8ManagerAPI) verticalScaler(ctx context.Context, id string) {
	ticker := time.NewTicker(time.Duration(168) * time.Hour)
	defer ticker.Stop()
	for range ticker.C {

		avg_cpu, avg_mem, err := listener.commonCreateFunc(ctx, id, "c")
		if err != nil {
			log.Printf("failed fetching data and retraining model from create lambda: %v\n", err)

			continue
		}

		if avg_cpu != "" && avg_mem != "" {
			deployment, namespace := helpers.DecomposeServiceID(id)
			final_cpu, _ := strconv.ParseFloat(avg_cpu, 32)
			final_mem, _ := strconv.ParseFloat(avg_mem, 32)

			// scale the values for cpu.
			err = listener.kubeClient.SetLimitValue(ctx, deployment, namespace, float32(final_cpu*3), float32(final_mem*3))
			if err != nil {
				log.Printf("failed setting limit value: %v\n", err)

				continue
			}

			err = listener.kubeClient.SetRequestValue(ctx, deployment, namespace, float32(final_cpu), float32(final_mem))
			if err != nil {
				log.Printf("failed setting request value: %v\n", err)

				continue
			}

			log.Println("successfully completed vertical scaling.")
		}
	}
}

func (listener *K8ManagerAPI) commonCreateFunc(ctx context.Context, id, funcName string) (string, string, error) {
	currTime := time.Now()
	startTime := currTime.AddDate(0, 0, -14)

	depl, ns := helpers.DecomposeServiceID(id)

	container, err := listener.kubeClient.GetContainerNameFromDeployment(ctx, depl, ns)
	if err != nil {
		return "", "", fmt.Errorf("failed fetching container name for service %s: %v", id, err)
	}

	var promCPU, promMem *models.PrometheusDataSetResponse
	var err1, err2 error

	// Fetch 2 weeks data from prom.
	if helpers.GetEnvOrDefault("PROM_MODE", constants.Local) == constants.Local {
		promCPU, err1 = controllers.GetCSVData(ctx, id+"-cpu-train.csv")
		promMem, err2 = controllers.GetCSVData(ctx, id+"-mem-train.csv")
	} else {
		promCPU, err1 = listener.promClient.GetPrometheusDataWithinRange(ctx, controllers.BuildPromQueryForCPU(ns, "2m", container), startTime, currTime, "20m", "cpu")
		promMem, err2 = listener.promClient.GetPrometheusDataWithinRange(ctx, controllers.BuildPromQueryForMemory(ns, "2m", container), startTime, currTime, "20m", "memory")
	}

	// promCPU, err = listener.promClient.GetPrometheusDataWithinRange(ctx, controllers.BuildPromQueryForCPU(ns, "2m", container), startTime, currTime, "20m")
	if err1 != nil {
		return "", "", fmt.Errorf("failed fetching data from prometheus: %v", err1)
	}

	// promMem, err = listener.promClient.GetPrometheusDataWithinRange(ctx, controllers.BuildPromQueryForMemory(ns, "2m", container), startTime, currTime, "20m")
	if err2 != nil {
		return "", "", fmt.Errorf("failed fetching data from prometheus: %v", err2)
	}

	promData := helpers.ProcessPromData(id, promCPU.PromItemList, promMem.PromItemList)

	// Push to DB.
	// err = listener.dbClient.AddData(data)
	err = listener.dbClient.AddDataBatch(promData, id)
	if err != nil {
		return "", "", fmt.Errorf("batch DB write failed: %v", err)
	}

	input := models.LambdaRequest{ServiceID: id, Params: models.TuningParams{}}

	// build the request.
	payload, err := json.Marshal(input)
	if err != nil {
		return "", "", fmt.Errorf("failed marshalling input to lambda function: %v", err)
	}

	var output *models.LambdaRespBody

	// Call create lambda.
	if helpers.GetEnvOrDefault("LAMBDA_MODE", constants.Local) == constants.Local {
		output, err = helpers.TriggerLocalCreateLambdaWithEvent(payload, funcName)
	} else {
		output, err = listener.lambdaClient.TriggerLambdaWithEvent(payload, funcName)
	}
	// output, err := listener.lambdaClient.TriggerLambdaWithEvent(payload, funcName)
	if err != nil {
		return "", "", fmt.Errorf("lambda trigger failed: %v", err)
	}

	return output.CPU, output.Memory, nil
}
