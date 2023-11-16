package clients

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

type NewLambdaInterface interface {
	TriggerLambdaWithEvent(data []byte, functionName string) (*models.LambdaRespBody, error)
}

type lambdaConfig struct {
	client lambdaiface.LambdaAPI
}

func NewLamdaClient(client lambdaiface.LambdaAPI) *lambdaConfig {
	return &lambdaConfig{
		client: client,
	}
}

func (e *lambdaConfig) TriggerLambdaWithEvent(data []byte, functionName string) (*models.LambdaRespBody, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      data,
	}

	output, err := e.client.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger lambda. %v", err)
	}

	var resp models.LambdaResponse

	err = json.Unmarshal(output.Payload, &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling lambda function response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed getting response from lambda, got statusCode: " + strconv.Itoa(resp.StatusCode))
	}

	return &resp.Body, nil
}

// TODO:
// UI - create the popup forms for both
// UI - integrations with the forms
// Backend - Complete code review
// Backend - Start testing end to end
// infrastructure - installed helm, python, k8s, mysql / kubectl - deployment.yaml
// Data - clean -  5 services
// Visualization - UI - live data in card
//                    - (ns) actual values in pod
// Visualization - UI - hardcoded - simulation
// Audit Logs    - UI
// VIDEO - 1 hour
// PRESENTATION - 10-12
