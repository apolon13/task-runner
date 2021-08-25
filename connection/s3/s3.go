package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	credentialsFactory "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"task-runner/config"
)

type Client struct {
	Session    *session.Session
	Connection *s3.S3
}

func NewClient(config config.Yaml) *Client {
	credentials := credentialsFactory.NewStaticCredentials(config.Connections.S3.Id, config.Connections.S3.Key, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      aws.String(config.Connections.S3.Region),
		Endpoint:    aws.String(config.Connections.S3.Entrypoint),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),

	})
	if err != nil {
		panic(err)
	}

	return &Client{
		Connection: s3.New(sess),
		Session:    sess,
	}
}
