package iam

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gcc798/lightning/application/iam/internal/domain/model"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/platform/jwt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	accessTokenKeyPrefix  = "auth:access:"
	refreshTokenKeyPrefix = "auth:refresh:"
	sessionKeyPrefix      = "auth:session:"
	userSessionsKeyPrefix = "auth:user-sessions:"
)

// TokenManager owns the complete access-token, refresh-token and login-session lifecycle.
type TokenManager interface {
	GenerateTokenPair(ctx context.Context, user *model.User, client *model.AuthClient) (accessToken, refreshToken string, accessExpiresIn, refreshExpiresIn int64, err error)
	ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string, client *model.AuthClient) (newAccessToken, newRefreshToken string, accessExpiresIn, refreshExpiresIn int64, err error)
	RevokeAccessToken(ctx context.Context, accessToken string) error
	RevokeUserSessions(ctx context.Context, userID int64, clientID string) error
}

type tokenManager struct {
	jwt    *jwt.Jwt
	redis  *redis.Client
	logger logging.Logger
}

type tokenSession struct {
	ID          string `json:"id"`
	UserID      int64  `json:"userId"`
	UserName    string `json:"userName"`
	ClientID    string `json:"clientId"`
	DeviceType  string `json:"deviceType"`
	AccessHash  string `json:"accessHash"`
	RefreshHash string `json:"refreshHash"`
}

func NewTokenManager(jwtService *jwt.Jwt, redisClient *redis.Client, logger logging.Logger) TokenManager {
	return &tokenManager{jwt: jwtService, redis: redisClient, logger: logger}
}

func (m *tokenManager) GenerateTokenPair(ctx context.Context, user *model.User, client *model.AuthClient) (string, string, int64, int64, error) {
	if user == nil || client == nil {
		return "", "", 0, 0, errors.New("user and client are required")
	}

	accessToken, accessTTL, err := m.jwt.GenerateToken(
		user.ID,
		user.UserName,
		client.ClientId,
		client.DeviceType,
		client.ActiveTimeout,
	)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("generate access token: %w", err)
	}
	claims, err := m.jwt.ValidateToken(accessToken)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("read generated access token: %w", err)
	}

	refreshToken, err := generateRandomToken(32)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshTTL := client.Timeout
	if refreshTTL <= 0 {
		return "", "", 0, 0, errors.New("client refresh token timeout must be positive")
	}
	session := tokenSession{
		ID:          claims.ID,
		UserID:      user.ID,
		UserName:    user.UserName,
		ClientID:    client.ClientId,
		DeviceType:  client.DeviceType,
		AccessHash:  tokenHash(accessToken),
		RefreshHash: tokenHash(refreshToken),
	}
	if err := m.storeSession(ctx, session, time.Duration(accessTTL)*time.Second, time.Duration(refreshTTL)*time.Second); err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, refreshToken, accessTTL, refreshTTL, nil
}

func (m *tokenManager) ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error) {
	claims, err := m.jwt.ValidateToken(token)
	if err != nil {
		return nil, errors.New("AccessToken 无效或已过期")
	}
	active, err := m.redis.Exists(ctx, accessTokenKey(tokenHash(token))).Result()
	if err != nil {
		m.logger.Error("validate access token session failed", zap.Error(err))
		return nil, errors.New("认证服务暂时不可用")
	}
	if active != 1 {
		return nil, errors.New("AccessToken 已失效")
	}
	return claims, nil
}

func (m *tokenManager) RefreshAccessToken(ctx context.Context, refreshToken string, client *model.AuthClient) (string, string, int64, int64, error) {
	if refreshToken == "" || client == nil {
		return "", "", 0, 0, errors.New("RefreshToken 和客户端不能为空")
	}

	// GETDEL makes a refresh token single-use, including under concurrent refresh requests.
	encoded, err := m.redis.GetDel(ctx, refreshTokenKey(tokenHash(refreshToken))).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", 0, 0, errors.New("RefreshToken 无效或已过期")
		}
		return "", "", 0, 0, fmt.Errorf("consume refresh token: %w", err)
	}

	var oldSession tokenSession
	if err := json.Unmarshal([]byte(encoded), &oldSession); err != nil {
		return "", "", 0, 0, errors.New("RefreshToken 会话数据无效")
	}
	if oldSession.ClientID != client.ClientId {
		return "", "", 0, 0, errors.New("客户端不匹配")
	}

	if err := m.deleteSession(ctx, oldSession, false); err != nil {
		m.logger.Warn("delete rotated token session failed", zap.Error(err))
	}
	user := &model.User{ID: oldSession.UserID, UserName: oldSession.UserName}
	return m.GenerateTokenPair(ctx, user, client)
}

func (m *tokenManager) RevokeAccessToken(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return nil
	}
	sessionID, err := m.redis.Get(ctx, accessTokenKey(tokenHash(accessToken))).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("read access token session: %w", err)
	}
	session, err := m.readSession(ctx, sessionID)
	if err != nil {
		_ = m.redis.Del(ctx, accessTokenKey(tokenHash(accessToken))).Err()
		return err
	}
	return m.deleteSession(ctx, session, true)
}

func (m *tokenManager) RevokeUserSessions(ctx context.Context, userID int64, clientID string) error {
	setKey := userSessionsKey(userID, clientID)
	sessionIDs, err := m.redis.SMembers(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("list user sessions: %w", err)
	}
	for _, sessionID := range sessionIDs {
		session, readErr := m.readSession(ctx, sessionID)
		if readErr != nil {
			m.logger.Warn("read user session during revocation failed", zap.String("sessionId", sessionID), zap.Error(readErr))
			continue
		}
		if deleteErr := m.deleteSession(ctx, session, true); deleteErr != nil {
			return deleteErr
		}
	}
	return m.redis.Del(ctx, setKey).Err()
}

func (m *tokenManager) storeSession(ctx context.Context, session tokenSession, accessTTL, refreshTTL time.Duration) error {
	encoded, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("encode token session: %w", err)
	}
	setKey := userSessionsKey(session.UserID, session.ClientID)
	pipe := m.redis.TxPipeline()
	pipe.Set(ctx, accessTokenKey(session.AccessHash), session.ID, accessTTL)
	pipe.Set(ctx, refreshTokenKey(session.RefreshHash), encoded, refreshTTL)
	pipe.Set(ctx, sessionKey(session.ID), encoded, refreshTTL)
	pipe.SAdd(ctx, setKey, session.ID)
	pipe.Expire(ctx, setKey, refreshTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("store token session: %w", err)
	}
	return nil
}

func (m *tokenManager) readSession(ctx context.Context, sessionID string) (tokenSession, error) {
	encoded, err := m.redis.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return tokenSession{}, fmt.Errorf("read token session: %w", err)
	}
	var session tokenSession
	if err := json.Unmarshal([]byte(encoded), &session); err != nil {
		return tokenSession{}, fmt.Errorf("decode token session: %w", err)
	}
	return session, nil
}

func (m *tokenManager) deleteSession(ctx context.Context, session tokenSession, deleteRefresh bool) error {
	keys := []string{accessTokenKey(session.AccessHash), sessionKey(session.ID)}
	if deleteRefresh {
		keys = append(keys, refreshTokenKey(session.RefreshHash))
	}
	pipe := m.redis.TxPipeline()
	pipe.Del(ctx, keys...)
	pipe.SRem(ctx, userSessionsKey(session.UserID, session.ClientID), session.ID)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("delete token session: %w", err)
	}
	return nil
}

func generateRandomToken(length int) (string, error) {
	value := make([]byte, length)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func accessTokenKey(hash string) string  { return accessTokenKeyPrefix + hash }
func refreshTokenKey(hash string) string { return refreshTokenKeyPrefix + hash }
func sessionKey(sessionID string) string { return sessionKeyPrefix + sessionID }
func userSessionsKey(userID int64, clientID string) string {
	return userSessionsKeyPrefix + strconv.FormatInt(userID, 10) + ":" + clientID
}
