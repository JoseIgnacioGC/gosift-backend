package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	DynamoDB *dynamodb.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	awsRegion := "us-east-1"

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if os.Getenv("USE_LOCALSTACK") == "true" {
			o.BaseEndpoint = aws.String("http://localhost:4566")
		}
	})

	return &Client{
		DynamoDB: dynamoClient,
	}, nil
}
