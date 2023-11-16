package models

import "time"

type PrometheusDataSetResponseItem struct {
	Timestamp time.Time `json:"timestamp"`
	Metric    float32   `json:"metric"`
}

type PrometheusDataSetResponse struct {
	PromItemList []PrometheusDataSetResponseItem `json:"items"`
}

type LambdaRequest struct {
	ServiceID string `json:"service_id"`
	Params TuningParams
	History PrometheusDataSetResponse
}

type TuningParams struct {
	Epochs string `json:"epochs"`
	HiddenLayers string `json:"hidden_layers"`
}

type LambdaRespBody struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	// Replicas string `json:"replicas"`
}

type LambdaResponse struct {
	StatusCode int            `json:"statusCode"`
	Body       LambdaRespBody `json:"body"`
}
