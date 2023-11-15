package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
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

	responseDataFrame, err := prom.GetPrometheusDataWithinRange(ctx, promQuery, startTime, endTime, "")
	if err != nil {
		return &models.PrometheusDataSetResponse{}, nil
	}
	return responseDataFrame, nil
}

func (prom *PrometheusInstance) GetPrometheusDataWithinRange(ctx context.Context, promQuery string, startTime time.Time, endTime time.Time, steps string) (*models.PrometheusDataSetResponse, error) {
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
		return nil, err
	}

	// Create a POST request to the Prometheus query_range API
	req, err := http.NewRequest("POST", prom.url, bytes.NewBuffer(queryRangeBytes))
	if err != nil {
		return nil, err
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into the QueryResult struct
	var result models.PrometheusDataSetResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
