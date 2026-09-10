package config

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/gcc798/microservice-kit/internal/platform/storage"
	"gopkg.in/yaml.v3"
)

func TestDevelopmentConfigStartsWithoutRequiredEnvironmentOverrides(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve config test path")
	}
	configPath := filepath.Join(filepath.Dir(sourceFile), "..", "..", "application", "iam", "conf.example.yaml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	appDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(appDir, "conf.dev.yaml"), content, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(AppEnvVar, "dev")
	t.Setenv("MS_K_DATABASE_DSN", "")
	t.Setenv("MS_K_REDIS_ADDR", "")
	t.Setenv("MS_K_REDIS_PASSWORD", "")
	t.Setenv("MS_K_JWT_SECRET", "")

	cfg, _, err := Load(appDir, ServiceIAM)
	if err != nil {
		t.Fatalf("Load(dev) error = %v", err)
	}
	if cfg.Database.DSN == "" || cfg.Redis.Addr == "" || cfg.JWT.Secret == "" {
		t.Fatal("development config is missing a required local value")
	}
	if cfg.AppDir != appDir {
		t.Fatalf("AppDir = %q, want %q", cfg.AppDir, appDir)
	}
	host, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if want := string(ServiceIAM) + "-" + host; cfg.Service.ID != want {
		t.Fatalf("Service.ID = %q, want %q", cfg.Service.ID, want)
	}
}

func TestServiceConfigsContainOnlyOwnedSections(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve config test path")
	}
	nativeDir := filepath.Join(filepath.Dir(sourceFile), "..", "..")
	services := []struct {
		name     string
		service  Service
		sections []string
	}{
		{name: "gateway", service: ServiceGateway, sections: []string{"cors", "gateway", "registry", "server", "service"}},
		{name: "iam", service: ServiceIAM, sections: []string{"auth", "cors", "database", "grpc", "jwt", "redis", "registry", "server", "service", "websocket"}},
		{name: "sys", service: ServiceSystem, sections: []string{"auth", "cors", "database", "grpc", "redis", "registry", "server", "service"}},
		{name: "resource", service: ServiceResource, sections: []string{"auth", "cors", "database", "grpc", "registry", "server", "service", "storage"}},
		{name: "scheduler", service: ServiceScheduler, sections: []string{"database", "registry", "service"}},
	}
	for _, tt := range services {
		t.Run(tt.name, func(t *testing.T) {
			for _, profile := range []string{"dev", "prod"} {
				t.Run(profile, func(t *testing.T) {
					examplePath := filepath.Join(nativeDir, "application", tt.name, "conf.example.yaml")
					data, err := os.ReadFile(examplePath)
					if err != nil {
						t.Fatal(err)
					}
					var document map[string]any
					if err := yaml.Unmarshal(data, &document); err != nil {
						t.Fatal(err)
					}
					sections := make([]string, 0, len(document))
					for section := range document {
						sections = append(sections, section)
					}
					slices.Sort(sections)
					if !slices.Equal(sections, tt.sections) {
						t.Fatalf("%s sections = %v, want %v", examplePath, sections, tt.sections)
					}
					dir := t.TempDir()
					if err := os.WriteFile(filepath.Join(dir, "conf."+profile+".yaml"), data, 0o600); err != nil {
						t.Fatal(err)
					}
					t.Setenv(AppEnvVar, profile)
					cfg, _, err := Load(dir, tt.service)
					if err != nil {
						t.Fatalf("Load(%s/%s) error = %v", tt.name, profile, err)
					}
					if tt.service != ServiceGateway && cfg.Database.DSN == "" {
						t.Fatalf("%s config is missing database.dsn", tt.name)
					}
					if tt.service == ServiceIAM && (cfg.Redis.Addr == "" || cfg.JWT.Secret == "") {
						t.Fatal("IAM config contains an empty required value")
					}
				})
			}
		})
	}
}

func TestLoadEnvironmentOverridesYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "conf.dev.yaml")
	content := []byte(`
server: { port: 9010 }
grpc: { port: 9110 }
registry: { driver: consul, address: "http://consul:8500", prefix: "/microservice-kit/services" }
service: { id: "", advertiseHost: "127.0.0.1" }
database: { dsn: yaml-dsn, maxOpenConns: 100, maxIdleConns: 10, connMaxLifetimeMinutes: 60, slowThreshold: 500 }
redis: { addr: "yaml-redis:6379", password: "", db: 0 }
jwt: { secret: yaml-secret-that-is-long-enough-for-tests, expire: 7200 }
auth: { tokenHeader: Authorization, allowConcurrent: false }
cors: { enabled: false }
websocket: { enabled: true, timeoutEnabled: true, readTimeoutSeconds: 60, writeTimeoutSeconds: 10, heartbeatEnabled: true, maxReadTimeouts: 3 }
`)
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv(AppEnvVar, "dev")
	t.Setenv("MS_K_SERVER_PORT", "8080")
	t.Setenv("MS_K_SERVICE_ID", "iam-explicit")
	t.Setenv("MS_K_DATABASE_DSN", "environment-dsn")
	t.Setenv("MS_K_REDIS_ADDR", "environment-redis:6379")
	t.Setenv("MS_K_JWT_SECRET", "environment-secret-with-at-least-32-characters")

	cfg, _, err := Load(dir, ServiceIAM)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("Server.Port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.Service.ID != "iam-explicit" {
		t.Fatalf("Service.ID = %q, want explicit environment value", cfg.Service.ID)
	}
	if cfg.Database.DSN != "environment-dsn" {
		t.Fatalf("Database.DSN = %q", cfg.Database.DSN)
	}
	if cfg.Redis.Addr != "environment-redis:6379" {
		t.Fatalf("Redis.Addr = %q", cfg.Redis.Addr)
	}
	if cfg.JWT.Secret != "environment-secret-with-at-least-32-characters" {
		t.Fatalf("JWT.Secret was not overridden by the environment")
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	valid := Config{
		Server:    Server{Port: 9009},
		GRPC:      Server{Port: 9100},
		Registry:  Registry{Driver: "consul", Address: "http://consul:8500", Prefix: "/microservice-kit/services"},
		Service:   ServiceEndpoint{AdvertiseHost: "127.0.0.1"},
		Database:  Database{DSN: "postgres-dsn", MaxOpenConns: 100, MaxIdleConns: 10, ConnMaxLifetimeMinutes: 60, SlowThreshold: 500},
		Redis:     Redis{Addr: "redis:6379"},
		JWT:       JWT{Secret: "a-secret-with-at-least-32-characters", Expire: 7200},
		Auth:      Auth{TokenHeader: "Authorization"},
		WebSocket: WebSocket{ReadTimeoutSeconds: 60, WriteTimeoutSeconds: 10, MaxReadTimeouts: 3},
		Storage: storage.Config{
			Endpoint:  "minio:9000",
			AccessKey: "access",
			SecretKey: "secret",
			Bucket:    "bucket",
			Region:    "us-east-1",
		},
	}
	if err := valid.Validate("dev", ServiceIAM); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	tweak := valid
	tweak.CORS.Enabled = true
	if err := tweak.Validate("prod", ServiceIAM); err == nil {
		t.Fatal("prod config with CORS enabled was accepted")
	}

	tweak = valid
	tweak.JWT.Secret = "short"
	if err := tweak.Validate("dev", ServiceIAM); err == nil {
		t.Fatal("short JWT secret was accepted")
	}

	scheduler := Config{
		Registry: Registry{Driver: "consul", Address: "http://consul:8500", Prefix: "/microservice-kit/services"},
		Database: Database{DSN: "postgres-dsn", MaxOpenConns: 20, MaxIdleConns: 5, ConnMaxLifetimeMinutes: 60, SlowThreshold: 500},
	}
	if err := scheduler.Validate("prod", ServiceScheduler); err != nil {
		t.Fatalf("minimal scheduler config rejected: %v", err)
	}
}

func TestCurrentEnvMustBeExplicit(t *testing.T) {
	t.Setenv(AppEnvVar, "")
	if got := CurrentEnv(); got != "" {
		t.Fatalf("CurrentEnv() = %q, want empty", got)
	}
	if _, _, err := Load(t.TempDir(), ServiceIAM); err == nil {
		t.Fatal("Load() accepted a missing MS_K_APP_ENV")
	}
}

func TestLoadRejectsImplicitConfiguration(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "conf.dev.yaml"), []byte("server:\n  port: 9009\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(AppEnvVar, "dev")
	_, _, err := Load(dir, ServiceGateway)
	if err == nil || !strings.Contains(err.Error(), "must be explicitly configured") {
		t.Fatalf("Load() error = %v, want explicit configuration error", err)
	}
}

func TestLoadRejectsUnknownProfile(t *testing.T) {
	t.Setenv(AppEnvVar, "staging")
	if _, _, err := Load(t.TempDir(), ServiceIAM); err == nil {
		t.Fatal("Load() accepted unsupported profile")
	}
}

func TestEnvironmentName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"database.maxOpenConns": "MS_K_DATABASE_MAX_OPEN_CONNS",
		"storage.useSSL":        "MS_K_STORAGE_USE_SSL",
	}
	for key, want := range tests {
		if got := environmentName(key); got != want {
			t.Errorf("environmentName(%q) = %q, want %q", key, got, want)
		}
	}
}
