package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/gcc798/lightning/internal/platform/storage"
	"github.com/spf13/viper"
)

const AppEnvVar = "LIGHTNING_APP_ENV"

type Service string

const (
	ServiceGateway     Service = "gateway"
	ServiceIAM         Service = "iam"
	ServiceSystem      Service = "sys"
	ServiceResource    Service = "resource"
	ServiceScheduler   Service = "scheduler"
	ServiceUserManager Service = "usermgr"
)

type Server struct {
	Port        int    `mapstructure:"port"`
	TLSCertFile string `mapstructure:"tlsCertFile"`
	TLSKeyFile  string `mapstructure:"tlsKeyFile"`
}

type Registry struct {
	Driver  string `mapstructure:"driver"`
	Address string `mapstructure:"address"`
	Prefix  string `mapstructure:"prefix"`
}

type ServiceEndpoint struct {
	ID            string `mapstructure:"id"`
	AdvertiseHost string `mapstructure:"advertiseHost"`
}

type Gateway struct {
	RateLimitPerMinute int `mapstructure:"rateLimitPerMinute"`
}

type Database struct {
	DSN                    string `mapstructure:"dsn"`
	MaxOpenConns           int    `mapstructure:"maxOpenConns"`
	MaxIdleConns           int    `mapstructure:"maxIdleConns"`
	ConnMaxLifetimeMinutes int    `mapstructure:"connMaxLifetimeMinutes"`
	SlowThreshold          int    `mapstructure:"slowThreshold"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWT struct {
	Secret string `mapstructure:"secret"`
	Expire int64  `mapstructure:"expire"`
}

type Auth struct {
	TokenHeader     string `mapstructure:"tokenHeader"`
	AllowConcurrent bool   `mapstructure:"allowConcurrent"`
}

type CORS struct {
	Enabled bool `mapstructure:"enabled"`
}

type WebSocket struct {
	Enabled             bool `mapstructure:"enabled"`
	TimeoutEnabled      bool `mapstructure:"timeoutEnabled"`
	ReadTimeoutSeconds  int  `mapstructure:"readTimeoutSeconds"`
	WriteTimeoutSeconds int  `mapstructure:"writeTimeoutSeconds"`
	HeartbeatEnabled    bool `mapstructure:"heartbeatEnabled"`
	MaxReadTimeouts     int  `mapstructure:"maxReadTimeouts"`
}

type Config struct {
	AppDir    string          `mapstructure:"-"`
	Server    Server          `mapstructure:"server"`
	GRPC      Server          `mapstructure:"grpc"`
	Registry  Registry        `mapstructure:"registry"`
	Service   ServiceEndpoint `mapstructure:"service"`
	Database  Database        `mapstructure:"database"`
	Redis     Redis           `mapstructure:"redis"`
	JWT       JWT             `mapstructure:"jwt"`
	Auth      Auth            `mapstructure:"auth"`
	CORS      CORS            `mapstructure:"cors"`
	Storage   storage.Config  `mapstructure:"storage"`
	WebSocket WebSocket       `mapstructure:"websocket"`
	Gateway   Gateway         `mapstructure:"gateway"`
}

func Load(configDir string, service Service) (*Config, *viper.Viper, error) {
	profile := CurrentEnv()
	if profile != "dev" && profile != "prod" {
		return nil, nil, fmt.Errorf("%s must be dev or prod", AppEnvVar)
	}
	if service != ServiceGateway && service != ServiceIAM && service != ServiceSystem && service != ServiceResource && service != ServiceScheduler && service != ServiceUserManager {
		return nil, nil, fmt.Errorf("unknown config service %q", service)
	}
	v := viper.New()
	v.SetEnvPrefix("LIGHTNING")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := bindEnvironment(v); err != nil {
		return nil, nil, err
	}

	configFileName := fmt.Sprintf("conf.%s.yaml", profile)
	foundPath, err := ResolveFilePath(configDir, configFileName)
	if err != nil {
		return nil, nil, err
	}
	v.SetConfigFile(foundPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("read config from %s: %w", foundPath, err)
	}
	if err := requireExplicitConfiguration(v, service); err != nil {
		return nil, nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("decode config: %w", err)
	}
	cfg.AppDir = filepath.Dir(foundPath)
	if err := cfg.Validate(profile, service); err != nil {
		return nil, nil, err
	}
	return &cfg, v, nil
}

func requireExplicitConfiguration(v *viper.Viper, service Service) error {
	database := []string{"database.dsn", "database.maxOpenConns", "database.maxIdleConns", "database.connMaxLifetimeMinutes", "database.slowThreshold"}
	registry := []string{"registry.driver", "registry.address", "registry.prefix"}
	server := []string{"server.port"}
	grpc := []string{"grpc.port", "service.id", "service.advertiseHost"}
	auth := []string{"auth.tokenHeader", "cors.enabled"}

	var keys []string
	switch service {
	case ServiceGateway:
		keys = append(server, "server.tlsCertFile", "server.tlsKeyFile", "gateway.rateLimitPerMinute", "cors.enabled")
		keys = append(keys, registry...)
	case ServiceIAM:
		keys = append(keys, server...)
		keys = append(keys, grpc...)
		keys = append(keys, registry...)
		keys = append(keys, database...)
		keys = append(keys, "redis.addr", "redis.password", "redis.db", "jwt.secret", "jwt.expire")
		keys = append(keys, auth...)
		keys = append(keys, "auth.allowConcurrent")
		keys = append(keys, "websocket.enabled", "websocket.timeoutEnabled", "websocket.readTimeoutSeconds", "websocket.writeTimeoutSeconds", "websocket.heartbeatEnabled", "websocket.maxReadTimeouts")
	case ServiceSystem:
		keys = append(keys, server...)
		keys = append(keys, grpc...)
		keys = append(keys, registry...)
		keys = append(keys, database...)
		keys = append(keys, "redis.addr", "redis.password", "redis.db")
		keys = append(keys, auth...)
	case ServiceResource:
		keys = append(keys, server...)
		keys = append(keys, grpc...)
		keys = append(keys, registry...)
		keys = append(keys, database...)
		keys = append(keys, "storage.endpoint", "storage.accessKey", "storage.secretKey", "storage.region", "storage.bucket", "storage.useSSL")
		keys = append(keys, auth...)
	case ServiceScheduler:
		keys = append(keys, database...)
		keys = append(keys, registry...)
	case ServiceUserManager:
		keys = []string{"database.dsn"}
	}
	for _, key := range keys {
		_, environmentSet := os.LookupEnv(environmentName(key))
		if !v.InConfig(key) && !environmentSet {
			return fmt.Errorf("%s must be explicitly configured", key)
		}
	}
	return nil
}

func bindEnvironment(v *viper.Viper) error {
	keys := []string{
		"server.port", "server.tlsCertFile", "server.tlsKeyFile",
		"grpc.port",
		"registry.driver", "registry.address", "registry.prefix",
		"service.id", "service.advertiseHost",
		"gateway.rateLimitPerMinute",
		"database.dsn", "database.maxOpenConns", "database.maxIdleConns", "database.connMaxLifetimeMinutes", "database.slowThreshold",
		"redis.addr", "redis.password", "redis.db",
		"jwt.secret", "jwt.expire",
		"auth.tokenHeader", "auth.allowConcurrent",
		"cors.enabled",
		"storage.endpoint", "storage.accessKey", "storage.secretKey", "storage.region", "storage.bucket", "storage.useSSL",
		"websocket.enabled", "websocket.timeoutEnabled", "websocket.readTimeoutSeconds", "websocket.writeTimeoutSeconds", "websocket.heartbeatEnabled", "websocket.maxReadTimeouts",
	}
	for _, key := range keys {
		if err := v.BindEnv(key, environmentName(key)); err != nil {
			return fmt.Errorf("bind environment for %s: %w", key, err)
		}
	}
	return nil
}

func environmentName(key string) string {
	runes := []rune(key)
	var name strings.Builder
	name.WriteString("LIGHTNING_")
	for i, current := range runes {
		switch {
		case current == '.':
			name.WriteByte('_')
		case unicode.IsUpper(current):
			previousIsLower := i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1]))
			startsWord := i > 0 && unicode.IsUpper(runes[i-1]) && i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if previousIsLower || startsWord {
				name.WriteByte('_')
			}
			name.WriteRune(current)
		default:
			name.WriteRune(unicode.ToUpper(current))
		}
	}
	return name.String()
}

func (c *Config) Validate(profile string, service Service) error {
	if service != ServiceGateway && c.Database.DSN == "" {
		return fmt.Errorf("database.dsn is required")
	}
	if service == ServiceUserManager {
		return nil
	}
	if service != ServiceGateway {
		if c.Database.MaxOpenConns < 1 || c.Database.MaxIdleConns < 0 || c.Database.MaxIdleConns > c.Database.MaxOpenConns {
			return fmt.Errorf("database connection limits are invalid")
		}
		if c.Database.ConnMaxLifetimeMinutes < 1 || c.Database.SlowThreshold < 1 {
			return fmt.Errorf("database lifetimes and thresholds must be positive")
		}
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		if service != ServiceScheduler {
			return fmt.Errorf("server.port must be between 1 and 65535")
		}
	}
	if (service == ServiceIAM || service == ServiceSystem || service == ServiceResource) && (c.GRPC.Port < 1 || c.GRPC.Port > 65535) {
		return fmt.Errorf("grpc.port must be between 1 and 65535")
	}
	if c.Registry.Driver == "" || c.Registry.Address == "" || c.Registry.Prefix == "" {
		return fmt.Errorf("registry driver, address and prefix are required")
	}
	if service == ServiceScheduler {
		return nil
	}
	if service == ServiceGateway {
		if (c.Server.TLSCertFile == "") != (c.Server.TLSKeyFile == "") {
			return fmt.Errorf("both server.tlsCertFile and server.tlsKeyFile are required for TLS")
		}
		if c.Gateway.RateLimitPerMinute < 0 {
			return fmt.Errorf("gateway.rateLimitPerMinute cannot be negative")
		}
		return nil
	}
	if c.Service.AdvertiseHost == "" {
		return fmt.Errorf("service.advertiseHost is required")
	}
	if (service == ServiceIAM || service == ServiceSystem) && c.Redis.Addr == "" {
		return fmt.Errorf("redis.addr is required")
	}
	if (service == ServiceIAM || service == ServiceSystem) && c.Redis.DB < 0 {
		return fmt.Errorf("redis.db cannot be negative")
	}
	if c.Auth.TokenHeader == "" {
		return fmt.Errorf("auth.tokenHeader is required")
	}
	if service == ServiceIAM {
		if len(c.JWT.Secret) < 32 {
			return fmt.Errorf("jwt.secret must contain at least 32 characters")
		}
		if c.JWT.Expire < 1 {
			return fmt.Errorf("jwt.expire must be positive")
		}
		if c.WebSocket.ReadTimeoutSeconds < 1 || c.WebSocket.WriteTimeoutSeconds < 1 || c.WebSocket.MaxReadTimeouts < 1 {
			return fmt.Errorf("websocket timeouts and retry limit must be positive")
		}
	}
	if profile == "prod" && c.CORS.Enabled {
		return fmt.Errorf("cors must be disabled in prod")
	}
	if service == ServiceResource {
		if c.Storage.Region == "" {
			return fmt.Errorf("storage.region is required")
		}
		return c.Storage.Validate()
	}
	return nil
}

func CurrentEnv() string {
	return os.Getenv(AppEnvVar)
}

func ResolveFilePath(configDir, fileName string) (string, error) {
	if strings.TrimSpace(configDir) == "" {
		return "", errors.New("config directory is required")
	}
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable: %w", err)
	}

	var candidates []string
	if filepath.IsAbs(configDir) {
		candidates = []string{filepath.Join(configDir, fileName)}
	} else {
		candidates = []string{
			filepath.Join(workDir, configDir, fileName),
			filepath.Join(filepath.Dir(executable), fileName),
		}
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("config file not found: tried %v", candidates)
}
