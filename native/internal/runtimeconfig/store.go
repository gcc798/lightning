package runtimeconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gcc798/microservice-kit/internal/platform/redislock"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	cacheKeyPrefix = "microservice-kit:config:"
	lockKeyPrefix  = "microservice-kit:lock:config:"
)

// Source loads persistent configuration when Redis has no value.
type Source interface {
	Load(ctx context.Context, code string) ([]byte, error)
}

type gormSource struct{ db *gorm.DB }

// configRow is the small persistence projection needed by shared runtime config;
// SYS owns the full domain model and is intentionally not imported here.
type configRow struct {
	Data []byte `gorm:"column:data"`
}

func (configRow) TableName() string { return "s_config" }

// NewGormSource creates an s_config-backed source.
func NewGormSource(db *gorm.DB) Source { return &gormSource{db: db} }

func (s *gormSource) Load(ctx context.Context, code string) ([]byte, error) {
	var config configRow
	if err := s.db.WithContext(ctx).Where("code = ?", code).First(&config).Error; err != nil {
		return nil, err
	}
	return append([]byte(nil), config.Data...), nil
}

// Store is a shared Redis runtime configuration store with database fallback.
type Store struct {
	redis  *redis.Client
	source Source
	locker *redislock.Locker
}

func NewStore(redisClient *redis.Client, source Source, locker *redislock.Locker) *Store {
	return &Store{redis: redisClient, source: source, locker: locker}
}

// Get loads and decodes one configuration code.
func (s *Store) Get(ctx context.Context, code string, target any) error {
	raw, err := s.GetRaw(ctx, code)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("decode runtime configuration %q: %w", code, err)
	}
	return nil
}

// GetRaw returns Redis data, loading s_config under a distributed lock on miss.
func (s *Store) GetRaw(ctx context.Context, code string) ([]byte, error) {
	if code == "" {
		return nil, errors.New("configuration code is empty")
	}
	raw, err := s.redis.Get(ctx, CacheKey(code)).Bytes()
	if err == nil {
		if validateErr := Validate(code, raw); validateErr == nil {
			return raw, nil
		}
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("read runtime configuration %q from Redis: %w", code, err)
	}
	var loaded []byte
	if err := s.WithCodeLock(ctx, code, func() error {
		cached, err := s.redis.Get(ctx, CacheKey(code)).Bytes()
		if err == nil {
			if validateErr := Validate(code, cached); validateErr == nil {
				loaded = cached
				return nil
			}
			if err := s.redis.Del(ctx, CacheKey(code)).Err(); err != nil {
				return fmt.Errorf("delete invalid runtime configuration %q from Redis: %w", code, err)
			}
		}
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("recheck runtime configuration %q in Redis: %w", code, err)
		}
		persistent, err := s.source.Load(ctx, code)
		if err != nil {
			return fmt.Errorf("load runtime configuration %q from s_config: %w", code, err)
		}
		if err := Validate(code, persistent); err != nil {
			return fmt.Errorf("validate runtime configuration %q: %w", code, err)
		}
		if err := s.redis.Set(ctx, CacheKey(code), persistent, 0).Err(); err != nil {
			return fmt.Errorf("cache runtime configuration %q: %w", code, err)
		}
		loaded = persistent
		return nil
	}); err != nil {
		return nil, err
	}
	return loaded, nil
}

// WithCodeLock serializes all cache load and mutation for one code.
func (s *Store) WithCodeLock(ctx context.Context, code string, fn func() error) error {
	return s.locker.WithLock(ctx, LockKey(code), fn)
}

// DeleteCache removes the shared runtime value. Call it while holding the code lock.
func (s *Store) DeleteCache(ctx context.Context, code string) error {
	if err := s.redis.Del(ctx, CacheKey(code)).Err(); err != nil {
		return fmt.Errorf("delete runtime configuration cache %q: %w", code, err)
	}
	return nil
}

// SetCache writes the shared runtime value. Call it while holding the code lock.
func (s *Store) SetCache(ctx context.Context, code string, data []byte) error {
	if err := Validate(code, data); err != nil {
		return fmt.Errorf("validate runtime configuration %q: %w", code, err)
	}
	if err := s.redis.Set(ctx, CacheKey(code), data, 0).Err(); err != nil {
		return fmt.Errorf("write runtime configuration cache %q: %w", code, err)
	}
	return nil
}

func CacheKey(code string) string { return cacheKeyPrefix + code }
func LockKey(code string) string  { return lockKeyPrefix + code }
