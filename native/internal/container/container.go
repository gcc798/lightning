package container

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/database"
	"github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/modules"
	"github.com/gcc798/lightning/internal/platform/jwt"
	redisclient "github.com/gcc798/lightning/internal/platform/redis"
	"github.com/gcc798/lightning/internal/platform/redislock"
	"github.com/gcc798/lightning/internal/platform/storage"
	"github.com/gcc798/lightning/internal/runtimeconfig"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/extra/redisotel/v9"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Container 依赖注入容器接口
type Container interface {
	modules.Container
	GetConfig() *config.Config
	GetViper() *viper.Viper
	GetJWT() *jwt.Jwt
	GetStorage() storage.Storage
	RegisterModules(ctx context.Context, candidates ...modules.Module) error
	StartModules(ctx context.Context) error
	StopModules(ctx context.Context) error
	RefreshModule(ctx context.Context, name string, req modules.ModuleRefreshRequest) error
}

type container struct {
	config        *config.Config
	viper         *viper.Viper
	db            *gorm.DB
	redis         *goredis.Client
	jwt           *jwt.Jwt
	logger        logger.Logger
	storage       storage.Storage
	runtimeConfig *runtimeconfig.Store

	moduleMu          sync.RWMutex
	moduleLifecycleMu sync.Mutex
	modules           map[string]modules.Module
	moduleOrder       []string
	moduleRefreshMu   map[string]*sync.Mutex
	startedModules    int
}

// Option composes process-specific infrastructure without coupling modules to cmd entries.
type Option func(*container) error

func WithIAMInfrastructure() Option {
	return func(c *container) error {
		if err := c.initRedis(); err != nil {
			return err
		}
		c.initRuntimeConfig()
		c.initJWT()
		return nil
	}
}

func WithSystemInfrastructure() Option {
	return func(c *container) error {
		if err := c.initRedis(); err != nil {
			return err
		}
		c.initRuntimeConfig()
		return nil
	}
}

func WithResourceInfrastructure() Option {
	return func(c *container) error { return c.initStorage() }
}

// NewEmpty 创建一个空容器，调用方可以按需初始化指定组件。
func NewEmpty(cfg *config.Config, v *viper.Viper, log logger.Logger) *container {
	return &container{
		config:          cfg,
		viper:           v,
		logger:          log,
		modules:         make(map[string]modules.Module),
		moduleOrder:     make([]string, 0),
		moduleRefreshMu: make(map[string]*sync.Mutex),
	}
}

// New 创建新的容器实例
func New(cfg *config.Config, v *viper.Viper, log logger.Logger, options ...Option) (Container, error) {
	c := NewEmpty(cfg, v, log)

	// 1. 初始化基础组件
	if err := c.initDB(); err != nil {
		return nil, err
	}
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// InitDBOnly 仅初始化数据库连接，适合一次性工具复用现有连库逻辑。
func (c *container) InitDBOnly() error {
	return c.initDB()
}

// initDB 初始化数据库
func (c *container) initDB() error {
	dsn := c.config.Database.DSN
	slowThreshold := time.Second
	if c.config.Database.SlowThreshold > 0 {
		slowThreshold = time.Duration(c.config.Database.SlowThreshold) * time.Millisecond
	}
	gormLogger := gormlogger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             slowThreshold,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	sqlDB, err := otelsql.Open("pgx", dsn, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))
	if err != nil {
		return fmt.Errorf("instrument database connection: %w", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		_ = sqlDB.Close()
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(c.config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.config.Database.ConnMaxLifetimeMinutes) * time.Minute)

	// 注册 ID 生成插件
	idGenPlugin := &database.IDGenPlugin{}
	if err := db.Use(idGenPlugin); err != nil {
		c.logger.Warn("failed to register ID generation plugin", zap.Error(err))
	}

	c.db = db
	return nil
}

// initRedis 初始化Redis
func (c *container) initRedis() error {
	redisClient := redisclient.NewRedis(c.config.Redis.Addr, c.config.Redis.Password, c.config.Redis.DB)
	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return fmt.Errorf("instrument redis tracing: %w", err)
	}
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	c.redis = redisClient
	return nil
}

// initJWT 初始化JWT
func (c *container) initJWT() {
	c.jwt = jwt.New(c.config.JWT.Secret, int64(c.config.JWT.Expire))
}

func (c *container) initRuntimeConfig() {
	locker := redislock.New(c.redis)
	c.runtimeConfig = runtimeconfig.NewStore(c.redis, runtimeconfig.NewGormSource(c.db), locker)
}

func (c *container) initStorage() error {
	store, err := storage.New(c.config.Storage)
	if err != nil {
		return fmt.Errorf("initialize storage: %w", err)
	}
	c.storage = store
	c.logger.Info("storage initialized", zap.String("bucket", c.config.Storage.Bucket))
	return nil
}

func (c *container) GetConfig() *config.Config              { return c.config }
func (c *container) GetViper() *viper.Viper                 { return c.viper }
func (c *container) GetDB() *gorm.DB                        { return c.db }
func (c *container) GetRedis() *goredis.Client              { return c.redis }
func (c *container) GetRuntimeConfig() *runtimeconfig.Store { return c.runtimeConfig }
func (c *container) GetJWT() *jwt.Jwt                       { return c.jwt }
func (c *container) GetLogger() logger.Logger               { return c.logger }
func (c *container) GetStorage() storage.Storage            { return c.storage }

// RegisterModules initializes and registers modules in deterministic order.
func (c *container) RegisterModules(ctx context.Context, candidates ...modules.Module) error {
	c.moduleLifecycleMu.Lock()
	defer c.moduleLifecycleMu.Unlock()

	c.moduleMu.Lock()
	if c.startedModules != 0 {
		c.moduleMu.Unlock()
		return errors.New("cannot register modules after module startup")
	}
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate == nil {
			c.moduleMu.Unlock()
			return errors.New("cannot register a nil module")
		}
		name := strings.TrimSpace(candidate.Name())
		if name == "" {
			c.moduleMu.Unlock()
			return errors.New("cannot register a module with an empty name")
		}
		if _, exists := c.modules[name]; exists {
			c.moduleMu.Unlock()
			return fmt.Errorf("module %q is already registered", name)
		}
		if _, exists := seen[name]; exists {
			c.moduleMu.Unlock()
			return fmt.Errorf("module %q appears more than once", name)
		}
		seen[name] = struct{}{}
	}
	for _, candidate := range candidates {
		name := strings.TrimSpace(candidate.Name())
		c.modules[name] = candidate
		c.moduleOrder = append(c.moduleOrder, name)
		c.moduleRefreshMu[name] = &sync.Mutex{}
	}
	c.moduleMu.Unlock()

	initialized := make([]modules.Module, 0, len(candidates))
	for _, candidate := range candidates {
		if err := candidate.Init(ctx, c); err != nil {
			for i := len(initialized) - 1; i >= 0; i-- {
				_ = initialized[i].Stop(ctx)
			}
			c.moduleMu.Lock()
			for _, rollback := range candidates {
				name := strings.TrimSpace(rollback.Name())
				delete(c.modules, name)
				delete(c.moduleRefreshMu, name)
			}
			c.moduleOrder = c.moduleOrder[:len(c.moduleOrder)-len(candidates)]
			c.moduleMu.Unlock()
			return fmt.Errorf("initialize module %q: %w", candidate.Name(), err)
		}
		initialized = append(initialized, candidate)
	}
	return nil
}

// GetModule returns a stable registered module by name.
func (c *container) GetModule(name string) modules.Module {
	c.moduleMu.RLock()
	defer c.moduleMu.RUnlock()
	return c.modules[name]
}

// StartModules starts registered modules and rolls back on partial failure.
func (c *container) StartModules(ctx context.Context) error {
	c.moduleLifecycleMu.Lock()
	defer c.moduleLifecycleMu.Unlock()

	c.moduleMu.RLock()
	if c.startedModules != 0 {
		c.moduleMu.RUnlock()
		return errors.New("modules are already started")
	}
	ordered := make([]modules.Module, 0, len(c.moduleOrder))
	for _, name := range c.moduleOrder {
		ordered = append(ordered, c.modules[name])
	}
	c.moduleMu.RUnlock()

	started := make([]modules.Module, 0, len(ordered))
	for _, candidate := range ordered {
		if err := candidate.Start(ctx); err != nil {
			var rollbackErrors []error
			for i := len(started) - 1; i >= 0; i-- {
				if stopErr := started[i].Stop(ctx); stopErr != nil {
					rollbackErrors = append(rollbackErrors, fmt.Errorf("rollback module %q: %w", started[i].Name(), stopErr))
				}
			}
			return errors.Join(append([]error{fmt.Errorf("start module %q: %w", candidate.Name(), err)}, rollbackErrors...)...)
		}
		started = append(started, candidate)
	}
	c.moduleMu.Lock()
	c.startedModules = len(started)
	c.moduleMu.Unlock()
	return nil
}

// StopModules stops all started modules in reverse registration order.
func (c *container) StopModules(ctx context.Context) error {
	c.moduleLifecycleMu.Lock()
	defer c.moduleLifecycleMu.Unlock()

	c.moduleMu.RLock()
	count := c.startedModules
	ordered := make([]modules.Module, 0, count)
	for i := 0; i < count; i++ {
		ordered = append(ordered, c.modules[c.moduleOrder[i]])
	}
	c.moduleMu.RUnlock()

	var stopErrors []error
	for i := len(ordered) - 1; i >= 0; i-- {
		if err := ordered[i].Stop(ctx); err != nil {
			stopErrors = append(stopErrors, fmt.Errorf("stop module %q: %w", ordered[i].Name(), err))
		}
	}
	c.moduleMu.Lock()
	c.startedModules = 0
	c.moduleMu.Unlock()
	return errors.Join(stopErrors...)
}

// RefreshModule serializes explicit refreshes for one module.
func (c *container) RefreshModule(ctx context.Context, name string, req modules.ModuleRefreshRequest) error {
	c.moduleMu.RLock()
	candidate := c.modules[name]
	refreshMu := c.moduleRefreshMu[name]
	c.moduleMu.RUnlock()
	if candidate == nil || refreshMu == nil {
		return fmt.Errorf("module %q is not registered", name)
	}
	refreshMu.Lock()
	defer refreshMu.Unlock()
	if err := candidate.Refresh(ctx, req); err != nil {
		return fmt.Errorf("refresh module %q: %w", name, err)
	}
	return nil
}
