package models

import "time"

type PrometheusDataSetResponseItem struct {
	Timestamp time.Time `json:"timestamp"`
	Metric    float32   `json:"metric"`
}

type PrometheusDataSetResponse struct {
	PromItemList []PrometheusDataSetResponseItem `json:"items"`
}
