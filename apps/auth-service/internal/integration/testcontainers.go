package integration

import (
	"context"
	"fmt"
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
		WaitingFor: wait.ForListeningPort("27017/tcp").WithStartupTimeout(90 * time.Second),
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

// Teardown stops and removes the MongoDB container.
func (mc *MongoContainer) Teardown(t *testing.T) {
	t.Helper()
	if mc == nil || mc.Container == nil {
		return
	}
	termCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := mc.Container.Terminate(termCtx); err != nil {
		t.Fatalf("failed to terminate mongo container: %v", err)
	}
}
