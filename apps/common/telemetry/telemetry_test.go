package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestInitTracerSetsProviderAndPropagator(t *testing.T) {
	shutdown, err := InitTelemetry(context.Background(), Config{
		ServiceName: "test-service",
		Token:       "test-token",
		Dataset:     "test-dataset",
	})
	if err != nil {
		t.Fatalf("InitTracer returned error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("InitTracer returned nil shutdown")
	}

	if _, ok := otel.GetTracerProvider().(*trace.TracerProvider); !ok {
		t.Fatalf("unexpected tracer provider type: %T", otel.GetTracerProvider())
	}

	fields := otel.GetTextMapPropagator().Fields()
	if !containsAll(fields, "traceparent", "tracestate", "baggage") {
		t.Fatalf("expected traceparent, tracestate, baggage in propagator fields, got: %v", fields)
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
