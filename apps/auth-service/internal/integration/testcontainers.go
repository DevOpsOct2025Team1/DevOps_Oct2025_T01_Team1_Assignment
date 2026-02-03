//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MongoContainer holds the running container and connection details.
type MongoContainer struct {
	Container tc.Container
	URI       string
}

// SetupMongoContainer starts a MongoDB docker container for tests and returns its URI.
// Note: This function registers t.Cleanup to terminate the container automatically.
func SetupMongoContainer(t *testing.T) *MongoContainer {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	req := tc.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE": "testdb",
		},
		// Avoid tailing logs; just wait until the port is listening.
		WaitingFor: wait.ForListeningPort("27017/tcp").WithStartupTimeout(2 * time.Minute),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start mongo container: %v", err)
	}

	// Ensure termination even if the test exits early.
	t.Cleanup(func() {
		termCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = container.Terminate(termCtx)
	})

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(context.Background())
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		_ = container.Terminate(context.Background())
		t.Fatalf("failed to get mapped port: %v", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%d", host, port.Int())

	return &MongoContainer{Container: container, URI: uri}
}

// waitForPortReady waits for a port to be ready.
func waitForPortReady(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		lastErr = err
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("port %s not ready: %v", addr, lastErr)
}
