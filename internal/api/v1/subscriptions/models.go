package subscriptions

import "time"

type Subscription struct {
	PK        string    `dynamodbav:"PK"`
	SK        string    `dynamodbav:"SK"`
	FeedURL   string    `dynamodbav:"FeedURL"`
	Title     string    `dynamodbav:"Title"`
	SiteURL   string    `dynamodbav:"SiteURL"`
	Category  string    `dynamodbav:"Category"`
	IsActive  bool      `dynamodbav:"IsActive"`
	CreatedAt time.Time `dynamodbav:"CreatedAt,unixtime"`
	UpdatedAt time.Time `dynamodbav:"UpdatedAt,unixtime"`
}

type SubscriptionUpdate struct {
	Title     *string    `dynamodbav:"Title,omitempty"`
	Category  *string    `dynamodbav:"Category,omitempty"`
	IsActive  *bool      `dynamodbav:"IsActive,omitempty"`
	UpdatedAt *time.Time `dynamodbav:"UpdatedAt,unixtime,omitempty"`
}

const (
	SKPrefix     = "SUB#"
	FeedURLIndex = "FeedURLIndex"
)
