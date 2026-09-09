package telemetry

import (
	"context"
	"errors"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// Init installs process-wide tracing. Without an OTLP endpoint spans still get
// valid IDs for log correlation, but are not exported.
func Init(ctx context.Context, serviceName, instanceID, environment string) (func(context.Context) error, error) {
	if strings.TrimSpace(serviceName) == "" {
		return nil, errors.New("telemetry service name is required")
	}
	if strings.TrimSpace(instanceID) == "" {
		return nil, errors.New("telemetry service instance ID is required")
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceInstanceID(instanceID),
		semconv.DeploymentEnvironmentName(environment),
	)
	options := []sdktrace.TracerProviderOption{sdktrace.WithResource(res)}
	if exporterConfigured() {
		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, err
		}
		options = append(options, sdktrace.WithBatcher(exporter))
	}

	provider := sdktrace.NewTracerProvider(options...)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return provider.Shutdown, nil
}

func exporterConfigured() bool {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("OTEL_TRACES_EXPORTER")), "none") {
		return false
	}
	return os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" || os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != ""
}

// TraceHTTPPath filters high-frequency infrastructure probes from tracing.
func TraceHTTPPath(path string) bool {
	switch path {
	case "/health", "/health/live", "/health/ready", "/health/startup", "/metrics":
		return false
	default:
		return true
	}
}
