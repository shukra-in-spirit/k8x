package database

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/shukra-in-spirit/k8x/internal/config"
)

const (
	tableName = "b_k8x_t"
)

type PromData struct {
	ServiceID string    `env:"service_id"`
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

func NewPromStore(conf config.AWSConfig) *promStore {
	newPromStore := promStore{}

	session, err := CreateSession(conf)
	if err != nil {
		panic(err)
	}
	newPromStore.client = session

	return &newPromStore
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
