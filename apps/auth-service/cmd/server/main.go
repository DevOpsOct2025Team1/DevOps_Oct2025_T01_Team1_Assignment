package main

import (
	"fmt"
	"log"
	"net"

	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/client"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/config"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/health"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/jwt"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/service"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Using user-service at: %s", cfg.UserServiceAddr)

	userClient, err := client.NewUserServiceClient(cfg.UserServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}
	defer userClient.Close()

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, service.NewAuthServiceServer(userClient, jwtManager))
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewHealthServer())
	reflection.Register(grpcServer)

	log.Printf("Auth service listening on :%s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
