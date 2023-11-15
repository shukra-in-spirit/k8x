package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/shukra-in-spirit/k8x/internal/config"
)

func CreateSession(awsConfig config.AWSConfig) (*session.Session, error) {
	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String(awsConfig.AWSRegion),
	// 	Credentials: credentials.NewStaticCredentials(awsConfig.AWSAccessID, awsConfig.AWSSecretKey, ""),
	// 	Endpoint:    aws.String(awsConfig.DBEndpoint),
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed creating a session. %v", err)
	// }

	// return sess, nil

	return session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithRegion(awsConfig.AWSRegion),
		SharedConfigState: session.SharedConfigEnable,
	})), nil
}
