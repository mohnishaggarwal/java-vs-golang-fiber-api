#!/bin/sh

# Wait for DynamoDB Local to be ready
until curl -s http://dynamodb-local:8000; do
  echo "Waiting for dynamodb-local..."
  sleep 1
done

echo "DynamoDB Local is up - proceeding with table creation"

# Create Products table for demo app to use
aws dynamodb create-table \
    --table-name Products \
    --attribute-definitions AttributeName=Id,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:8000 \
    --region us-east-1

# Start the Spring Boot application
echo "Starting up application!"
java -jar /app/app.jar
