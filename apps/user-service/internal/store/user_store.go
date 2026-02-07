package store

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type User struct {
	Id             string `bson:"_id,omitempty" json:"id"`
	Username       string `bson:"username" json:"username"`
	HashedPassword string `bson:"hashedPassword" json:"hashedPassword"`
	Role           string `bson:"role" json:"role"`
}

type UserStore struct {
	database *mongo.Database
}

func NewUserStore(database *mongo.Database) *UserStore {
	return &UserStore{
		database: database,
	}
}

func (u *UserStore) CreateUser(ctx context.Context, user *User) (string, error) {
	collection := u.database.Collection("users")

	existing, _ := u.GetUserByUsername(ctx, user.Username)
	if existing != nil {
		return "", ErrUserExists
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", errors.New("failed to get inserted ID")
	}

	return oid.Hex(), nil
}

func (u *UserStore) GetUserByID(ctx context.Context, id string) (*User, error) {
	collection := u.database.Collection("users")

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var user User
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Id = oid.Hex()
	return &user, nil
}

func (u *UserStore) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	collection := u.database.Collection("users")

	var result struct {
		ID             bson.ObjectID `bson:"_id"`
		Username       string        `bson:"username"`
		HashedPassword string        `bson:"hashedPassword"`
		Role           string        `bson:"role"`
	}

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user := &User{
		Id:             result.ID.Hex(),
		Username:       result.Username,
		HashedPassword: result.HashedPassword,
		Role:           result.Role,
	}

	return user, nil
}

func (u *UserStore) DeleteUserByID(ctx context.Context, id string) error {
	collection := u.database.Collection("users")

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrUserNotFound
	}

	err = collection.FindOneAndDelete(ctx, bson.M{"_id": oid}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrUserNotFound
	}
	return err
}

func (u *UserStore) ListUsers(ctx context.Context, roleFilter string, usernameFilter string) ([]*User, error) {
	collection := u.database.Collection("users")

	filter := bson.M{}
	if roleFilter != "" {
		filter["role"] = roleFilter
	}
	if usernameFilter != "" {
		filter["username"] = bson.M{"$regex": usernameFilter, "$options": "i"}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID             bson.ObjectID `bson:"_id"`
		Username       string        `bson:"username"`
		HashedPassword string        `bson:"hashedPassword"`
		Role           string        `bson:"role"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	users := make([]*User, len(results))
	for i, result := range results {
		users[i] = &User{
			Id:             result.ID.Hex(),
			Username:       result.Username,
			HashedPassword: result.HashedPassword,
			Role:           result.Role,
		}
	}

	return users, nil
}
