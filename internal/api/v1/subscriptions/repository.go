package subscriptions

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/users"
)

type Repository struct {
	db *dynamodb.Client
}

func NewRepository(db *dynamodb.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, sub *Subscription) error {
	item, err := attributevalue.MarshalMap(sub)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}

	isPkAbsent := expression.AttributeNotExists(expression.Name("PK"))
	expr, err := expression.NewBuilder().WithCondition(isPkAbsent).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 aws.String(users.TableName),
		Item:                      item,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to put subscription item: %w", err)
	}

	return nil
}

func (r *Repository) ListByUser(ctx context.Context, userPK string) ([]Subscription, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(userPK)),
		expression.Key("SK").BeginsWith(SKPrefix),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query expression: %w", err)
	}

	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(users.TableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}

	subs := make([]Subscription, 0, len(result.Items))
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &subs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscriptions: %w", err)
	}

	return subs, nil
}

func (r *Repository) FindByUserAndFeedURL(ctx context.Context, userPK string, feedURL string) (*Subscription, error) {
	keyCond := expression.Key("FeedURL").Equal(expression.Value(feedURL))
	filter := expression.Name("PK").Equal(expression.Value(userPK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithFilter(filter).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query expression: %w", err)
	}

	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(users.TableName),
		IndexName:                 aws.String(FeedURLIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by feed URL: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var sub Subscription
	if err := attributevalue.UnmarshalMap(result.Items[0], &sub); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return &sub, nil
}

func (r *Repository) UpdateFields(ctx context.Context, pk string, sk string, fields *SubscriptionUpdate) error {
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
		TableName: aws.String(users.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to update subscription fields: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, pk string, sk string) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(users.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
