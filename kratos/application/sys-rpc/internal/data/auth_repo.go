package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent/authclient"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent/casbinrule"
)

const (
	authClientCachePrefix = "auth:client:key:"
	authClientCacheTTL    = 30 * 24 * time.Hour
)

type AuthClientInfo struct {
	ClientID      string
	GrantType     string
	DeviceType    string
	Timeout       int64
	ActiveTimeout int64
}

func (r *Resources) FindUserByAccount(ctx context.Context, account string) (*v1.UserItem, string, error) {
	account = strings.TrimSpace(account)
	if account == "" {
		return nil, "", nil
	}
	items, err := r.activeUsers(ctx)
	if err != nil {
		return nil, "", err
	}
	for _, item := range items {
		if item.UserName == account || item.Phonenumber == account || item.Email == account {
			return userEntityToItem(item), item.Password, nil
		}
	}
	return nil, "", nil
}

func (r *Resources) AuthenticateClient(ctx context.Context, clientKey, clientSecret, grantType string) (*AuthClientInfo, error) {
	clientKey = strings.TrimSpace(clientKey)
	clientSecret = strings.TrimSpace(clientSecret)
	grantType = strings.TrimSpace(grantType)
	if clientKey == "" || clientSecret == "" {
		return nil, errors.New("clientKey和clientSecret不能为空")
	}
	client, err := r.loadAuthClient(ctx, clientKey)
	if err != nil {
		return nil, err
	}
	if client.Status != 0 {
		return nil, errors.New("客户端已停用")
	}
	if client.ClientSecret != clientSecret {
		return nil, errors.New("客户端密钥错误")
	}
	if grantType != "" && !grantTypeAllowed(client.GrantType, grantType) {
		return nil, fmt.Errorf("客户端不支持授权类型: %s", grantType)
	}
	return &AuthClientInfo{
		ClientID:      client.ClientID,
		GrantType:     client.GrantType,
		DeviceType:    client.DeviceType,
		Timeout:       client.Timeout,
		ActiveTimeout: client.ActiveTimeout,
	}, nil
}

func (r *Resources) loadAuthClient(ctx context.Context, clientKey string) (*entpkg.AuthClient, error) {
	if cached, err := r.authClientFromCache(ctx, clientKey); err == nil && cached != nil {
		return cached, nil
	}
	client, err := r.Ent.AuthClient.Query().
		Where(authclient.ClientKey(clientKey), authclient.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, errors.New("客户端不存在")
		}
		return nil, err
	}
	_ = r.storeAuthClientCache(ctx, client)
	return client, nil
}

func (r *Resources) authClientFromCache(ctx context.Context, clientKey string) (*entpkg.AuthClient, error) {
	if r == nil || r.Redis == nil {
		return nil, nil
	}
	payload, err := r.Redis.Get(ctx, authClientCacheKey(clientKey)).Result()
	if err != nil || strings.TrimSpace(payload) == "" {
		return nil, err
	}
	var client entpkg.AuthClient
	if err = json.Unmarshal([]byte(payload), &client); err != nil {
		_ = r.Redis.Del(ctx, authClientCacheKey(clientKey)).Err()
		return nil, err
	}
	return &client, nil
}

func (r *Resources) storeAuthClientCache(ctx context.Context, client *entpkg.AuthClient) error {
	if r == nil || r.Redis == nil || client == nil {
		return nil
	}
	payload, err := json.Marshal(client)
	if err != nil {
		return err
	}
	return r.Redis.Set(ctx, authClientCacheKey(client.ClientKey), string(payload), authClientCacheTTL).Err()
}

func authClientCacheKey(clientKey string) string {
	return authClientCachePrefix + strings.TrimSpace(clientKey)
}

func grantTypeAllowed(supported, grantType string) bool {
	if grantType == "" {
		return true
	}
	for _, item := range strings.Split(supported, ",") {
		if strings.TrimSpace(item) == grantType {
			return true
		}
	}
	return false
}

func (r *Resources) ValidateAccessToken(ctx context.Context, accessToken string) (*v1.ValidateAccessTokenReply, error) {
	session, err := r.SessionDetails(ctx, accessToken)
	if err != nil {
		return &v1.ValidateAccessTokenReply{Valid: false}, nil
	}
	userID := session.UserID
	user, err := r.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &v1.ValidateAccessTokenReply{Valid: false}, nil
	}
	return &v1.ValidateAccessTokenReply{
		Valid:    true,
		UserId:   userID,
		UserName: user.GetUserName(),
		ClientId: session.ClientID,
	}, nil
}

func (r *Resources) UpdateUserLoginState(ctx context.Context, userID int64, ip string) error {
	if userID <= 0 {
		return nil
	}
	_, err := r.Ent.User.UpdateOneID(userID).
		SetLoginIP(strings.TrimSpace(ip)).
		SetLoginDate(time.Now().Unix()).
		SetUpdatedTime(time.Now()).
		Save(ctx)
	return err
}

func (r *Resources) CreateLoginLogEntry(ctx context.Context, userName, ip, clientID, userAgent string, status int32, msg string) error {
	browser, osName := parseUserAgent(userAgent)
	_, err := r.Ent.LoginLog.Create().
		SetID(nextID()).
		SetUserName(strings.TrimSpace(userName)).
		SetIpaddr(strings.TrimSpace(ip)).
		SetLoginLocation("").
		SetBrowser(browser).
		SetOs(osName).
		SetStatus(status).
		SetMsg(strings.TrimSpace(msg)).
		SetLoginTime(time.Now()).
		SetClientID(strings.TrimSpace(clientID)).
		Save(ctx)
	return err
}

func (r *Resources) EnsureUserRoleRule(ctx context.Context, userID int64, roleKey string) error {
	if userID <= 0 || strings.TrimSpace(roleKey) == "" {
		return nil
	}
	_, err := r.Ent.CasbinRule.Query().
		Where(
			casbinrule.Ptype("g"),
			casbinrule.V0(fmt.Sprintf("user::%d", userID)),
			casbinrule.V1(fmt.Sprintf("role::%s", roleKey)),
		).
		Only(ctx)
	if err == nil {
		return nil
	}
	if err != nil && !entpkg.IsNotFound(err) {
		return err
	}
	_, err = r.Ent.CasbinRule.Create().
		SetID(nextID()).
		SetPtype("g").
		SetV0(fmt.Sprintf("user::%d", userID)).
		SetV1(fmt.Sprintf("role::%s", roleKey)).
		Save(ctx)
	return err
}

func (r *Resources) RemoveUserRoleRule(ctx context.Context, userID int64, roleKey string) error {
	if userID <= 0 || strings.TrimSpace(roleKey) == "" {
		return nil
	}
	_, err := r.Ent.CasbinRule.Delete().
		Where(
			casbinrule.Ptype("g"),
			casbinrule.V0(fmt.Sprintf("user::%d", userID)),
			casbinrule.V1(fmt.Sprintf("role::%s", roleKey)),
		).
		Exec(ctx)
	return err
}

func parseUserAgent(value string) (string, string) {
	lower := strings.ToLower(strings.TrimSpace(value))
	browser := "Unknown"
	switch {
	case strings.Contains(lower, "chrome"):
		browser = "Chrome"
	case strings.Contains(lower, "safari"):
		browser = "Safari"
	case strings.Contains(lower, "firefox"):
		browser = "Firefox"
	case strings.Contains(lower, "edge"):
		browser = "Edge"
	case strings.Contains(lower, "msie") || strings.Contains(lower, "trident"):
		browser = "IE"
	}
	osName := "Unknown"
	switch {
	case strings.Contains(lower, "windows"):
		osName = "Windows"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macos"):
		osName = "macOS"
	case strings.Contains(lower, "android"):
		osName = "Android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ios"):
		osName = "iOS"
	case strings.Contains(lower, "linux"):
		osName = "Linux"
	}
	return browser, osName
}
