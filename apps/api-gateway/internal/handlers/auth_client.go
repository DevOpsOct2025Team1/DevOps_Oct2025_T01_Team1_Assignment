package handlers

import (
	"context"

	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient interface {
	SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error)
	Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
	ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error)
	Close() error
}

type grpcAuthClient struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func NewGRPCAuthClient(addr string) (AuthServiceClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcAuthClient{
		conn:   conn,
		client: authv1.NewAuthServiceClient(conn),
	}, nil
}

func (c *grpcAuthClient) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	return c.client.SignUp(ctx, req)
}

func (c *grpcAuthClient) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

func (c *grpcAuthClient) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, req)
}

func (c *grpcAuthClient) Close() error {
	return c.conn.Close()
}
