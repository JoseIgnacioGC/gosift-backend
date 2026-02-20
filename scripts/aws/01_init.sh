#!/bin/bash
echo "LocalStack custom script: Initializing LocalStack..."

# DynamoDB 
awslocal dynamodb create-table \
    --table-name NewsItems \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

# S3 Bucket (Media Storage)
awslocal s3 mb s3://gosift-media

echo "LocalStack custom script: Resources created successfully."