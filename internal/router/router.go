package router

import (
	"github.com/JoseIgnacioGC/gosift-backend/internal/api/health"
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, awsClient *aws.Client) {
	health.RegisterRoutes(router, awsClient)

	// v1Group := router.Group("/api/v1")
	// something.RegisterRoutes(v1Group, awsClient)

}
