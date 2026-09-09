package data

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	entpkg "github.com/gcc798/lightning/kratos/application/sys-rpc/ent"
)

type localStorageBackend struct {
	basePath  string
	urlPrefix string
}

func newLocalStorageBackend(env *entpkg.StorageEnv) (storageBackend, error) {
	if env == nil {
		return nil, errors.New("storage env not found")
	}
	basePath := storageConfigString(env.Config, []string{"basePath"}, filepath.Join("runtime", "attachments", defaultLocalStorageDir(env.Code)))
	basePath = strings.TrimSpace(basePath)
	if basePath == "" {
		basePath = filepath.Join("runtime", "attachments", defaultLocalStorageDir(env.Code))
	}
	basePath = filepath.Clean(basePath)
	urlPrefix := strings.TrimSpace(storageConfigString(env.Config, []string{"urlPrefix"}, ""))
	backend := &localStorageBackend{basePath: basePath, urlPrefix: urlPrefix}
	if err := backend.Ping(context.Background()); err != nil {
		return nil, err
	}
	return backend, nil
}

func (b *localStorageBackend) Type() string { return "local" }

func (b *localStorageBackend) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	fullPath := b.fullPath(key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	return err
}

func (b *localStorageBackend) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	file, err := os.Open(b.fullPath(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("attachment file not found")
		}
		return nil, err
	}
	return file, nil
}

func (b *localStorageBackend) Delete(ctx context.Context, key string) error {
	err := os.Remove(b.fullPath(key))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (b *localStorageBackend) GetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if strings.TrimSpace(b.urlPrefix) == "" {
		return "", errors.New("local storage urlPrefix is not configured")
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(b.urlPrefix, "/"), strings.TrimLeft(cleanStorageObjectKey(key), "/")), nil
}

func (b *localStorageBackend) Exists(ctx context.Context, key string) (bool, error) {
	_, err := os.Stat(b.fullPath(key))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (b *localStorageBackend) Stat(ctx context.Context, key string) (*storageObjectInfo, error) {
	stat, err := os.Stat(b.fullPath(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("attachment file not found")
		}
		return nil, err
	}
	return &storageObjectInfo{Key: cleanStorageObjectKey(key), Size: stat.Size(), ContentType: "application/octet-stream", LastModified: stat.ModTime()}, nil
}

func (b *localStorageBackend) Ping(ctx context.Context) error {
	if err := os.MkdirAll(b.basePath, 0o755); err != nil {
		return err
	}
	probe, err := os.CreateTemp(b.basePath, ".storage-probe-*")
	if err != nil {
		return err
	}
	name := probe.Name()
	if closeErr := probe.Close(); closeErr != nil {
		_ = os.Remove(name)
		return closeErr
	}
	return os.Remove(name)
}

func (b *localStorageBackend) fullPath(key string) string {
	return filepath.Join(b.basePath, filepath.FromSlash(cleanStorageObjectKey(key)))
}

func defaultLocalStorageDir(code string) string {
	code = normalizeStorageCode(code)
	if code == "" {
		return "default"
	}
	return code
}
