package data

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret        []byte
	defaultExpire time.Duration
	issuer        string
}

type JWTClaims struct {
	UserID     int64  `json:"userId"`
	UserName   string `json:"userName"`
	ClientID   string `json:"clientId"`
	DeviceType string `json:"deviceType"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, expireSeconds int64, issuer string) *JWTManager {
	if strings.TrimSpace(secret) == "" {
		secret = "lightning-kratos-secret"
	}
	if strings.TrimSpace(issuer) == "" {
		issuer = "NTZ-go"
	}
	expire := time.Duration(expireSeconds) * time.Second
	if expireSeconds <= 0 {
		expire = 30 * time.Minute
	}
	return &JWTManager{
		secret:        []byte(secret),
		defaultExpire: expire,
		issuer:        issuer,
	}
}

func (m *JWTManager) GenerateToken(userID int64, userName, clientID, deviceType string, expire time.Duration) (string, int64, error) {
	if expire <= 0 {
		expire = m.defaultExpire
	}
	now := time.Now()
	claims := JWTClaims{
		UserID:     userID,
		UserName:   strings.TrimSpace(userName),
		ClientID:   strings.TrimSpace(clientID),
		DeviceType: strings.TrimSpace(deviceType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, err
	}
	return tokenStr, int64(expire.Seconds()), nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
