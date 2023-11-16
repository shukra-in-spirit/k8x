package main

import (
	"fmt"

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
	kubeHandler := controllers.NewKubeClient()

	kubeManager := api.NewK8Manager(dbHandler, promHandler, lambdaHandler, kubeHandler)

	// create the endpoints.
	router.POST("/:service_id", kubeManager.AddServiceTok8x)

	router.POST("/:service_id/start", kubeManager.StartPredictionService)

	// unused..
	router.POST("/:service_id/hscale", kubeManager.ScaleOnPredict)

	router.GET("/:service_id/stats", kubeManager.GetInfo)

	// run it on port 8585
	router.Run(":8585")
}
