package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/JoseIgnacioGC/gosift-backend/internal/api/handlers"
	"github.com/JoseIgnacioGC/gosift-backend/internal/config"
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
)

func main() {
	ginConfig := config.Get()

	ctx := context.Background()
	awsClient, err := aws.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize AWS: %v", err)
	}

	router := gin.Default()
	config.ConfigureProxies(router, ginConfig)

	router.GET("/health", handlers.HealthCheck(awsClient))

	log.Printf("gosift-backend running on http://localhost:%v\n", ginConfig.Port)
	router.Run(":" + ginConfig.Port)
}
