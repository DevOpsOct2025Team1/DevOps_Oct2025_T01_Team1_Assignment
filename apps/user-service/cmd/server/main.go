package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/provsalt/DOP_P01_Team1/common/telemetry"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/config"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/health"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/service"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/store"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.AxiomToken == "" {
		log.Printf("Tracing disabled: AXIOM_API_TOKEN is empty")
	} else {
		shutdown, err := telemetry.InitTracer(context.Background(), telemetry.Config{
			ServiceName: "user-service",
			Environment: cfg.Environment,
			Token:       cfg.AxiomToken,
			Endpoint:    cfg.AxiomEndpoint,
			Dataset:     cfg.AxiomDataset,
		})
		if err != nil {
			log.Printf("Tracing disabled: failed to initialize tracer: %v", err)
		} else {
			defer shutdown(context.Background())
		}
	}

	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	database := client.Database(cfg.MongoDBDatabase)
	userStore := store.NewUserStore(database)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	userv1.RegisterUserServiceServer(grpcServer, service.NewUserServiceServer(userStore))
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewHealthServer())
	reflection.Register(grpcServer)

	log.Printf("User service listening on :%s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
