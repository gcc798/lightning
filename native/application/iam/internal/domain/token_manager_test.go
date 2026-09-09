package iam

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gcc798/lightning/application/iam/internal/domain/model"
	logging "github.com/gcc798/lightning/internal/logger"
	jwtservice "github.com/gcc798/lightning/internal/platform/jwt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type tokenTestLogger struct{}

func (tokenTestLogger) Get() *zap.Logger                   { return zap.NewNop() }
func (tokenTestLogger) Debug(string, ...zap.Field)         {}
func (tokenTestLogger) Info(string, ...zap.Field)          {}
func (tokenTestLogger) Warn(string, ...zap.Field)          {}
func (tokenTestLogger) Error(string, ...zap.Field)         {}
func (tokenTestLogger) Fatal(string, ...zap.Field)         {}
func (l tokenTestLogger) With(...zap.Field) logging.Logger { return l }

func TestTokenManagerLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager, closeRedis := newTokenManagerForTest(t)
	defer closeRedis()
	user := &model.User{ID: 42, UserName: "tester"}
	client := &model.AuthClient{ClientId: "web", DeviceType: "web", ActiveTimeout: 1800, Timeout: 604800}

	access, refresh, _, _, err := manager.GenerateTokenPair(ctx, user, client)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}
	if _, err := manager.ValidateAccessToken(ctx, access); err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	newAccess, newRefresh, _, _, err := manager.RefreshAccessToken(ctx, refresh, client)
	if err != nil {
		t.Fatalf("RefreshAccessToken() error = %v", err)
	}
	if _, _, _, _, err := manager.RefreshAccessToken(ctx, refresh, client); err == nil {
		t.Fatal("old refresh token remained usable after rotation")
	}
	if _, err := manager.ValidateAccessToken(ctx, access); err == nil {
		t.Fatal("old access token remained usable after refresh rotation")
	}
	if _, err := manager.ValidateAccessToken(ctx, newAccess); err != nil {
		t.Fatalf("new access token is invalid: %v", err)
	}

	if err := manager.RevokeAccessToken(ctx, newAccess); err != nil {
		t.Fatalf("RevokeAccessToken() error = %v", err)
	}
	if _, err := manager.ValidateAccessToken(ctx, newAccess); err == nil {
		t.Fatal("revoked access token remained usable")
	}
	if _, _, _, _, err := manager.RefreshAccessToken(ctx, newRefresh, client); err == nil {
		t.Fatal("refresh token remained usable after logout")
	}
}

func TestTokenManagerRevokeUserSessions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager, closeRedis := newTokenManagerForTest(t)
	defer closeRedis()
	user := &model.User{ID: 7, UserName: "parallel-user"}
	client := &model.AuthClient{ClientId: "mobile", DeviceType: "mobile", ActiveTimeout: 1800, Timeout: 604800}

	firstAccess, firstRefresh, _, _, err := manager.GenerateTokenPair(ctx, user, client)
	if err != nil {
		t.Fatal(err)
	}
	secondAccess, secondRefresh, _, _, err := manager.GenerateTokenPair(ctx, user, client)
	if err != nil {
		t.Fatal(err)
	}

	if err := manager.RevokeUserSessions(ctx, user.ID, client.ClientId); err != nil {
		t.Fatalf("RevokeUserSessions() error = %v", err)
	}
	for _, access := range []string{firstAccess, secondAccess} {
		if _, err := manager.ValidateAccessToken(ctx, access); err == nil {
			t.Fatal("access token remained usable after revoking all user sessions")
		}
	}
	for _, refresh := range []string{firstRefresh, secondRefresh} {
		if _, _, _, _, err := manager.RefreshAccessToken(ctx, refresh, client); err == nil {
			t.Fatal("refresh token remained usable after revoking all user sessions")
		}
	}
}

func TestTokenHash(t *testing.T) {
	t.Parallel()

	if got, want := tokenHash("token"), "3c469e9d6c5875d37a43f353d4f88e61fcf812c66eee3457465a40b0da4153e0"; got != want {
		t.Fatalf("tokenHash() = %q, want %q", got, want)
	}
	if tokenHash("token") == tokenHash("other") {
		t.Fatal("different tokens produced the same test hash")
	}
}

func newTokenManagerForTest(t *testing.T) (TokenManager, func()) {
	t.Helper()

	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	manager := NewTokenManager(jwtservice.New("test-secret-with-at-least-32-bytes", 3600), client, tokenTestLogger{})
	return manager, func() {
		_ = client.Close()
		server.Close()
	}
}
