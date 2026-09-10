package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	// ClientCacheKeyPrefix 定义业务常量。
	ClientCacheKeyPrefix = "s_auth_client:"
	// ClientCacheTTL 定义业务常量。
	ClientCacheTTL = 30 * 24 * time.Hour
)

// ClientService 定义业务数据结构。
type ClientService interface {
	// AuthenticateClientID validates an active client and a login grant type.
	AuthenticateClientID(ctx context.Context, clientID, grantType string) (*model.AuthClient, error)
	// GetActiveClient returns an active client for non-login operations such as token refresh.
	GetActiveClient(ctx context.Context, clientID string) (*model.AuthClient, error)
}

type clientService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger logging.Logger
}

// NewClientService 创建组件实例。
func NewClientService(db *gorm.DB, redis *redis.Client, logger logging.Logger) ClientService {
	return &clientService{db: db, redis: redis, logger: logger}
}

// AuthenticateClientID 根据 clientId 认证客户端。
func (s *clientService) AuthenticateClientID(ctx context.Context, clientID, grantType string) (*model.AuthClient, error) {
	if grantType == "" {
		return nil, fmt.Errorf("grantType不能为空")
	}
	client, err := s.GetActiveClient(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if !client.IsGrantTypeSupported(grantType) {
		return nil, fmt.Errorf("客户端不支持该授权类型")
	}
	return client, nil
}

func (s *clientService) GetActiveClient(ctx context.Context, clientID string) (*model.AuthClient, error) {
	if clientID == "" {
		return nil, fmt.Errorf("clientId不能为空")
	}

	cacheKey := ClientCacheKeyPrefix + "id:" + clientID
	val, err := s.redis.Get(ctx, cacheKey).Result()
	var client *model.AuthClient
	if err == nil {
		var cached model.AuthClient
		if json.Unmarshal([]byte(val), &cached) == nil {
			client = &cached
		}
	}

	if client == nil {
		var m model.AuthClient
		c, err := m.FindByClientId(s.db, clientID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("客户端不存在")
			}
			return nil, fmt.Errorf("查询客户端失败: %w", err)
		}
		client = c
		b, _ := json.Marshal(client)
		_ = s.redis.Set(ctx, cacheKey, string(b), ClientCacheTTL).Err()
	}

	if !client.IsActive() {
		return nil, fmt.Errorf("客户端已停用")
	}
	return client, nil
}
