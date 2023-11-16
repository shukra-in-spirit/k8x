package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shukra-in-spirit/k8x/internal/helpers"
)

func (listener *K8ManagerAPI) GetInfo(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("got the service id %s", serviceID)

	ctx := c.Request.Context()
	depl, ns := helpers.DecomposeServiceID(serviceID)

	container, err := listener.kubeClient.GetContainerNameFromDeployment(ctx, depl, ns)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed fetching container name for deployment %s: %s", depl, err.Error()),
		})
	}

	req_cpu, req_mem, err := listener.kubeClient.GetRequestValue(ctx, depl, ns)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed fetching request values for deployment %s: %s", depl, err.Error()),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"container":   container,
		"request_cpu": req_cpu,
		"request_mem": req_mem,
	})

}

func (listener *K8ManagerAPI) ScaleOnPredict(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("got the service id %s", serviceID)

	// ctx := c.Request.Context()

	// ar req models.LambdaRespBody
	// if err := ginCtx.BindJSON(&req); err != nil {
	// 	hpeError.SetBadRequest(ginCtx, hpeError.ErrorResponse{
	// 		Message:            err.Error(),
	// 		RecommendedActions: []string{constants.CheckRequestBody},
	// 		ErrorCode:          constants.BadRequest,
	// 	})

	// 	return
	// }

	c.JSON(http.StatusAccepted, gin.H{
		"message": "autoscaling pods",
	})

	return
}
