package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mohnishaggarwal/products/database"
	"github.com/mohnishaggarwal/products/models"
	"log"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*models.ProductOutput, error)
	CreateProduct(ctx context.Context, product *models.Product) (string, error)
	UpdateProduct(ctx context.Context, id string, product *models.Product) (string, error)
}

type productRepository struct {
	db *dynamodb.Client
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		db: database.DynamoClient,
	}
}

func (r *productRepository) GetProductByID(ctx context.Context, id string) (*models.ProductOutput, error) {
	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("goservice_products_table"),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		log.Println("Error getting product " + err.Error())
		return nil, err
	}

	var product models.Product
	if len(result.Item) == 0 {
		return nil, errors.New("product not found")
	}
	err = attributevalue.UnmarshalMap(result.Item, &product)
	if err != nil {
		return nil, err
	}

	return transformProduct(&product), nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) (string, error) {
	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		log.Println("Did we fail to unwrap?")
		return "", err
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("goservice_products_table"),
		Item:      item,
	})
	if err != nil {
		return "", err
	}
	return "Product added successfully", nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, id string, product *models.Product) (string, error) {
	product.Id = id // Ensure the product ID is set to the provided ID
	_, err := r.CreateProduct(ctx, product)
	if err != nil {
		return "", err
	} else {
		return "Product updated successfully", err
	}
}

func transformProduct(product *models.Product) *models.ProductOutput {
	expenseCategory := getExpenseCasegory(product.Price)
	url := fmt.Sprintf("https://example.com/product/%s", product.Id)
	return &models.ProductOutput{
		Product:         product,
		Url:             url,
		ExpenseCategory: expenseCategory,
	}
}

func getExpenseCasegory(price float32) models.ExpenseCategory {
	if price < 10 {
		return models.VeryCheap
	} else if price < 100 {
		return models.Cheap
	} else {
		return models.Expensive
	}
}
