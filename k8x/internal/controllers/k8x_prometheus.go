package controllers

import (
	"context"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/constants"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

// based on https://prometheus.io/docs/prometheus/latest/querying/api/

type PrometheusFunctions interface {
	GetPrometheusData(ctx context.Context, promQuery string) (*models.PrometheusDataSetResponse, error)
	GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time, steps string) (*models.PrometheusDataSetResponse, error)
}

type PrometheusInstance struct {
	url string
}

func NewPrometheusInstance(promUrl string) *PrometheusInstance {
	return &PrometheusInstance{
		url: promUrl,
	}
}

func BuildPromQueryForCPU(namespace_name string, rate string, container_name string) string {
	pq_cpu := "sum(rate(container_cpu_usage_seconds_total{container=\"" + container_name + "\",namespace=\"" + namespace_name + "\"}[" + rate + "]))"
	return pq_cpu
}

func BuildPromQueryForMemory(namespace_name string, rate string, container_name string) string {
	pq_memory := "sum(rate(container_memory_usage_bytes{container=\"" + container_name + "\",namespace=\"" + namespace_name + "\"}[" + rate + "]))"
	return pq_memory
}

func (prom *PrometheusInstance) GetPrometheusData(ctx context.Context, promQuery string) (*models.PrometheusDataSetResponse, error) {
	endTime := time.Now()
	startTimeInt := endTime.Unix() - constants.DefaultPrometheusTimeRange
	startTime := time.Unix(startTimeInt, 0)

	responseDataFrame, err := prom.GetPrometheusDataWithinRange(ctx, promQuery, startTime, endTime)
	if err != nil {
		return &models.PrometheusDataSetResponse{}, nil
	}
	return responseDataFrame, nil
}

func (prom *PrometheusInstance) GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time) (*models.PrometheusDataSetResponse, error) {
	responseDataFrame := models.PrometheusDataSetResponse{}
	return &responseDataFrame, nil
}
