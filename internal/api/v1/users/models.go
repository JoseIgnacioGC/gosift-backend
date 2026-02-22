package users

import "time"

type User struct {
	PK            string    `dynamodbav:"PK"`
	SK            string    `dynamodbav:"SK"`
	Email         string    `dynamodbav:"Email"`
	Username      string    `dynamodbav:"Username"`
	PasswordHash  string    `dynamodbav:"PasswordHash"`
	Name          string    `dynamodbav:"Name"`
	AvatarURL     string    `dynamodbav:"AvatarURL"`
	Timezone      string    `dynamodbav:"Timezone"`
	Language      string    `dynamodbav:"Language"`
	EmailVerified bool      `dynamodbav:"EmailVerified"`
	LastLoginAt   time.Time `dynamodbav:"LastLoginAt,unixtime"`
	CreatedAt     time.Time `dynamodbav:"CreatedAt,unixtime"`
	UpdatedAt     time.Time `dynamodbav:"UpdatedAt,unixtime"`
}

const (
	TableName  = "Users"
	PKPrefix   = "USER#"
	SKMetadata = "METADATA"
	EmailIndex = "EmailIndex"
)
