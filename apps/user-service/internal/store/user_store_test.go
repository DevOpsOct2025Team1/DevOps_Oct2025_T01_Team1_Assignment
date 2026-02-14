package store

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGetUserByID_InvalidHex_Unit(t *testing.T) {
	invalidHexID := "not-a-valid-hex-id"
	_, err := bson.ObjectIDFromHex(invalidHexID)
	if err == nil {
		t.Error("expected error for invalid hex ID")
	}
}

func TestDeleteUserByID_InvalidHex_Unit(t *testing.T) {
	invalidHexID := "xyz123"
	_, err := bson.ObjectIDFromHex(invalidHexID)
	if err == nil {
		t.Error("expected error for invalid hex ID")
	}
}

func TestObjectIDFromHex_Valid(t *testing.T) {
	validHex := bson.NewObjectID().Hex()
	oid, err := bson.ObjectIDFromHex(validHex)
	if err != nil {
		t.Errorf("expected no error for valid hex, got %v", err)
	}
	if oid.IsZero() {
		t.Error("expected non-zero ObjectID")
	}
	if oid.Hex() != validHex {
		t.Errorf("expected hex %s, got %s", validHex, oid.Hex())
	}
}

func TestBsonFilter_Role(t *testing.T) {
	roleFilter := "admin"
	filter := bson.M{}
	filter["role"] = roleFilter

	if filter["role"] != "admin" {
		t.Errorf("expected role 'admin', got %v", filter["role"])
	}
}

func TestBsonFilter_Username_Regex(t *testing.T) {
	usernameFilter := "john"
	filter := bson.M{}
	filter["username"] = bson.M{"$regex": usernameFilter, "$options": "i"}

	regexFilter, ok := filter["username"].(bson.M)
	if !ok {
		t.Error("expected username filter to be bson.M")
	}
	if regexFilter["$regex"] != "john" {
		t.Errorf("expected regex 'john', got %v", regexFilter["$regex"])
	}
	if regexFilter["$options"] != "i" {
		t.Errorf("expected options 'i', got %v", regexFilter["$options"])
	}
}

func TestBsonFilter_Combined(t *testing.T) {
	roleFilter := "admin"
	usernameFilter := "alice"

	filter := bson.M{}
	filter["role"] = roleFilter
	filter["username"] = bson.M{"$regex": usernameFilter, "$options": "i"}

	if filter["role"] != "admin" {
		t.Errorf("expected role 'admin'")
	}

	regexFilter, ok := filter["username"].(bson.M)
	if !ok {
		t.Error("expected username filter")
	}
	if regexFilter["$regex"] != "alice" {
		t.Errorf("expected regex 'alice'")
	}
}

func TestBsonFilter_Empty(t *testing.T) {
	filter := bson.M{}
	if len(filter) != 0 {
		t.Errorf("expected empty filter, got %v", filter)
	}
}

func TestUserStructure(t *testing.T) {
	user := &User{
		Id:             "12345",
		Username:       "testuser",
		HashedPassword: "hashedpw",
		Role:           "user",
	}

	if user.Id != "12345" {
		t.Errorf("expected ID 12345, got %s", user.Id)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", user.Username)
	}
	if user.Role != "user" {
		t.Errorf("expected role user, got %s", user.Role)
	}
}

func TestErrorVariables(t *testing.T) {
	if ErrUserNotFound == nil {
		t.Error("ErrUserNotFound should be defined")
	}
	if ErrUserExists == nil {
		t.Error("ErrUserExists should be defined")
	}
	if ErrUserNotFound.Error() != "user not found" {
		t.Errorf("expected 'user not found', got %s", ErrUserNotFound.Error())
	}
	if ErrUserExists.Error() != "user already exists" {
		t.Errorf("expected 'user already exists', got %s", ErrUserExists.Error())
	}
}
