package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shukra-in-spirit/k8x/internal/clients"
	"github.com/shukra-in-spirit/k8x/internal/constants"
	"github.com/shukra-in-spirit/k8x/internal/controllers"
	"github.com/shukra-in-spirit/k8x/internal/helpers"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

type K8Manager interface {
	AddServiceTok8x(ginCtx *gin.Context)
	StartPredictionService(ginCtx *gin.Context)
	ScaleOnPredict(ginCtx *gin.Context)
	GetInfo(ginCtx *gin.Context)
}

type K8ManagerAPI struct {
	dbClient     clients.PromStorer
	promClient   controllers.PrometheusFunctions
	lambdaClient clients.NewLambdaInterface
	kubeClient   controllers.ScalingFunctions
}

func NewK8Manager(
	dbClient clients.PromStorer,
	promClient controllers.PrometheusFunctions,
	lambdaClient clients.NewLambdaInterface,
	kubeClient controllers.ScalingFunctions,
) *K8ManagerAPI {
	return &K8ManagerAPI{
		dbClient:     dbClient,
		promClient:   promClient,
		lambdaClient: lambdaClient,
		kubeClient:   kubeClient,
	}
}

func (listener *K8ManagerAPI) AddServiceTok8x(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("got the service id %s", serviceID)

	ctx := c.Request.Context()

	_, _, err := listener.commonCreateFunc(ctx, serviceID, "c")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "done training model",
	})

	return

	// currTime := time.Now()
	// startTime := currTime.AddDate(0, 0, -14)

	// // Fetch 2 weeks data from prom.
	// promData, err := listener.promClient.GetPrometheusDataWithinRange(ctx, "", startTime, currTime, "")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "failed fetching data from prometheus: " + err.Error(),
	// 	})

	// 	return
	// }

	// // Push to DB.
	// // err = listener.dbClient.AddData(data)
	// err = listener.dbClient.AddDataBatch(&promData.PromItemList, serviceID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "batch DB write failed: " + err.Error(),
	// 	})

	// 	return
	// }

	// // Call create lambda and drop the response.
	// // payload, err := json.Marshal(data)
	// // if err != nil {
	// // 	return fmt.Errorf("failed marshalling input to lambda function: %v", err)
	// // }

	// _, err = listener.lambdaClient.TriggerCreateLambdaWithEvent([]byte("{}"), "change-lambda-func-name")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "lambda trigger failed: " + err.Error(),
	// 	})

	// 	return
	// }
}

func (listener *K8ManagerAPI) StartPredictionService(c *gin.Context) {
	serviceID := c.Param("service_id")
	message := fmt.Sprintf("pocessing the service id %s", serviceID)

	// Runs every 20 mins
	go listener.triggerPredict(c.Request.Context(), serviceID)

	// Runs every week
	go listener.verticalScaler(c.Request.Context(), serviceID)

	c.JSON(http.StatusAccepted, gin.H{
		"message": message,
	})
}

func (listener *K8ManagerAPI) triggerPredict(ctx context.Context, id string) {
	ticker := time.NewTicker(time.Duration(20) * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		listener.performOp(ctx, id)
	}
}

func (listener *K8ManagerAPI) performOp(ctx context.Context, id string) {
	deployment, namespace := helpers.DecomposeServiceID(id)

	container, err := listener.kubeClient.GetContainerNameFromDeployment(ctx, deployment, namespace)
	if err != nil {
		log.Printf("failed fetching container name for service %s: %v\n", id, err)

		return
	}

	var promCPU, promMem *models.PrometheusDataSetResponse
	var err1, err2 error

	currTime := time.Now()
	startTime := currTime.Add(-3 * time.Hour)

	// Fetch 3 hours data from prom.
	if helpers.GetEnvOrDefault("PROM_MODE", constants.PrometheusMode) == "local" {
		promCPU, err1 = controllers.GetCSVData(ctx, id+"-cpu-pred.csv")
		promMem, err2 = controllers.GetCSVData(ctx, id+"-mem-pred.csv")
	} else {
		promCPU, err1 = listener.promClient.GetPrometheusDataWithinRange(ctx, controllers.BuildPromQueryForCPU(namespace, "2m", container), startTime, currTime, constants.StepsMinutesInterval*time.Minute, "cpu")
		promMem, err2 = listener.promClient.GetPrometheusDataWithinRange(ctx, controllers.BuildPromQueryForMemory(namespace, "2m", container), startTime, currTime, constants.StepsMinutesInterval*time.Minute, "memory")
	}

	// promCPU, err := listener.promClient.GetPrometheusData(ctx, controllers.BuildPromQueryForCPU(namespace, "2m", container))
	if err1 != nil {
		log.Printf("failed fetching cpu data from prometheus: %v\n", err1)

		return
	}

	// promMem, err := listener.promClient.GetPrometheusData(ctx, controllers.BuildPromQueryForMemory(namespace, "2m", container))
	if err2 != nil {
		log.Printf("failed fetching memory data from prometheus: %v\n", err2)

		return
	}

	promData := helpers.PrepareHistoryData(promCPU.PromItemList, promMem.PromItemList)

	input := models.LambdaRequest{ServiceID: id, Params: models.TuningParams{}, History: promData}

	// Get the current resource utilization -> request cpu
	request_cpu, _, err := listener.kubeClient.GetRequestValue(ctx, deployment, namespace)
	if err != nil {
		log.Printf("failed getting current request values: %v\n", err)

		return
	}

	// build the request.
	payload, err := json.Marshal(input)
	if err != nil {
		log.Printf("failed marshalling input to lambda function: %v\n", err)

		return
	}

	var output *models.LambdaRespBody
	// Call predict lambda
	if helpers.GetEnvOrDefault("LAMBDA_MODE", constants.Local) == constants.Local {
		output, err = helpers.TriggerLocalCreateLambdaWithEvent(id, "p")
	} else {
		output, err = listener.lambdaClient.TriggerLambdaWithEvent(payload, "p")
	}
	if err != nil {
		log.Printf("lambda trigger failed: %v\n", err)

		return
	}

	if output.CPU != "" {
		predicted_cpu, _ := strconv.Atoi(output.CPU)
		int_request_cpu, _ := strconv.Atoi(request_cpu)

		// replica = predicted cpu/request cpu
		err = listener.kubeClient.SetReplicaValue(ctx, deployment, namespace, int32(predicted_cpu/int_request_cpu))
		if err != nil {
			log.Printf("failed setting limit value: %v\n", err)

			return
		}
	}

	log.Println("finished calling predict lambda and horizontally scaled the pods")

	return
}
