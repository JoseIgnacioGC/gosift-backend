package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/JoseIgnacioGC/gosift-backend/internal/config"
	"github.com/JoseIgnacioGC/gosift-backend/internal/platform/aws"
	"github.com/JoseIgnacioGC/gosift-backend/internal/router"
)

func main() {
	ginConfig := config.Get()

	awsClient, err := aws.NewClient(context.Background())
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize AWS: %v", err)
	}

	mainRouter := gin.Default()

	config.ConfigureProxies(mainRouter, ginConfig)
	router.RegisterRoutes(mainRouter, awsClient, ginConfig.JWTSecret)

	log.Printf("[INFO] running on http://localhost:%v\n", ginConfig.Port)
	mainRouter.Run(":" + ginConfig.Port)

}
