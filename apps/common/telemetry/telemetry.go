package telemetry

import (
	"context"
	"crypto/tls"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Config holds the configuration for initializing the tracer.
type Config struct {
	ServiceName string
	Environment string
	Endpoint    string
	Token       string
	Dataset     string
}

// InitTracer initializes OpenTelemetry with Axiom as the backend.
// Returns a shutdown function that should be deferred in main().
func InitTracer(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "us-east-1.aws.edge.axiom.co"
	}

	environment := cfg.Environment
	if environment == "" {
		environment = "development"
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization":   "Bearer " + cfg.Token,
			"X-AXIOM-DATASET": cfg.Dataset,
		}),
		otlptracehttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", environment),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
