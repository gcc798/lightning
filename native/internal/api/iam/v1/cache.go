package iamv1

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/gcc798/microservice-kit/internal/platform/jwt"
)

type cacheEntry[T any] struct {
	value   T
	expires time.Time
}

type Cached struct {
	next        API
	ttl         time.Duration
	mu          sync.RWMutex
	tokens      map[string]cacheEntry[*jwt.Claims]
	permissions map[string]cacheEntry[bool]
}

func NewCached(next API, ttl time.Duration) *Cached {
	return &Cached{next: next, ttl: ttl, tokens: make(map[string]cacheEntry[*jwt.Claims]), permissions: make(map[string]cacheEntry[bool])}
}

func (c *Cached) ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error) {
	key := tokenKey(token)
	if value, ok := cachedValue(c, c.tokens, key); ok {
		return value, nil
	}
	value, err := c.next.ValidateAccessToken(ctx, token)
	if err == nil {
		storeValue(c, c.tokens, key, value)
	}
	return value, err
}

func (c *Cached) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	key := fmt.Sprintf("%d:%s:%s", userID, resource, action)
	if value, ok := cachedValue(c, c.permissions, key); ok {
		return value, nil
	}
	value, err := c.next.CheckPermission(ctx, userID, resource, action)
	if err == nil {
		storeValue(c, c.permissions, key, value)
	}
	return value, err
}

func cachedValue[T any](c *Cached, values map[string]cacheEntry[T], key string) (T, bool) {
	c.mu.RLock()
	entry, ok := values[key]
	c.mu.RUnlock()
	if ok && time.Now().Before(entry.expires) {
		return entry.value, true
	}
	var zero T
	return zero, false
}

func storeValue[T any](c *Cached, values map[string]cacheEntry[T], key string, value T) {
	c.mu.Lock()
	values[key] = cacheEntry[T]{value: value, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func tokenKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
