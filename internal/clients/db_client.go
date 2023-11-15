package clients

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/shukra-in-spirit/k8x/internal/models"
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
	AddDataBatch(dataList *[]models.PrometheusDataSetResponseItem, serviceID string) error
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
		return fmt.Errorf("failed marshalling input to dynamoDB attribute format. %v", err)
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

func (e *promStore) AddDataBatch(dataList *[]models.PrometheusDataSetResponseItem, serviceID string) error {
	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	for start < len(*dataList) {
		var writeReqs []*dynamodb.WriteRequest
		if end > len(*dataList) {
			end = len(*dataList)
		}

		for _, data := range (*dataList)[start:end] {
			dataItem := &PromData{
				ServiceID: serviceID,
				Timestamp: data.Timestamp,
				CPU:       data.CPU,
				Memory:    data.Memory,
			}
			item, err := dynamodbattribute.MarshalMap(dataItem)
			if err != nil {
				log.Printf("error in marshalling data %v for batch writing: %v\n", dataItem, err)
			} else {
				writeReqs = append(
					writeReqs,
					&dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: item}},
				)
			}
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{tableName: writeReqs},
		}

		_, err := e.client.BatchWriteItem(input)
		if err != nil {
			log.Printf("failed to perform batch write: %v\n", err)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	log.Printf("successfully performed batch write on %v, got= %v\n", written, len(*dataList))

	return nil
}
