package users

import (
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, awsClient *aws.Client) {
	// users := group.Group("/users")
}
