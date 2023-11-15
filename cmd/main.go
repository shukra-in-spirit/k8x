package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shukra-in-spirit/k8x/internal/api"
	"github.com/shukra-in-spirit/k8x/internal/config"
	"github.com/shukra-in-spirit/k8x/internal/controllers"
	"github.com/shukra-in-spirit/k8x/internal/database"
)

func main() {
	// create the default gin router
	router := gin.Default()

	conf, err := config.NewServiceConfig()
	if err != nil {
		panic(fmt.Sprintf("failed retrieving svc config: %v", err))
	}

	dbHandler := database.NewPromStore(conf.AWSConfig)
	promHandler := controllers.NewPrometheusInstance(conf.PromUrl)

	kubeManager := api.NewK8Manager(dbHandler, promHandler)

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
