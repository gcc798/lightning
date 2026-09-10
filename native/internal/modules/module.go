package modules

import (
	"context"

	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Module is a runtime capability that can be composed into any application entrypoint.
type Module interface {
	Name() string
	Init(ctx context.Context, cont Container) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Refresh(ctx context.Context, req ModuleRefreshRequest) error
}

// ModuleRefreshRequest describes an explicit local refresh. It intentionally has no
// version: each configuration code only retains its latest value.
type ModuleRefreshRequest struct {
	Codes  []string
	Reason string
}

// Container exposes only process-neutral dependencies to modules.
type Container interface {
	GetDB() *gorm.DB
	GetRedis() *redis.Client
	GetLogger() logging.Logger
	GetRuntimeConfig() *runtimeconfig.Store
	GetModule(name string) Module
}
