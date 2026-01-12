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

func (c *UserServiceClient) CreateUser(ctx context.Context, username, hashedPassword, role string) (*userv1.User, error) {
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

func (c *UserServiceClient) GetUserByUsername(ctx context.Context, username string) (*userv1.User, string, error) {
	resp, err := c.client.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, "", err
	}
	return resp.User, resp.HashedPassword, nil
}
