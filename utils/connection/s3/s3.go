package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	credentialsFactory "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type Client struct {
	Session    *session.Session
	Connection *s3.S3
}

func NewClient(path string) *Client {
	s3Config := viper.GetStringMapString("connections.s3." + path)
	credentials := credentialsFactory.NewStaticCredentials(s3Config["id"], s3Config["key"], "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      aws.String(s3Config["region"]),
		Endpoint:    aws.String(s3Config["entrypoint"]),
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
