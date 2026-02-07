package handlers

import (
	"context"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient interface {
	GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
	DeleteAccount(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error)
	ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error)
	Close() error
}

type grpcUserClient struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func NewGRPCUserClient(addr string) (UserServiceClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcUserClient{
		conn:   conn,
		client: userv1.NewUserServiceClient(conn),
	}, nil
}

func (c *grpcUserClient) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return c.client.GetUser(ctx, req)
}

func (c *grpcUserClient) DeleteAccount(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	return c.client.DeleteUser(ctx, req)
}

func (c *grpcUserClient) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, req)
}

func (c *grpcUserClient) Close() error {
	return c.conn.Close()
}
