package store

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// EnsureDefaultAdmin creates a default admin user if the database is empty.
// This should be called during service startup to ensure there's always an admin user available.
func (u *UserStore) EnsureDefaultAdmin(ctx context.Context, username, password string) error {
	collection := u.database.Collection("users")

	err := collection.FindOne(ctx, bson.M{}).Err()
	if err == nil {
		return nil // Found a user, skip initialization
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err // Real error occurred
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	_, err = u.CreateUser(ctx, &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Role:           "admin",
	})

	return err
}
