package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Jwt 定义业务数据结构。
type Jwt struct {
	secret        []byte
	defaultExpire time.Duration
}

// Claims 定义业务数据结构。
type Claims struct {
	UserId     int64  `json:"userId"`
	UserName   string `json:"userName"`
	ClientId   string `json:"clientId"`
	DeviceType string `json:"deviceType"`
	jwt.RegisteredClaims
}

// New 创建组件实例。
func New(secret string, expireSeconds int64) *Jwt {
	exp := time.Duration(expireSeconds) * time.Second
	if expireSeconds <= 0 {
		exp = 2 * time.Hour
	}
	return &Jwt{secret: []byte(secret), defaultExpire: exp}
}

func (s *Jwt) GenerateToken(userId int64, userName, clientId, deviceType string, expireSeconds ...int64) (string, int64, error) {
	expire := s.defaultExpire
	if len(expireSeconds) > 0 && expireSeconds[0] > 0 {
		expire = time.Duration(expireSeconds[0]) * time.Second
	}
	now := time.Now()
	claims := Claims{
		UserId:     userId,
		UserName:   userName,
		ClientId:   clientId,
		DeviceType: deviceType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "lightning",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}
	return tokenStr, int64(expire.Seconds()), nil
}

func (s *Jwt) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("lightning"),
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func (s *Jwt) DefaultExpireSeconds() int64 { return int64(s.defaultExpire.Seconds()) }
