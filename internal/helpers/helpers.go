package helpers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/models"
)

func DecomposeServiceID(service_id string) (string, string) {
	// Split the input string based on the '-' delimiter
	parts := strings.Split(service_id, "-")

	// Check if there are at least two parts (containerName and namespace)
	if len(parts) < 2 {
		return "", ""
	}

	// Extract containerName and namespace
	containerName := parts[0]
	namespace := parts[1]

	return containerName, namespace
}

// call the local lambda methods
func TriggerLocalCreateLambdaWithEvent(data string, functionName string) (*models.LambdaRespBody, error) {
	var lambdaFilePath string
	// Command to run the Python script
	if functionName == "c" {
		lambdaFilePath = "./lambda/create_lambda/lambda_function.py"
	} else {
		lambdaFilePath = "./lambda/predict_lambda/local_trigger.py"
	}
	cmd := exec.Command("python", lambdaFilePath, data)

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
	randomID := 0
	for index, cpu := range l1 {
		// for _, mem := range l2 {
		// 	ID =
		// 	if cpu.ID == mem.ID {
		// 		promData = append(promData, &models.PromData{ServiceID: id, Timestamp: cpu.Timestamp, CPU: cpu.Metric, Memory: mem.Metric})
		// 	}
		mem := l2[index]
		randomID = randomID + 1
		fmt.Printf("CPU VALUE: %v, MEMORY VALUE: %v, ID: %v\n\n", cpu.Metric, mem.Metric, randomID)
		promData = append(promData, &models.PromData{ServiceID: id, ID: strconv.Itoa(randomID), CPU: cpu.Metric, Memory: mem.Metric})
		// }
	}

	return promData
}

func PrepareHistoryData(l1 []models.PrometheusDataSetResponseItem, l2 []models.PrometheusDataSetResponseItem) []*models.History {
	promData := []*models.History{}

	for _, cpu := range l1 {
		for _, mem := range l2 {
			// if cpu.Timestamp == mem.Timestamp {
			// 	promData = append(promData, &models.History{Timestamp: cpu.Timestamp, CPU: cpu.Metric, Memory: mem.Metric})
			// }
			promData = append(promData, &models.History{CPU: cpu.Metric, Memory: mem.Metric})
		}
	}

	return promData
}

// GetEnvOrDefault will return env value of given variable.
func GetEnvOrDefault(envVar, defaultValue string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}

	return defaultValue
}

// GenerateRandomCode method is used to generate a random character of input length
// it is being used to generate a sort key in the dynamodb table
func GenerateRandomCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
