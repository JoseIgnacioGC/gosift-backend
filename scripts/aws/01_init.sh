#!/bin/bash
echo "LocalStack custom script: Initializing LocalStack..."

# DynamoDB — Users table
awslocal dynamodb create-table \
    --table-name Users \
    --attribute-definitions \
        AttributeName=PK,AttributeType=S \
        AttributeName=SK,AttributeType=S \
        AttributeName=Email,AttributeType=S \
    --key-schema \
        AttributeName=PK,KeyType=HASH \
        AttributeName=SK,KeyType=RANGE \
    --global-secondary-indexes \
        '[
            {
                "IndexName": "EmailIndex",
                "KeySchema": [
                    {"AttributeName": "Email", "KeyType": "HASH"}
                ],
                "Projection": {
                    "ProjectionType": "INCLUDE",
                    "NonKeyAttributes": ["PK", "PasswordHash"]
                },
                "ProvisionedThroughput": {
                    "ReadCapacityUnits": 5,
                    "WriteCapacityUnits": 5
                }
            }
        ]' \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

# S3 Bucket (Media Storage)
awslocal s3 mb s3://gosift-media

echo "LocalStack custom script: Resources created successfully."