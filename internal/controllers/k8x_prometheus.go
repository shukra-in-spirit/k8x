package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	promModel "github.com/prometheus/common/model"
	"github.com/shukra-in-spirit/k8x/internal/constants"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

// based on https://prometheus.io/docs/prometheus/latest/querying/api/

type PrometheusFunctions interface {
	GetPrometheusData(ctx context.Context, promQuery string, queryType string) (*models.PrometheusDataSetResponse, error)
	GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time, steps time.Duration, queryType string) (*models.PrometheusDataSetResponse, error)
}

type PrometheusInstance struct {
	PromClientAPI v1.API
}

func NewPrometheusInstance(promUrl string) *PrometheusInstance {
	// Create a new Prometheus API client
	client, err := api.NewClient(api.Config{
		Address: promUrl,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating client:", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)

	return &PrometheusInstance{
		PromClientAPI: v1api,
	}
}

func BuildPromQueryForCPU(namespace_name string, rate string, container_name string) string {
	pq_cpu := "rate(eagle_pod_container_resource_usage_cpu_cores{exported_container=\"" + container_name + "\",exported_namespace=\"" + namespace_name + "\"}[" + rate + "])"
	return pq_cpu
}

func BuildPromQueryForMemory(namespace_name string, rate string, container_name string) string {
	pq_memory := "rate(eagle_pod_container_resource_usage_memory_bytes{exported_container=\"" + container_name + "\",exported_namespace=\"" + namespace_name + "\"}[" + rate + "])"
	return pq_memory
}

func (prom *PrometheusInstance) GetPrometheusData(ctx context.Context, promQuery string, queryType string) (*models.PrometheusDataSetResponse, error) {
	endTime := time.Now()
	startTimeInt := endTime.Unix() - constants.DefaultPrometheusTimeRange
	startTime := time.Unix(startTimeInt, 0)

	responseDataFrame, err := prom.GetPrometheusDataWithinRange(ctx, promQuery, startTime, endTime, constants.StepsMinutesInterval*time.Minute, queryType)
	if err != nil {
		fmt.Errorf("%v", err)
		return &models.PrometheusDataSetResponse{}, nil
	}
	return responseDataFrame, nil
}

func (prom *PrometheusInstance) GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time, steps time.Duration, queryType string) (*models.PrometheusDataSetResponse, error) {

	queryRangeJSON := v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  steps,
	}

	result, warnings, err := prom.PromClientAPI.QueryRange(ctx, promQuery, queryRangeJSON, v1.WithTimeout(constants.PrometheusRequestTimeOut*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	list := MarshalPromResult(result, queryType)

	return list, nil
}

func MarshalPromResult(PromResult promModel.Value, queryType string) *models.PrometheusDataSetResponse {
	// Process the result
	dataSetResponse := models.PrometheusDataSetResponse{
		PromDataType: queryType,
		PromItemList: make([]models.PrometheusDataSetResponseItem, 0),
	}

	// Assuming result is of type model.Value
	switch data := PromResult.(type) {
	case model.Matrix: // Check if the result is a Matrix (range vector)
		for _, stream := range data {
			for _, value := range stream.Values {
				dataSetResponse.PromItemList = append(dataSetResponse.PromItemList, models.PrometheusDataSetResponseItem{
					Metric: float32(value.Value),
				})
			}
		}
	// Handle other types (model.Vector, model.Scalar, etc.) as needed
	default:
		fmt.Println("Unknown result format")
	}

	// Print the parsed data
	//fmt.Printf("Parsed Data: %+v\n", dataSetResponse)
	return &dataSetResponse
}
