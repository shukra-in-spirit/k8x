package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gin-gonic/gin"
	"github.com/shukra-in-spirit/k8x/internal/api"
	"github.com/shukra-in-spirit/k8x/internal/clients"
	"github.com/shukra-in-spirit/k8x/internal/config"
	"github.com/shukra-in-spirit/k8x/internal/controllers"
)

func main() {
	// create the default gin router
	router := gin.Default()

	conf, err := config.NewServiceConfig()
	if err != nil {
		panic(fmt.Sprintf("failed retrieving svc config: %v", err))
	}

	session, err := clients.CreateSession(conf.AWSConfig)
	if err != nil {
		panic(err)
	}

	dbHandler := clients.NewPromStore(dynamodb.New(session))
	promHandler := controllers.NewPrometheusInstance(conf.PromUrl)
	lambdaHandler := clients.NewLamdaClient(lambda.New(session))

	kubeManager := api.NewK8Manager(dbHandler, promHandler, lambdaHandler)

	// create the endpoints
	router.POST("/:service_id", kubeManager.AddServiceTok8x)
	router.POST("/:service_id/start", startPredictionOfService)

	// run it on port 8585
	router.Run(":8585")
}

func startPredictionOfService(c *gin.Context) {
	serviceID := c.Param("service_id")
	message := fmt.Sprintf("got the service id %s", serviceID)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
