package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/constants"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

// based on https://prometheus.io/docs/prometheus/latest/querying/api/

type PrometheusFunctions interface {
	GetPrometheusData(ctx context.Context, promQuery string, queryType string) (*models.PrometheusDataSetResponse, error)
	GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time, steps string, queryType string) (*models.PrometheusDataSetResponse, error)
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

	responseDataFrame, err := prom.GetPrometheusDataWithinRange(ctx, promQuery, startTime, endTime, "", queryType)
	if err != nil {
		fmt.Errorf("%v", err)
		return &models.PrometheusDataSetResponse{}, nil
	}
	return responseDataFrame, nil
}

func (prom *PrometheusInstance) GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time, steps string, queryType string) (*models.PrometheusDataSetResponse, error) {
	// Create a JSON payload with the query_range parameters
	queryRangeJSON := map[string]interface{}{
		"query":  promQuery,
		"start":  startTime,
		"end":    endTime, // Set end time to the current time
		"step":   steps,   // Step size for the range vector in seconds (e.g., 20s)
		"format": "json",
	}
	queryRangeBytes, err := json.Marshal(queryRangeJSON)
	if err != nil {
		fmt.Errorf("%v", err)
		return nil, err
	}

	// Create a POST request to the Prometheus query_range API
	req, err := http.NewRequest("POST", prom.url, bytes.NewBuffer(queryRangeBytes))
	if err != nil {
		fmt.Errorf("%v", err)
		return nil, err
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("%v", err)
		return nil, err
	}

	// Create an instance of prometheusDataItemList
	list := models.PrometheusDataSetResponse{
		PromItemList: make([]models.PrometheusDataSetResponseItem, 0, len(body)),
		PromDataType: queryType,
	}

	// Iterate over the byte slice and fill the list
	for _, b := range body {
		list.PromItemList = append(list.PromItemList, models.PrometheusDataSetResponseItem{Metric: float32(b)})
	}

	return &list, nil
}
