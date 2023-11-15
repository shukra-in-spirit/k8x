package models

import "time"

type PrometheusDataSetResponse struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float32   `json:"cpu"`
	Memory    float32   `json:"memory"`
}
