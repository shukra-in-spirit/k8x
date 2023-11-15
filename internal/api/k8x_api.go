package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shukra-in-spirit/k8x/internal/controllers"
	"github.com/shukra-in-spirit/k8x/internal/database"
)

type K8Manager interface {
	AddServiceTok8x(ginCtx *gin.Context)
}

type K8ManagerAPI struct {
	dbClient   database.PromStorer
	promClient controllers.PrometheusFunctions
}

func NewK8Manager(
	dbClient database.PromStorer,
	promClient controllers.PrometheusFunctions,
) *K8ManagerAPI {
	return &K8ManagerAPI{
		dbClient:   dbClient,
		promClient: promClient,
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
	}

	data := &database.PromData{
		ServiceID: serviceID,
		Timestamp: promData.Timestamp,
		CPU:       promData.CPU,
		Memory:    promData.Memory,
	}

	// Push to DB.
	err = listener.dbClient.AddData(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "DB write failed: " + err.Error(),
		})
	}

	// Call create lambda.

	c.JSON(http.StatusOK, gin.H{
		"message": "done",
	})
}
