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
	Params    TuningParams
	History   []*History `json:"history"`
}

type History struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float32   `json:"cpu"`
	Memory    float32   `json:"memory"`
}

type TuningParams struct {
	Epochs       string `json:"epochs"`
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

type PromData struct {
	ServiceID string    `json:"service_id"`
	Timestamp time.Time `json:"timestamp"`
	CPU       float32   `json:"cpu"`
	Memory    float32   `json:"memory"`
}
