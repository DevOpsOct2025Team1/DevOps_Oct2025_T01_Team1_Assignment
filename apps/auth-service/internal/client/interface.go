package client

import (
	"context"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

// UserClient defines operations for user management.
// CreateUser expects the provided password to be already hashed (e.g., bcrypt).
type UserClient interface {
	CreateUser(ctx context.Context, username, hashedPassword string, role userv1.Role) (*userv1.User, error)
	VerifyPassword(ctx context.Context, username, password string) (*userv1.User, bool, error)
	Close() error
}
