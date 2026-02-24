package subscriptions

import (
	"github.com/JoseIgnacioGC/gosift-backend/internal/middleware"
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, awsClient *aws.Client, jwtSecret string) {
	subRepo := NewRepository(awsClient.DynamoDB)
	svc := newService(subRepo)

	subs := group.Group("/subscriptions")
	subs.Use(middleware.Auth(jwtSecret))

	subs.POST("", createSubscription(svc))
	subs.GET("", listSubscriptions(svc))
	subs.PATCH("/:id", updateSubscription(svc))
	subs.DELETE("/:id", deleteSubscription(svc))
}
