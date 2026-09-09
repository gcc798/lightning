package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestInitWithoutExporterStillCreatesTraceIDs(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	previousProvider := otel.GetTracerProvider()
	previousPropagator := otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTracerProvider(previousProvider)
		otel.SetTextMapPropagator(previousPropagator)
	})

	shutdown, err := Init(context.Background(), "iam", "iam-test", "dev")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = shutdown(context.Background()) })

	_, span := otel.Tracer("test").Start(context.Background(), "login")
	defer span.End()
	if !span.SpanContext().IsValid() {
		t.Fatal("span has no valid trace_id/span_id")
	}
}

func TestInitRequiresInstanceID(t *testing.T) {
	if _, err := Init(context.Background(), "iam", "", "dev"); err == nil {
		t.Fatal("Init accepted an empty instance ID")
	}
}

func TestTraceHTTPPath(t *testing.T) {
	if TraceHTTPPath("/health/live") || !TraceHTTPPath("/login") {
		t.Fatal("unexpected HTTP tracing filter result")
	}
}
