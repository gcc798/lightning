package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gopkg.in/yaml.v3"
)

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
