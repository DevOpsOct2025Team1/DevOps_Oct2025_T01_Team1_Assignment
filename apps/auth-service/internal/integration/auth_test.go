//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/client"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/jwt"
	authsvc "github.com/provsalt/DOP_P01_Team1/auth-service/internal/service"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
)

// startUserService starts an in-process gRPC user service bound to an OS-assigned free port.
func startUserService(t *testing.T, dbURI string) (addr string, stop func()) {
	t.Helper()

	// Ask OS for a free port, then release it for the child process to bind.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to acquire free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	address := fmt.Sprintf("127.0.0.1:%d", port)
	_ = l.Close()

	// Run the local user-service entrypoint to avoid module path/network resolution.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", "..", ".."))
	cmd := exec.Command("go", "run", "./cmd/server")
	cmd.Dir = filepath.Join(repoRoot, "apps", "user-service")
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

	// Wait for the gRPC port to be ready instead of a fixed sleep
	if err := waitForPortReady(address, 15*time.Second); err != nil {
		stopFn := func() {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(os.Interrupt)
				time.Sleep(500 * time.Millisecond)
				_ = cmd.Process.Kill()
				_, _ = cmd.Process.Wait()
			}
		}
		stopFn()
		t.Fatalf("user-service didn't become ready on %s: %v", address, err)
	}

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
	// Rely on t.Cleanup in SetupMongoContainer; no manual teardown.

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
