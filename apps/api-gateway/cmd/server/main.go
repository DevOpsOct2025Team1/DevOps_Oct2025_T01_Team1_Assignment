package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/config"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Using auth-service at: %s", cfg.AuthServiceAddr)

	router := gin.Default()

	authHandler := handlers.NewAuthHandler(cfg.AuthServiceAddr)
	router.POST("/api/signup", authHandler.SignUp)
	router.POST("/api/login", authHandler.Login)
	router.POST("/api/validate", authHandler.ValidateToken)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("API Gateway listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
