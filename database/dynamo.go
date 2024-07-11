package database

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

var DynamoClient *dynamodb.Client

func InitDynamo() {
	//env := os.Getenv("env")
	env := "local"
	//region := os.Getenv("region")
	region := "us-east-1"

	if env == "local" {
		fmt.Println("Getting here")
		endpoint := "http://localhost:8000" // Ensure DynamoDB Local is running at this endpoint
		DynamoClient = dynamodb.New(dynamodb.Options{
			Region:           region,
			Credentials:      aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
			EndpointResolver: dynamodb.EndpointResolverFromURL(endpoint),
		})
		log.Println("Using local DynamoDB instance")
	} else {
		fmt.Println("we really getting here?")
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		DynamoClient = dynamodb.NewFromConfig(cfg)
		log.Println("Using AWS DynamoDB instance")
	}
}
