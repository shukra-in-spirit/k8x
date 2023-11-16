package helpers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/shukra-in-spirit/k8x/internal/clients"
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
func TriggerLocalCreateLambdaWithEvent(data []byte, functionName string) (*clients.LambdaRespBody, error) {
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
		return &clients.LambdaRespBody{}, fmt.Errorf("error running local lambda: %v", err)
	}

	var response clients.LambdaRespBody
	err = json.Unmarshal(output, &response)
	if err != nil {
		return &clients.LambdaRespBody{}, fmt.Errorf("error while unmarshalling lambda output: %v", err)
	}
	return &response, nil
}
