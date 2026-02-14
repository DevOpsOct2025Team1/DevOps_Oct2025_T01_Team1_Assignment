package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestInitTelemetrySetsProviderAndPropagator(t *testing.T) {
	shutdown, err := InitTelemetry(context.Background(), Config{
		ServiceName: "test-service",
		Token:       "test-token",
		Dataset:     "test-dataset",
	})
	if err != nil {
		t.Fatalf("InitTelemetry returned error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("InitTelemetry returned nil shutdown")
	}
	defer shutdown(context.Background())

	if _, ok := otel.GetTracerProvider().(*trace.TracerProvider); !ok {
		t.Fatalf("unexpected tracer provider type: %T", otel.GetTracerProvider())
	}

	fields := otel.GetTextMapPropagator().Fields()
	if !containsAll(fields, "traceparent", "tracestate", "baggage") {
		t.Fatalf("expected traceparent, tracestate, baggage in propagator fields, got: %v", fields)
	}
}

func TestInitTelemetry_EmptyServiceName(t *testing.T) {
	_, err := InitTelemetry(context.Background(), Config{
		Token:   "tok",
		Dataset: "ds",
	})
	if err == nil || err.Error() != "telemetry: ServiceName is required" {
		t.Fatalf("expected ServiceName error, got %v", err)
	}
}

func TestInitTelemetry_EmptyToken(t *testing.T) {
	_, err := InitTelemetry(context.Background(), Config{
		ServiceName: "svc",
		Dataset:     "ds",
	})
	if err == nil || err.Error() != "telemetry: Token is required" {
		t.Fatalf("expected Token error, got %v", err)
	}
}

func TestInitTelemetry_EmptyDataset(t *testing.T) {
	_, err := InitTelemetry(context.Background(), Config{
		ServiceName: "svc",
		Token:       "tok",
	})
	if err == nil || err.Error() != "telemetry: Dataset is required" {
		t.Fatalf("expected Dataset error, got %v", err)
	}
}

func containsAll(fields []string, want ...string) bool {
	set := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		set[field] = struct{}{}
	}
	for _, key := range want {
		if _, ok := set[key]; !ok {
			return false
		}
	}
	return true
}
