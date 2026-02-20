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

// NewClient creates an AWS client that automatically detects LocalStack
func NewClient(ctx context.Context) (*Client, error) {
	// Custom Resolver: Redirects AWS calls to LocalStack if env var is set
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if os.Getenv("USE_LOCALSTACK") == "true" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: "us-east-1",
			}, nil
		}
		// Returning EndpointNotFoundError allows the SDK to fallback to real AWS
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Load Config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	return &Client{
		DynamoDB: dynamodb.NewFromConfig(cfg),
	}, nil
}
