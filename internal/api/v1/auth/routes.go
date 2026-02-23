package auth

import (
	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/users"
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, awsClient *aws.Client, jwtSecret string) {
	userRepo := users.NewRepository(awsClient.DynamoDB)
	authService := newService(userRepo, jwtSecret)

	auth := group.Group("/auth")

	auth.POST("/register", register(authService))
	auth.POST("/login", login(authService))

}
