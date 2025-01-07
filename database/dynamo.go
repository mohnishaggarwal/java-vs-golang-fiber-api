package database

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"os"
)

var DynamoClient *dynamodb.Client

func InitDynamo() {
	env := os.Getenv("env")
	region := os.Getenv("region")

	if env == "local" {
		log.Println("Using local DynamoDB instance ")

		endpoint := "http://localhost:8000" // Ensure DynamoDB Local is running at this endpoint
		DynamoClient = dynamodb.New(dynamodb.Options{
			Region:       region,
			Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
			BaseEndpoint: aws.String(endpoint),
		})

		tableName := "Products"
		if !DoesTableExist(&tableName) {
			err := CreateDynamoTable(&tableName)
			if err != nil {
				panic("Unable to create table")
			}
		}
		if DoesTableExist(&tableName) {
			log.Println("Table exists, all set for requests in dev mode")
		} else {
			panic("Table does not exist.")
		}
	} else {
		fmt.Println("we really getting here?")
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		DynamoClient = dynamodb.NewFromConfig(cfg)
		log.Println("Using AWS DynamoDB instance")
	}
}

func DoesTableExist(tableName *string) bool {
	out, err := DynamoClient.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("failed to list tables: %v", err)
		return false
	}
	return contains(out.TableNames, tableName)
}

func CreateDynamoTable(tableName *string) error {
	createInput := &dynamodb.CreateTableInput{
		TableName: aws.String(*tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       types.KeyTypeHash, // Partition key
			},
		},
		// For local DynamoDB, provisioned throughput is still required in the request.
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := DynamoClient.CreateTable(context.TODO(), createInput)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Successfully created table 'Products'.")
	return nil
}

func contains(slice []string, item *string) bool {
	for _, v := range slice {
		if v == *item {
			return true
		}
	}
	return false
}
