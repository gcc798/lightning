package runtimeconfig

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/gcc798/microservice-kit/internal/platform/redislock"
	"github.com/redis/go-redis/v9"
)

type fakeSource struct {
	data  []byte
	loads atomic.Int32
}

func (s *fakeSource) Load(context.Context, string) ([]byte, error) {
	s.loads.Add(1)
	time.Sleep(10 * time.Millisecond)
	return append([]byte(nil), s.data...), nil
}

func TestStoreConcurrentMissLoadsSourceOnce(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	source := &fakeSource{data: []byte(`{"enabled":false,"appId":"","secret":"","templateId":""}`)}
	locker := redislock.New(client, redislock.WithRetryRange(time.Millisecond, 2*time.Millisecond))
	store := NewStore(client, source, locker)

	const callers = 12
	var wg sync.WaitGroup
	errors := make(chan error, callers)
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var cfg WeChatConfig
			errors <- store.Get(context.Background(), CodeWeChat, &cfg)
		}()
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := source.loads.Load(); got != 1 {
		t.Fatalf("source loads = %d, want 1", got)
	}
}

func TestStoreUsesCachedValue(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	source := &fakeSource{data: []byte(`{"enabled":false,"host":"","port":0,"username":"","password":"","from":""}`)}
	store := NewStore(client, source, redislock.New(client))

	for range 2 {
		var cfg EmailConfig
		if err := store.Get(context.Background(), CodeEmail, &cfg); err != nil {
			t.Fatal(err)
		}
	}
	if got := source.loads.Load(); got != 1 {
		t.Fatalf("source loads = %d, want 1", got)
	}
	if !server.Exists(CacheKey(CodeEmail)) {
		t.Fatal("configuration was not cached")
	}
}

func TestValidateKnownConfigurationRejectsUnknownField(t *testing.T) {
	data := []byte(`{"enabled":false,"legacy":true}`)
	if err := Validate(CodeWeChat, data); err == nil {
		t.Fatal("Validate() accepted an unknown field")
	}
}

func TestStoreReplacesInvalidCachedValueFromSource(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	valid := []byte(`{"enabled":false,"appId":"","secret":"","templateId":""}`)
	source := &fakeSource{data: valid}
	store := NewStore(client, source, redislock.New(client))
	if err := client.Set(context.Background(), CacheKey(CodeWeChat), `{"legacy":true}`, 0).Err(); err != nil {
		t.Fatal(err)
	}
	var cfg WeChatConfig
	if err := store.Get(context.Background(), CodeWeChat, &cfg); err != nil {
		t.Fatal(err)
	}
	if source.loads.Load() != 1 {
		t.Fatalf("source loads = %d, want 1", source.loads.Load())
	}
	if got, _ := server.Get(CacheKey(CodeWeChat)); got != string(valid) {
		t.Fatalf("cached value = %s, want %s", got, valid)
	}
}
