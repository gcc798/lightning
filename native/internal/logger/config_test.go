package logger

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/yaml.v3"
)

func TestWithContextAddsTraceIdentifiers(t *testing.T) {
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceID, SpanID: spanID}))
	core, observed := observer.New(zap.InfoLevel)

	WithContext(ctx, &zapLogger{logger: zap.New(core)}).Info("correlated")
	fields := observed.All()[0].ContextMap()
	if fields["trace_id"] != traceID.String() || fields["span_id"] != spanID.String() {
		t.Fatalf("unexpected trace fields: %#v", fields)
	}
}

func TestServiceLoggerExamples(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve logger test path")
	}
	nativeDir := filepath.Join(filepath.Dir(sourceFile), "..", "..")
	for _, service := range []string{"gateway", "iam", "sys", "resource", "scheduler"} {
		path := filepath.Join(nativeDir, "application", service, "zaplogger.example.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var document map[string]any
		if err := yaml.Unmarshal(data, &document); err != nil {
			t.Fatal(err)
		}
		if _, exists := document["file"]; exists {
			t.Fatalf("%s configures unused file output", path)
		}
		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatalf("LoadConfig(%s) error = %v", path, err)
		}
		if cfg.Level == "" || cfg.Output != "console" || cfg.Encoding == "" {
			t.Fatalf("%s contains invalid logger configuration", path)
		}
	}
}
