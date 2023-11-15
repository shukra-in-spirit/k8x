package models

import "time"

type PrometheusDataSetResponseItem struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float32   `json:"cpu"`
	Memory    float32   `json:"memory"`
}

type PrometheusDataSetResponse struct {
	PromItemList []PrometheusDataSetResponseItem `json:"items"`
}
