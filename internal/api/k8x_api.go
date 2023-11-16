package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shukra-in-spirit/k8x/internal/clients"
	"github.com/shukra-in-spirit/k8x/internal/controllers"
)

type K8Manager interface {
	AddServiceTok8x(ginCtx *gin.Context)
}

type K8ManagerAPI struct {
	dbClient     clients.PromStorer
	promClient   controllers.PrometheusFunctions
	lambdaClient clients.NewLambdaInterface
}

func NewK8Manager(
	dbClient clients.PromStorer,
	promClient controllers.PrometheusFunctions,
	lambdaClient clients.NewLambdaInterface,
) *K8ManagerAPI {
	return &K8ManagerAPI{
		dbClient:     dbClient,
		promClient:   promClient,
		lambdaClient: lambdaClient,
	}
}

func (listener *K8ManagerAPI) AddServiceTok8x(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("got the service id %s", serviceID)

	ctx := c.Request.Context()

	currTime := time.Now()
	startTime := currTime.AddDate(0, 0, -14)

	// Fetch 2 weeks data from prom.
	promData, err := listener.promClient.GetPrometheusDataWithinRange(ctx, "", startTime, currTime, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed fetching data from prometheus: " + err.Error(),
		})

		return
	}

	// Push to DB.
	// err = listener.dbClient.AddData(data)
	err = listener.dbClient.AddDataBatch(&promData.PromItemList, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "batch DB write failed: " + err.Error(),
		})

		return
	}

	// Call create lambda and drop the response.
	// payload, err := json.Marshal(data)
	// if err != nil {
	// 	return fmt.Errorf("failed marshalling input to lambda function: %v", err)
	// }

	_, err = listener.lambdaClient.TriggerCreateLambdaWithEvent([]byte("{}"), "c")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "lambda trigger failed: " + err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "done",
	})

	return
}
