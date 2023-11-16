package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
