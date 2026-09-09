package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gcc798/lightning/application/sys-rpc/pkg/storage"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// StorageManager manages lazy-initialized Storage instances keyed by storage env ID.
type StorageManager struct {
	db       sqlx.SqlConn
	storages map[int64]storage.Storage
	mu       sync.RWMutex
}

// NewStorageManager creates a StorageManager that reads storage env config from DB.
func NewStorageManager(db sqlx.SqlConn) *StorageManager {
	return &StorageManager{
		db:       db,
		storages: make(map[int64]storage.Storage),
	}
}

// GetStorage returns the Storage instance for the given env ID, creating it lazily.
func (m *StorageManager) GetStorage(ctx context.Context, envID int64) (storage.Storage, error) {
	m.mu.RLock()
	s, ok := m.storages[envID]
	m.mu.RUnlock()
	if ok {
		return s, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	// Double-check after acquiring write lock
	if s, ok = m.storages[envID]; ok {
		return s, nil
	}

	var row struct {
		StorageType string `db:"storage_type"`
		Config      string `db:"config"`
	}
	err := m.db.QueryRowCtx(ctx, &row,
		`select storage_type, config from public.s_storage_env where id = $1 and status = 0 limit 1`,
		envID,
	)
	if err != nil {
		return nil, fmt.Errorf("storage env %d not found: %w", envID, err)
	}

	s, err = m.createStorage(row.StorageType, row.Config)
	if err != nil {
		return nil, err
	}
	m.storages[envID] = s
	return s, nil
}

// GetStorageByCode resolves an env by code and returns its Storage instance.
func (m *StorageManager) GetStorageByCode(ctx context.Context, code string) (storage.Storage, error) {
	var envID int64
	err := m.db.QueryRowCtx(ctx, &envID,
		`select id from public.s_storage_env where code = $1 and status = 0 limit 1`,
		code,
	)
	if err != nil {
		return nil, fmt.Errorf("storage env %q not found: %w", code, err)
	}
	return m.GetStorage(ctx, envID)
}

// GetDefaultStorage returns the Storage instance for the default env.
func (m *StorageManager) GetDefaultStorage(ctx context.Context) (storage.Storage, error) {
	var envID int64
	err := m.db.QueryRowCtx(ctx, &envID,
		`select id from public.s_storage_env where is_default = true and status = 0 order by id desc limit 1`,
	)
	if err != nil {
		return nil, fmt.Errorf("no default storage env: %w", err)
	}
	return m.GetStorage(ctx, envID)
}

func (m *StorageManager) createStorage(storageType, configJSON string) (storage.Storage, error) {
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("parse storage config failed: %w", err)
	}
	switch storageType {
	case "local":
		return storage.NewLocalStorage(storage.LocalConfig{
			BasePath:  getString(cfg, "basePath", "runtime/attachments"),
			URLPrefix: getString(cfg, "urlPrefix", "http://localhost:9002/files"),
		})
	case "s3", "minio", "oss":
		return storage.NewS3Storage(storage.S3Config{
			Endpoint:  getString(cfg, "endpoint", "localhost:9000"),
			AccessKey: getString(cfg, "accessKey", ""),
			SecretKey: getString(cfg, "secretKey", ""),
			Bucket:    getString(cfg, "bucket", "lightning"),
			Region:    getString(cfg, "region", "us-east-1"),
			UseSSL:    getBool(cfg, "useSSL", false),
		})
	default:
		// Fall back to local storage
		return storage.NewLocalStorage(storage.LocalConfig{
			BasePath:  getString(cfg, "basePath", "runtime/attachments"),
			URLPrefix: getString(cfg, "urlPrefix", "http://localhost:9002/files"),
		})
	}
}

func getString(m map[string]interface{}, key, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getBool(m map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}
