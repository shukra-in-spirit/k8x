package clients

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type NewLambdaInterface interface {
	TriggerCreateLambdaWithEvent(data []byte, functionName string) (*LambdaRespBody, error)
}

type lambdaConfig struct {
	client lambdaiface.LambdaAPI
}

type LambdaRespBody struct {
	CPU      string `json:"cpu"`
	Memory   string `json:"memory"`
	Replicas string `json:"replicas"`
}

type LambdaResponse struct {
	StatusCode int            `json:"statusCode"`
	Body       LambdaRespBody `json:"body"`
}

func NewLamdaClient(client lambdaiface.LambdaAPI) *lambdaConfig {
	return &lambdaConfig{
		client: client,
	}
}

func (e *lambdaConfig) TriggerCreateLambdaWithEvent(data []byte, functionName string) (*LambdaRespBody, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      data,
	}

	output, err := e.client.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger lambda. %v", err)
	}

	var resp LambdaResponse

	err = json.Unmarshal(output.Payload, &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling lambda function response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed getting response from lambda, got statusCode: " + strconv.Itoa(resp.StatusCode))
	}

	return &resp.Body, nil
}