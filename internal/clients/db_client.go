package clients

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/shukra-in-spirit/k8x/internal/models"
)

const (
	tableName = "b_k8x_t"
)

type PromStorer interface {
	AddData(data *models.PromData) error
	AddDataBatch(dataList []*models.PromData, serviceID string) error
}

type promStore struct {
	client dynamodbiface.DynamoDBAPI
}

func NewPromStore(client dynamodbiface.DynamoDBAPI) *promStore {
	return &promStore{
		client: client,
	}
}

func (e *promStore) AddData(data *models.PromData) error {
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

func (e *promStore) AddDataBatch(dataList []*models.PromData, serviceID string) error {
	// for _, data := range dataList {
	// 	fmt.Printf("ID: %s, ServiceID: %s, CPU: %.2f, Memory: %.2f\n", data.ID, data.ServiceID, data.CPU, data.Memory)
	// }
	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize

	for start < len(dataList) {
		var writeReqs []*dynamodb.WriteRequest
		if end > len(dataList) {
			end = len(dataList)
		}

		for _, data := range dataList[start:end] {
			item, err := dynamodbattribute.MarshalMap(data)
			if err != nil {
				return fmt.Errorf("error in marshalling data %v for batch writing: %v", data, err)
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
			return fmt.Errorf("failed to perform batch write: %v", err)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	log.Printf("successfully performed batch write on %v, got= %v\n", written, len(dataList))

	return nil
}
