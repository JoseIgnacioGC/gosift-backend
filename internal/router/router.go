package router

import (
	"github.com/JoseIgnacioGC/gosift-backend/internal/api/health"
	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/auth"
	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/subscriptions"
	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/users"
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, awsClient *aws.Client, jwtSecret string) {
	health.RegisterRoutes(router, awsClient)

	v1Group := router.Group("/api/v1")
	auth.RegisterRoutes(v1Group, awsClient, jwtSecret)
	users.RegisterRoutes(v1Group, awsClient)
	subscriptions.RegisterRoutes(v1Group, awsClient, jwtSecret)
}
