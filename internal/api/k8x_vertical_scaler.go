package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/models"
)

func (listener *K8ManagerAPI) verticalScaler(ctx context.Context, id string) {
	ticker := time.NewTicker(time.Duration(168) * time.Hour)
	defer ticker.Stop()
	for range ticker.C {

		avg_cpu, avg_mem, err := listener.commonCreateFunc(ctx, id, "modified_create_lambda_function")
		if err != nil {
			log.Printf("failed fetching data and retraining model from create lambda: %v\n", err)

			continue
		}

		if avg_cpu != "" && avg_mem != "" {
			final_cpu, _ := strconv.ParseFloat(avg_cpu, 32)
			final_mem, _ := strconv.ParseFloat(avg_mem, 32)

			// scale the values
			err = listener.kubeClient.SetLimitValue(ctx, "", "", float32(final_cpu*3), float32(final_mem*3))
			if err != nil {
				log.Printf("failed setting limit value: %v\n", err)

				continue
			}

			err = listener.kubeClient.SetRequestValue(ctx, "", "", float32(final_cpu), float32(final_mem))
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

	// Fetch 2 weeks data from prom.
	promData, err := listener.promClient.GetPrometheusDataWithinRange(ctx, "", startTime, currTime, "")
	if err != nil {
		return "", "", fmt.Errorf("failed fetching data from prometheus: %v", err)
	}

	// Push to DB.
	// err = listener.dbClient.AddData(data)
	err = listener.dbClient.AddDataBatch(&promData.PromItemList, id)
	if err != nil {
		return "", "", fmt.Errorf("batch DB write failed: %v", err)
	}

	input := models.LambdaRequest{ServiceID: id, Params: models.TuningParams{}}

	// build the request.
	payload, err := json.Marshal(input)
	if err != nil {
		return "", "", fmt.Errorf("failed marshalling input to lambda function: %v", err)
	}

	// Call create lambda.
	output, err := listener.lambdaClient.TriggerCreateLambdaWithEvent(payload, funcName)
	if err != nil {
		return "", "", fmt.Errorf("lambda trigger failed: %v", err)
	}

	return output.CPU, output.Memory, nil
}
