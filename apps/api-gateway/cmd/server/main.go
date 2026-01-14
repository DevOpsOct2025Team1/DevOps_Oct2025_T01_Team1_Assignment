package main

import (
	"log"

	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/config"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Using auth-service at: %s", cfg.AuthServiceAddr)

	authClient, err := handlers.NewGRPCAuthClient(cfg.AuthServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	userClient, err := handlers.NewGRPCUserClient(cfg.UserServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}

	srv := server.New(authClient, userClient)
	defer srv.Close()

	log.Printf("API Gateway listening on :%s", cfg.Port)
	if err := srv.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
