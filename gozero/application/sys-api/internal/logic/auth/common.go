package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	commonauth "github.com/gcc798/microservice-kit/common/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenKeyPrefix   = "refresh_token:"
	refreshTokenIndexPrefix = "refresh_token_index:"
)

type authClient struct {
	ClientId      string
	ClientKey     string
	DeviceType    string
	Timeout       int64
	ActiveTimeout int64
}

type loginUser struct {
	Id          int64
	UserName    string
	NickName    string
	Email       string
	Phonenumber string
	Avatar      string
	UserType    int64
	OrgID       int64
	Roles       []string
	Permissions []string
}

func loginWithRPC(ctx context.Context, svcCtx *svc.ServiceContext, req *types.LoginReq) (*loginUser, *authClient, error) {
	code := req.Code
	if code == "" && req.SmsCode != "" {
		code = req.SmsCode
	}
	resp, err := svcCtx.SysRpcClient.AuthLogin(ctx, &sysservice.AuthLoginReq{
		ClientKey:   req.ClientId,
		GrantType:   req.GrantType,
		Username:    req.Username,
		Password:    req.Password,
		Code:        code,
		Phonenumber: req.Phonenumber,
		Email:       req.Email,
		WxCode:      req.WxCode,
		Uuid:        req.Uuid,
	})
	if err != nil {
		return nil, nil, err
	}
	if resp.UserInfo == nil {
		return nil, nil, fmt.Errorf("登录失败")
	}
	user := &loginUser{
		Id:          resp.UserInfo.UserId,
		UserName:    resp.UserInfo.Username,
		NickName:    resp.UserInfo.Nickname,
		Email:       resp.UserInfo.Email,
		Phonenumber: resp.UserInfo.Phonenumber,
		Avatar:      resp.UserInfo.Avatar,
		UserType:    int64(resp.UserInfo.UserType),
		OrgID:       resp.UserInfo.OrgId,
	}
	if err := enrichLoginUserAuthContext(ctx, svcCtx, user); err != nil {
		return nil, nil, err
	}
	return user, &authClient{
		ClientId:      resp.ClientId,
		ClientKey:     resp.ClientKey,
		DeviceType:    resp.DeviceType,
		Timeout:       resp.Timeout,
		ActiveTimeout: resp.ActiveTimeout,
	}, nil
}

func buildLoginResponse(ctx context.Context, svcCtx *svc.ServiceContext, user *loginUser, client *authClient) (*types.CommonResp, error) {
	accessToken, accessExpiresIn, err := generateAccessToken(user, client, svcCtx.Config.Jwt.Secret)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败")
	}
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("生成Token失败")
	}
	if err := storeRefreshToken(ctx, svcCtx, user, client, refreshToken); err != nil {
		return nil, fmt.Errorf("生成Token失败")
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: map[string]interface{}{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"expires_in":         accessExpiresIn,
		"refresh_expires_in": client.Timeout,
		"user_info": map[string]interface{}{
			"userId":      user.Id,
			"username":    user.UserName,
			"nickname":    user.NickName,
			"phonenumber": user.Phonenumber,
			"email":       user.Email,
			"avatar":      user.Avatar,
			"userType":    user.UserType,
			"orgId":       user.OrgID,
			"roles":       user.Roles,
			"permissions": user.Permissions,
		},
	}}, nil
}

func refreshLoginToken(ctx context.Context, svcCtx *svc.ServiceContext, refreshToken string, clientId string) (*types.CommonResp, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("RefreshToken 无效或已过期")
	}
	indexKey := refreshTokenIndexPrefix + hashToken(refreshToken)
	refreshKey, err := svcCtx.Redis.Get(ctx, indexKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("RefreshToken 无效或已过期")
		}
		return nil, fmt.Errorf("RefreshToken 无效或已过期")
	}
	refreshData, err := svcCtx.Redis.HGetAll(ctx, refreshKey).Result()
	if err != nil || len(refreshData) == 0 {
		return nil, fmt.Errorf("RefreshToken 无效或已过期")
	}
	if refreshData["token"] != refreshToken {
		return nil, fmt.Errorf("RefreshToken 无效")
	}
	// refresh flow uses cached client metadata; clientId must match the stored token owner
	if refreshData["clientKey"] != "" && refreshData["clientKey"] != clientId {
		return nil, fmt.Errorf("客户端不匹配")
	}
	userId, _ := strconv.ParseInt(refreshData["userId"], 10, 64)
	userResp, err := svcCtx.SysRpcClient.UserProfile(ctx, &sysservice.IdReq{Id: userId})
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	client := &authClient{ClientId: refreshData["clientId"], ClientKey: clientId, DeviceType: refreshData["deviceType"], Timeout: parseInt64(refreshData["timeout"]), ActiveTimeout: parseInt64(refreshData["activeTimeout"])}
	user := &loginUser{Id: userResp.UserId, UserName: userResp.UserName, NickName: userResp.NickName, Email: userResp.Email, Phonenumber: userResp.Phonenumber, Avatar: userResp.Avatar, UserType: int64(userResp.UserType), OrgID: userResp.OrgId}
	if err := enrichLoginUserAuthContext(ctx, svcCtx, user); err != nil {
		return nil, fmt.Errorf("用户权限上下文获取失败")
	}
	accessToken, accessExpiresIn, err := generateAccessToken(user, client, svcCtx.Config.Jwt.Secret)
	if err != nil {
		return nil, fmt.Errorf("生成新Token失败")
	}
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("生成新Token失败")
	}
	if err := storeRefreshData(ctx, svcCtx, refreshKey, user, client, newRefreshToken); err != nil {
		return nil, fmt.Errorf("更新RefreshToken失败")
	}
	_ = svcCtx.Redis.Del(ctx, indexKey).Err()
	_ = svcCtx.Redis.Set(ctx, refreshTokenIndexPrefix+hashToken(newRefreshToken), refreshKey, time.Duration(client.Timeout)*time.Second).Err()
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: map[string]interface{}{"access_token": accessToken, "refresh_token": newRefreshToken, "expires_in": accessExpiresIn, "refresh_expires_in": client.Timeout}}, nil
}

func invalidateByToken(ctx context.Context, svcCtx *svc.ServiceContext, token string) {
	if token == "" {
		return
	}
	claims, err := parseAccessToken(token, svcCtx.Config.Jwt.Secret)
	if err != nil {
		return
	}
	refreshKey := buildRefreshKey(claims.UserID, claims.ClientID)
	refreshData, err := svcCtx.Redis.HGetAll(ctx, refreshKey).Result()
	if err == nil && len(refreshData) > 0 && refreshData["token"] != "" {
		_ = svcCtx.Redis.Del(ctx, refreshTokenIndexPrefix+hashToken(refreshData["token"])).Err()
	}
	_ = svcCtx.Redis.Del(ctx, refreshKey).Err()
}

func generateAccessToken(user *loginUser, client *authClient, secret string) (string, int64, error) {
	expireSeconds := client.ActiveTimeout
	if expireSeconds <= 0 {
		expireSeconds = 1800
	}
	expireAt := time.Now().Add(time.Duration(expireSeconds) * time.Second)
	claims := commonauth.AccessClaims{
		UserID:      user.Id,
		UserName:    user.UserName,
		ClientID:    client.ClientId,
		DeviceType:  client.DeviceType,
		OrgID:       user.OrgID,
		Roles:       user.Roles,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "MS_K-gozero",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expireSeconds, nil
}

func parseAccessToken(tokenString, secret string) (*commonauth.AccessClaims, error) {
	return commonauth.ParseAccessToken(tokenString, secret)
}

func generateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func storeRefreshToken(ctx context.Context, svcCtx *svc.ServiceContext, user *loginUser, client *authClient, refreshToken string) error {
	refreshKey := buildRefreshKey(user.Id, client.ClientId)
	if err := storeRefreshData(ctx, svcCtx, refreshKey, user, client, refreshToken); err != nil {
		return err
	}
	return svcCtx.Redis.Set(ctx, refreshTokenIndexPrefix+hashToken(refreshToken), refreshKey, time.Duration(client.Timeout)*time.Second).Err()
}

func storeRefreshData(ctx context.Context, svcCtx *svc.ServiceContext, refreshKey string, user *loginUser, client *authClient, refreshToken string) error {
	data := map[string]interface{}{"token": refreshToken, "userId": strconv.FormatInt(user.Id, 10), "userName": user.UserName, "orgId": strconv.FormatInt(user.OrgID, 10), "clientId": client.ClientId, "clientKey": client.ClientKey, "deviceType": client.DeviceType, "timeout": strconv.FormatInt(client.Timeout, 10), "activeTimeout": strconv.FormatInt(client.ActiveTimeout, 10)}
	if err := svcCtx.Redis.HSet(ctx, refreshKey, data).Err(); err != nil {
		return err
	}
	return svcCtx.Redis.Expire(ctx, refreshKey, time.Duration(client.Timeout)*time.Second).Err()
}

func enrichLoginUserAuthContext(ctx context.Context, svcCtx *svc.ServiceContext, user *loginUser) error {
	authCtx, err := svcCtx.SysRpcClient.UserAuthContext(ctx, &sysservice.UserAuthContextReq{UserId: user.Id})
	if err != nil {
		return fmt.Errorf("获取用户权限上下文失败: %w", err)
	}
	user.OrgID = authCtx.OrgId
	user.Roles = authCtx.Roles
	user.Permissions = authCtx.Permissions
	return nil
}

func buildRefreshKey(userId int64, clientId string) string {
	return refreshTokenKeyPrefix + strconv.FormatInt(userId, 10) + ":" + clientId
}
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
func parseInt64(v string) int64 { n, _ := strconv.ParseInt(v, 10, 64); return n }
