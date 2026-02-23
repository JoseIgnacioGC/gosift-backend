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

type UserUpdate struct {
	Email         *string    `dynamodbav:"Email,omitempty"`
	Username      *string    `dynamodbav:"Username,omitempty"`
	PasswordHash  *string    `dynamodbav:"PasswordHash,omitempty"`
	Name          *string    `dynamodbav:"Name,omitempty"`
	AvatarURL     *string    `dynamodbav:"AvatarURL,omitempty"`
	Timezone      *string    `dynamodbav:"Timezone,omitempty"`
	Language      *string    `dynamodbav:"Language,omitempty"`
	EmailVerified *bool      `dynamodbav:"EmailVerified,omitempty"`
	LastLoginAt   *time.Time `dynamodbav:"LastLoginAt,unixtime,omitempty"`
	UpdatedAt     *time.Time `dynamodbav:"UpdatedAt,unixtime,omitempty"`
}

const (
	TableName  = "Users"
	PKPrefix   = "USER#"
	SKMetadata = "METADATA"
	EmailIndex = "EmailIndex"
)
