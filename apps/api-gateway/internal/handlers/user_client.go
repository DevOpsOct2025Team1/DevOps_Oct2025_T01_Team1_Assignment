package handlers

import (
	"context"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient interface {
	DeleteAccount(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error)
	Close() error
}

type grpcUserClient struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func NewGRPCUserClient(addr string) (UserServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcUserClient{
		conn:   conn,
		client: userv1.NewUserServiceClient(conn),
	}, nil
}

func (c *grpcUserClient) DeleteAccount(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	return c.client.DeleteUserByUserId(ctx, req)
}

func (c *grpcUserClient) Close() error {
	return c.conn.Close()
}
