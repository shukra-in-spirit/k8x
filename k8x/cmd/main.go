package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// create the default gin router
	router := gin.Default()

	// create the endpoints
	router.POST("/:service_id", addServiceTok8x)
	router.POST("/:service_id/start", startPredictionOfService)

	// run it on port 8585
	router.Run(":8585")
}

func addServiceTok8x(c *gin.Context) {
	serviceID := c.Param("service_id")
	message := fmt.Sprintf("got the service id %s", serviceID)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func startPredictionOfService(c *gin.Context) {
	serviceID := c.Param("service_id")
	message := fmt.Sprintf("got the service id %s", serviceID)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
