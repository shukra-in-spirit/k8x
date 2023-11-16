package helpers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/shukra-in-spirit/k8x/internal/models"
)

func DecomposeServiceID(service_id string) (string, string) {
	// Split the input string based on the '-' delimiter
	parts := strings.Split(service_id, "-")

	// Check if there are at least two parts (serviceName and namespace)
	if len(parts) < 2 {
		return "", ""
	}

	// Extract serviceName and namespace
	serviceName := parts[0]
	namespace := parts[1]

	return serviceName, namespace
}

// call the local lambda methods
func TriggerLocalCreateLambdaWithEvent(data []byte, functionName string) (*models.LambdaRespBody, error) {
	var lambdaFilePath string
	// Command to run the Python script
	if functionName == "c" {
		lambdaFilePath = "./lambda/create_lambda/local_trigger.py"
	} else {
		lambdaFilePath = "./lambda/predict_lambda/local_trigger.py"
	}
	eventString := string(data)
	cmd := exec.Command("python", lambdaFilePath, eventString)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &models.LambdaRespBody{}, fmt.Errorf("error running local lambda: %v", err)
	}

	var response models.LambdaRespBody
	err = json.Unmarshal(output, &response)
	if err != nil {
		return &models.LambdaRespBody{}, fmt.Errorf("error while unmarshalling lambda output: %v", err)
	}
	return &response, nil
}

func ProcessPromData(id string, l1 []models.PrometheusDataSetResponseItem, l2 []models.PrometheusDataSetResponseItem) []*models.PromData {
	promData := []*models.PromData{}

	for _, cpu := range l1 {
		for _, mem := range l2 {
			if cpu.Timestamp == mem.Timestamp {
				promData = append(promData, &models.PromData{ServiceID: id, Timestamp: cpu.Timestamp, CPU: cpu.Metric, Memory: mem.Metric})
			}
		}
	}

	return promData
}

func PrepareHistoryData(l1 []models.PrometheusDataSetResponseItem, l2 []models.PrometheusDataSetResponseItem) []*models.History {
	promData := []*models.History{}

	for _, cpu := range l1 {
		for _, mem := range l2 {
			if cpu.Timestamp == mem.Timestamp {
				promData = append(promData, &models.History{Timestamp: cpu.Timestamp, CPU: cpu.Metric, Memory: mem.Metric})
			}
		}
	}

	return promData
}