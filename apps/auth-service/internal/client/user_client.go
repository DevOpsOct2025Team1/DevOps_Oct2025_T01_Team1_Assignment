package client

import (
	"context"
	"fmt"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	client userv1.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userv1.NewUserServiceClient(conn)

	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

func (c *UserServiceClient) CreateUser(ctx context.Context, username, hashedPassword string, role userv1.Role) (*userv1.User, error) {
	resp, err := c.client.CreateUser(ctx, &userv1.CreateUserRequest{
		Username:       username,
		HashedPassword: hashedPassword,
		Role:           role,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (c *UserServiceClient) VerifyPassword(ctx context.Context, username, password string) (*userv1.User, bool, error) {
	resp, err := c.client.VerifyPassword(ctx, &userv1.VerifyPasswordRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, false, err
	}
	return resp.User, resp.Valid, nil
}
