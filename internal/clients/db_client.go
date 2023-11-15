package clients

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	tableName = "b_k8x_t"
)

type PromData struct {
	ServiceID string    `json:"service_id"`
	Timestamp time.Time `json:"timestamp"`
	CPU       float32   `json:"cpu"`
	Memory    float32   `json:"memory"`
}

type PromStorer interface {
	AddData(data *PromData) error
}

type promStore struct {
	client dynamodbiface.DynamoDBAPI
}

func NewPromStore(client dynamodbiface.DynamoDBAPI) *promStore {
	return &promStore{
		client: client,
	}
}

func (e *promStore) AddData(data *PromData) error {
	inputItems, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return fmt.Errorf("failed marshalling input to dyanmoDB attribute format. %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      inputItems,
	}

	_, err = e.client.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put items. %v", err)
	}

	return nil
}
