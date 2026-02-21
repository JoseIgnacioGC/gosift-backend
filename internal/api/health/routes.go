package health

import (
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, awsClient *aws.Client) {
	router.GET("/health", getHealthCheck(awsClient))
}
