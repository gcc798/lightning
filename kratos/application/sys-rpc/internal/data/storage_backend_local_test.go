package data

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"

	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
)

func TestLocalStorageBackendLifecycle(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	backend, err := newLocalStorageBackend(&entpkg.StorageEnv{
		Code: "local-test",
		Config: map[string]any{
			"basePath":  root,
			"urlPrefix": "http://example.com/files",
		},
	})
	if err != nil {
		t.Fatalf("newLocalStorageBackend() error = %v", err)
	}

	ctx := context.Background()
	key := filepath.ToSlash(filepath.Join("nested", "hello.txt"))
	content := []byte("hello kratos storage")
	if err := backend.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain"); err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	exists, err := backend.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Fatalf("Exists() = false, want true")
	}

	stat, err := backend.Stat(ctx, key)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if stat.Size != int64(len(content)) {
		t.Fatalf("Stat().Size = %d, want %d", stat.Size, len(content))
	}

	reader, err := backend.Download(ctx, key)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	defer reader.Close()
	actual, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(actual) != string(content) {
		t.Fatalf("downloaded content = %q, want %q", string(actual), string(content))
	}

	url, err := backend.GetURL(ctx, key, 0)
	if err != nil {
		t.Fatalf("GetURL() error = %v", err)
	}
	if url != "http://example.com/files/nested/hello.txt" {
		t.Fatalf("GetURL() = %q", url)
	}

	if err := backend.Delete(ctx, key); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	deleted, err := backend.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists() after delete error = %v", err)
	}
	if deleted {
		t.Fatalf("Exists() after delete = true, want false")
	}
}
