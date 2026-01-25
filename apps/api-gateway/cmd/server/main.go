package main

import (
	"context"
	"log"

	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/config"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/server"
	"github.com/provsalt/DOP_P01_Team1/common/telemetry"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.AxiomToken == "" {
		log.Printf("Tracing disabled: AXIOM_API_TOKEN is empty")
	} else {
		shutdown, err := telemetry.InitTelemetry(context.Background(), telemetry.Config{
			ServiceName:    "api-gateway",
			Environment:    cfg.Environment,
			Token:          cfg.AxiomToken,
			Endpoint:       cfg.AxiomEndpoint,
			Dataset:        cfg.AxiomDataset,
			MetricsDataset: cfg.AxiomMetricsDataset,
		})
		if err != nil {
			log.Printf("Tracing disabled: failed to initialize tracer: %v", err)
		} else {
			defer shutdown(context.Background())
		}
	}

	log.Printf("Using auth-service at: %s", cfg.AuthServiceAddr)
	log.Printf("Using user-service at: %s", cfg.UserServiceAddr)

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
