package data

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	entpkg "github.com/gcc798/lightning/kratos/application/sys-rpc/ent"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/ent/storageenv"
)

type storageObjectInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
}

type storageBackend interface {
	Type() string
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string, expires time.Duration) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Stat(ctx context.Context, key string) (*storageObjectInfo, error)
	Ping(ctx context.Context) error
}

type storageFactory func(env *entpkg.StorageEnv) (storageBackend, error)

type StorageManager struct {
	ent       *entpkg.Client
	mu        sync.RWMutex
	factories map[string]storageFactory
	envs      map[int64]*entpkg.StorageEnv
	codeToID  map[string]int64
	backends  map[int64]storageBackend
	defaultID int64
}

func NewStorageManager(entClient *entpkg.Client) *StorageManager {
	mgr := &StorageManager{
		ent:       entClient,
		factories: make(map[string]storageFactory),
		envs:      make(map[int64]*entpkg.StorageEnv),
		codeToID:  make(map[string]int64),
		backends:  make(map[int64]storageBackend),
	}
	mgr.RegisterFactory("local", newLocalStorageBackend)
	mgr.RegisterFactory("minio", newS3CompatibleStorageBackend("minio"))
	mgr.RegisterFactory("s3", newS3CompatibleStorageBackend("s3"))
	mgr.RegisterFactory("oss", newS3CompatibleStorageBackend("oss"))
	return mgr
}

func (m *StorageManager) RegisterFactory(storageType string, factory storageFactory) {
	if m == nil || factory == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[strings.ToLower(strings.TrimSpace(storageType))] = factory
}

func (m *StorageManager) InvalidateAll() {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.envs = make(map[int64]*entpkg.StorageEnv)
	m.codeToID = make(map[string]int64)
	m.backends = make(map[int64]storageBackend)
	m.defaultID = 0
}

func (m *StorageManager) InvalidateEnv(envID int64) {
	if m == nil || envID <= 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if env, ok := m.envs[envID]; ok {
		delete(m.codeToID, normalizeStorageCode(env.Code))
		if env.IsDefault {
			m.defaultID = 0
		}
	}
	delete(m.envs, envID)
	delete(m.backends, envID)
}

func (m *StorageManager) ResolveActiveEnv(ctx context.Context, envCode string) (*entpkg.StorageEnv, error) {
	if m == nil || m.ent == nil {
		return nil, errors.New("storage manager is not initialized")
	}
	envCode = normalizeStorageCode(envCode)
	if envCode == "" {
		return m.defaultEnv(ctx)
	}

	m.mu.RLock()
	if envID, ok := m.codeToID[envCode]; ok {
		if env, ok := m.envs[envID]; ok {
			m.mu.RUnlock()
			return env, nil
		}
	}
	m.mu.RUnlock()

	env, err := m.ent.StorageEnv.Query().
		Where(
			storageenv.Code(envCode),
			storageenv.Status(0),
			storageenv.DeletedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, errors.New("storage env not found")
		}
		return nil, err
	}
	m.cacheEnv(env)
	return env, nil
}

func (m *StorageManager) ActiveEnvByID(ctx context.Context, envID int64) (*entpkg.StorageEnv, error) {
	if m == nil || m.ent == nil {
		return nil, errors.New("storage manager is not initialized")
	}
	if envID <= 0 {
		return nil, errors.New("storage env not found")
	}

	m.mu.RLock()
	if env, ok := m.envs[envID]; ok {
		m.mu.RUnlock()
		return env, nil
	}
	m.mu.RUnlock()

	env, err := m.ent.StorageEnv.Query().
		Where(
			storageenv.ID(envID),
			storageenv.Status(0),
			storageenv.DeletedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, errors.New("storage env not found")
		}
		return nil, err
	}
	m.cacheEnv(env)
	return env, nil
}

func (m *StorageManager) Upload(ctx context.Context, envCode, key string, reader io.Reader, size int64, contentType string) (*entpkg.StorageEnv, error) {
	env, backend, err := m.backendByCode(ctx, envCode)
	if err != nil {
		return nil, err
	}
	if err := backend.Upload(ctx, key, reader, size, contentType); err != nil {
		return nil, err
	}
	return env, nil
}

func (m *StorageManager) Download(ctx context.Context, envID int64, key string) (io.ReadCloser, error) {
	_, backend, err := m.backendByID(ctx, envID)
	if err != nil {
		return nil, err
	}
	return backend.Download(ctx, key)
}

func (m *StorageManager) Delete(ctx context.Context, envID int64, key string) error {
	_, backend, err := m.backendByID(ctx, envID)
	if err != nil {
		return err
	}
	return backend.Delete(ctx, key)
}

func (m *StorageManager) AttachmentURL(ctx context.Context, item *entpkg.Attachment, expires time.Duration) (string, error) {
	if item == nil || item.DeletedAt != nil {
		return "", errors.New("attachment not found")
	}
	if item.IsPublic && strings.TrimSpace(item.AccessURL) != "" && expires <= 0 {
		return strings.TrimSpace(item.AccessURL), nil
	}
	_, backend, err := m.backendByID(ctx, item.EnvID)
	if err != nil {
		return "", err
	}
	if backend.Type() == "local" {
		if url, urlErr := backend.GetURL(ctx, item.FileKey, expires); urlErr == nil && strings.TrimSpace(url) != "" {
			return url, nil
		}
		return fmt.Sprintf("/api/v1/attachment/%d/download", item.ID), nil
	}
	return backend.GetURL(ctx, item.FileKey, expires)
}

func (m *StorageManager) TestConnection(ctx context.Context, envID int64) error {
	_, backend, err := m.backendByID(ctx, envID)
	if err != nil {
		return err
	}
	return backend.Ping(ctx)
}

func (m *StorageManager) backendByCode(ctx context.Context, envCode string) (*entpkg.StorageEnv, storageBackend, error) {
	env, err := m.ResolveActiveEnv(ctx, envCode)
	if err != nil {
		return nil, nil, err
	}
	backend, err := m.loadBackend(env)
	if err != nil {
		return nil, nil, err
	}
	return env, backend, nil
}

func (m *StorageManager) backendByID(ctx context.Context, envID int64) (*entpkg.StorageEnv, storageBackend, error) {
	env, err := m.ActiveEnvByID(ctx, envID)
	if err != nil {
		return nil, nil, err
	}
	backend, err := m.loadBackend(env)
	if err != nil {
		return nil, nil, err
	}
	return env, backend, nil
}

func (m *StorageManager) loadBackend(env *entpkg.StorageEnv) (storageBackend, error) {
	if env == nil {
		return nil, errors.New("storage env not found")
	}

	m.mu.RLock()
	if backend, ok := m.backends[env.ID]; ok {
		m.mu.RUnlock()
		return backend, nil
	}
	factory, ok := m.factories[strings.ToLower(strings.TrimSpace(env.StorageType))]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unsupported storage type: %s", env.StorageType)
	}

	backend, err := factory(env)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	m.backends[env.ID] = backend
	m.cacheEnvLocked(env)
	m.mu.Unlock()
	return backend, nil
}

func (m *StorageManager) defaultEnv(ctx context.Context) (*entpkg.StorageEnv, error) {
	m.mu.RLock()
	if m.defaultID > 0 {
		if env, ok := m.envs[m.defaultID]; ok {
			m.mu.RUnlock()
			return env, nil
		}
	}
	m.mu.RUnlock()

	env, err := m.ent.StorageEnv.Query().
		Where(
			storageenv.IsDefault(true),
			storageenv.Status(0),
			storageenv.DeletedAtIsNil(),
		).
		First(ctx)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, errors.New("default storage env not found")
		}
		return nil, err
	}
	m.cacheEnv(env)
	return env, nil
}

func (m *StorageManager) cacheEnv(env *entpkg.StorageEnv) {
	if m == nil || env == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheEnvLocked(env)
}

func (m *StorageManager) cacheEnvLocked(env *entpkg.StorageEnv) {
	if env == nil {
		return
	}
	m.envs[env.ID] = env
	if code := normalizeStorageCode(env.Code); code != "" {
		m.codeToID[code] = env.ID
	}
	if env.IsDefault {
		m.defaultID = env.ID
	}
}

func normalizeStorageCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func cleanStorageObjectKey(key string) string {
	key = filepath.ToSlash(strings.TrimSpace(key))
	key = strings.TrimLeft(key, "/")
	return key
}
