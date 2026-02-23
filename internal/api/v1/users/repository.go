package users

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository struct {
	db *dynamodb.Client
}

func NewRepository(db *dynamodb.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	isPkAbsent := expression.AttributeNotExists(expression.Name("PK"))
	expr, err := expression.NewBuilder().WithCondition(isPkAbsent).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 aws.String(TableName),
		Item:                      item,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to put user item: %w", err)
	}

	return nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	keyCond := expression.Key("Email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query expression: %w", err)
	}

	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		IndexName:                 aws.String(EmailIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by email: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	fullResult, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"PK": result.Items[0]["PK"],
			"SK": &types.AttributeValueMemberS{Value: SKMetadata},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get full user item: %w", err)
	}

	if fullResult.Item == nil {
		return nil, nil
	}

	var user User
	if err := attributevalue.UnmarshalMap(fullResult.Item, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UpdateFields(ctx context.Context, pk string, fields *UserUpdate) error {
	item, err := attributevalue.MarshalMap(fields)
	if err != nil {
		return fmt.Errorf("failed to marshal update fields: %w", err)
	}

	update := expression.UpdateBuilder{}
	for key, val := range item {
		update = update.Set(expression.Name(key), expression.Value(val))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build update expression: %w", err)
	}

	_, err = r.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: SKMetadata},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to update user fields: %w", err)
	}

	return nil
}
