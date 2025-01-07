package repository

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mohnishaggarwal/products/database"
	"github.com/mohnishaggarwal/products/models"
	"log"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, id string, product *models.Product) error
}

type productRepository struct {
	db *dynamodb.Client
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		db: database.DynamoClient,
	}
}

func (r *productRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
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
	log.Println(result.Item)
	if len(result.Item) == 0 {
		return nil, errors.New("no product found")
	}
	err = attributevalue.UnmarshalMap(result.Item, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		log.Println("Did we fail to unwrap?")
		return err
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("goservice_products_table"),
		Item:      item,
	})
	return err
}

func (r *productRepository) UpdateProduct(ctx context.Context, id string, product *models.Product) error {
	product.Id = id // Ensure the product ID is set to the provided ID
	return r.CreateProduct(ctx, product)
}
