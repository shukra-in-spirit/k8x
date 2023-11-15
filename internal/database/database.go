package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/shukra-in-spirit/k8x/internal/config"
)

func CreateSession(awsConfig config.AWSConfig) (dynamodbiface.DynamoDBAPI, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsConfig.AWSRegion),
		Credentials: credentials.NewStaticCredentials(awsConfig.AWSAccessID, awsConfig.AWSSecretKey, ""),
		Endpoint:    aws.String(awsConfig.DBEndpoint),
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating a session. %v", err)
	}

	return dynamodb.New(sess), nil
}
