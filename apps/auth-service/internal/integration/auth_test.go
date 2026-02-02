package integration

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/client"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/jwt"
	authsvc "github.com/provsalt/DOP_P01_Team1/auth-service/internal/service"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
)

// startUserService starts an in-process gRPC user service bound to a random port.
func startUserService(t *testing.T, dbURI string) (addr string, stop func()) {
	t.Helper()

	port := 18080 + rand.Intn(1000)
	address := fmt.Sprintf("localhost:%d", port)

	cmd := exec.Command("go", "run", "github.com/provsalt/DOP_P01_Team1/user-service/cmd/server")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		fmt.Sprintf("MONGODB_URI=%s", dbURI),
		"MONGODB_DATABASE=user_service_test",
		"AXIOM_API_TOKEN=",
		"ENVIRONMENT=test",
	)

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start user-service: %v", err)
	}

	// Give service time to boot
	time.Sleep(2 * time.Second)

	stopFn := func() {
		if cmd.Process != nil {
			_ = cmd.Process.Signal(os.Interrupt)
			time.Sleep(500 * time.Millisecond)
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait() 
		}
	}

	return address, stopFn
}


func TestAuthService_SignUpLoginValidate(t *testing.T) {
	// 1) Start MongoDB testcontainer
	mongoC := SetupMongoContainer(t)
	defer mongoC.Teardown(t)

	// 2) Start in-process user-service backed by the container DB
	userAddr, stopUser := startUserService(t, mongoC.URI)
	defer stopUser()

	// 3) Create auth-service dependencies
	userClient, err := client.NewUserServiceClient(userAddr)
	if err != nil {
		t.Fatalf("failed to create user client: %v", err)
	}
	defer userClient.Close()

	jwtManager := jwt.NewJWTManager("testsecret", 24*time.Hour)
	authServer := authsvc.NewAuthServiceServer(userClient, jwtManager)

	ctx := context.Background()

	// 4) SignUp new user
	signupResp, err := authServer.SignUp(ctx, &authv1.SignUpRequest{
		Username: "alice",
		Password: "s3cret",
	})
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}
	if signupResp.User == nil || signupResp.User.Id == "" {
		t.Fatalf("SignUp returned invalid user: %+v", signupResp.User)
	}
	if signupResp.Token == "" {
		t.Fatalf("SignUp returned empty token")
	}

	// 5) Login with same credentials
	loginResp, err := authServer.Login(ctx, &authv1.LoginRequest{
		Username: "alice",
		Password: "s3cret",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginResp.User == nil || loginResp.User.Id == "" {
		t.Fatalf("Login returned invalid user: %+v", loginResp.User)
	}
	if loginResp.Token == "" {
		t.Fatalf("Login returned empty token")
	}

	// 6) Validate token
	validateResp, err := authServer.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		Token: loginResp.Token,
	})
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if !validateResp.Valid {
		t.Fatalf("expected token to be valid")
	}
	if validateResp.User == nil || validateResp.User.Username != "alice" {
		t.Fatalf("unexpected validated user: %+v", validateResp.User)
	}
}
