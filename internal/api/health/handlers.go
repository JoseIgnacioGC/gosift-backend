package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/gin-gonic/gin"
)

func getHealthCheck(awsClient *aws.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := awsClient.DynamoDB.ListTables(ctx, nil)
		dbStatus := "connected"
		if err != nil {
			dbStatus = fmt.Sprintf("[ERROR]: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "active",
			"platform": "LocalStack Pro",
			"database": dbStatus,
		})
	}
}
