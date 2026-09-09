package iam

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gcc798/lightning/application/iam/internal/domain/model"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	"github.com/gcc798/lightning/internal/config"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/platform/captcha"
	"github.com/gcc798/lightning/internal/platform/thirdparty/wechat"
	"github.com/gcc798/lightning/internal/utils"
	"github.com/gcc798/lightning/internal/utils/idgen"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	miniProgramUserType int32 = 1
	passwordErrorLimit        = 5
	passwordLockTTL           = 10 * time.Minute
)

type WeChatSessionProvider interface {
	Code2Session(ctx context.Context, wxCode string) (*wechat.Code2SessionResponse, error)
}

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Refresh(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, accessToken string) error
}

type authService struct {
	db             *gorm.DB
	redis          *redis.Client
	config         *config.Config
	logger         logging.Logger
	clients        ClientService
	tokens         TokenManager
	captcha        CaptchaService
	wechat         WeChatSessionProvider
	logs           sysv1.API
	authenticators map[string]loginAuthenticator
}

type loginAuthenticator interface {
	Authenticate(ctx context.Context, req *LoginRequest) (*model.User, error)
}

func NewAuthService(
	db *gorm.DB,
	redisClient *redis.Client,
	cfg *config.Config,
	logger logging.Logger,
	clients ClientService,
	tokens TokenManager,
	captchaService CaptchaService,
	wechatProvider WeChatSessionProvider,
	logs sysv1.API,
) AuthService {
	s := &authService{
		db:      db,
		redis:   redisClient,
		config:  cfg,
		logger:  logger,
		clients: clients,
		tokens:  tokens,
		captcha: captchaService,
		wechat:  wechatProvider,
		logs:    logs,
	}
	s.authenticators = map[string]loginAuthenticator{
		"password": &passwordAuthenticator{service: s},
		"email":    &emailAuthenticator{service: s},
		"sms":      &smsAuthenticator{service: s},
		"wechat":   &wechatAuthenticator{service: s},
	}
	return s
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	client, err := s.clients.AuthenticateClientID(ctx, req.ClientID, req.GrantType)
	if err != nil {
		s.recordLogin(ctx, req, resolveLoginAccount(req), req.ClientID, 1, err.Error())
		return nil, err
	}
	authenticator, ok := s.authenticators[req.GrantType]
	if !ok {
		err := fmt.Errorf("不支持的授权类型: %s", req.GrantType)
		s.recordLogin(ctx, req, resolveLoginAccount(req), client.ClientId, 1, err.Error())
		return nil, err
	}
	user, err := authenticator.Authenticate(ctx, req)
	if err != nil {
		s.recordLogin(ctx, req, resolveLoginAccount(req), client.ClientId, 1, err.Error())
		return nil, err
	}
	if !s.config.Auth.AllowConcurrent {
		if err := s.tokens.RevokeUserSessions(ctx, user.ID, client.ClientId); err != nil {
			return nil, fmt.Errorf("清理旧会话失败: %w", err)
		}
	}
	accessToken, refreshToken, accessTTL, refreshTTL, err := s.tokens.GenerateTokenPair(ctx, user, client)
	if err != nil {
		s.recordLogin(ctx, req, user.UserName, client.ClientId, 1, "生成 Token 失败")
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}
	if err := (&model.User{}).UpdateLoginInfo(s.db.WithContext(ctx), user.ID, req.LoginIP, time.Now().Unix()); err != nil {
		s.logger.Warn("update login information failed", zap.Error(err))
	}
	s.recordLogin(ctx, req, user.UserName, client.ClientId, 0, "登录成功")
	userInfo := newLoginUserInfo(user)
	return &LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessTTL,
		RefreshExpiresIn: refreshTTL,
		ClientID:         client.ClientId,
		OpenID:           userInfo.OpenID,
		UserInfo:         userInfo,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	client, err := s.clients.GetActiveClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}
	access, refresh, accessTTL, refreshTTL, err := s.tokens.RefreshAccessToken(ctx, req.RefreshToken, client)
	if err != nil {
		return nil, err
	}
	return &RefreshTokenResponse{
		AccessToken:      access,
		RefreshToken:     refresh,
		ExpiresIn:        accessTTL,
		RefreshExpiresIn: refreshTTL,
	}, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string) error {
	return s.tokens.RevokeAccessToken(ctx, accessToken)
}

type passwordAuthenticator struct{ service *authService }

func (a *passwordAuthenticator) Authenticate(ctx context.Context, req *LoginRequest) (*model.User, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}
	imageCaptchaEnabled, err := a.service.captcha.ImageEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("读取验证码配置失败: %w", err)
	}
	if imageCaptchaEnabled {
		if err := a.service.verifyCaptcha(ctx, captcha.CaptchaTypeImage, req.Uuid, req.Code); err != nil {
			return nil, err
		}
	}
	if err := a.checkBruteForce(ctx, req.Username); err != nil {
		return nil, err
	}
	user, err := (&model.User{}).FindByUsername(a.service.db.WithContext(ctx), req.Username)
	if err != nil {
		a.incrementErrorCount(ctx, req.Username)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		a.incrementErrorCount(ctx, req.Username)
		return nil, errors.New("用户名或密码错误")
	}
	if !user.CanLogin() {
		return nil, errors.New("用户已被停用")
	}
	a.clearErrorCount(ctx, req.Username)
	return user, nil
}

func (a *passwordAuthenticator) checkBruteForce(ctx context.Context, username string) error {
	key := "auth:password-errors:" + username
	count, err := a.service.redis.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("读取登录限制失败: %w", err)
	}
	if count < passwordErrorLimit {
		return nil
	}
	ttl, _ := a.service.redis.TTL(ctx, key).Result()
	return fmt.Errorf("密码错误次数过多，请%d分钟后再试", int(ttl.Minutes())+1)
}

func (a *passwordAuthenticator) incrementErrorCount(ctx context.Context, username string) {
	key := "auth:password-errors:" + username
	pipe := a.service.redis.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, passwordLockTTL)
	_, _ = pipe.Exec(ctx)
}

func (a *passwordAuthenticator) clearErrorCount(ctx context.Context, username string) {
	_ = a.service.redis.Del(ctx, "auth:password-errors:"+username).Err()
}

type emailAuthenticator struct{ service *authService }

func (a *emailAuthenticator) Authenticate(ctx context.Context, req *LoginRequest) (*model.User, error) {
	if req.Email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	if err := a.service.verifyCaptcha(ctx, captcha.CaptchaTypeEmail, req.Uuid, req.Code); err != nil {
		return nil, err
	}
	user, err := (&model.User{}).FindByEmail(a.service.db.WithContext(ctx), req.Email)
	return validateLoginUser(user, err, "邮箱或验证码错误")
}

type smsAuthenticator struct{ service *authService }

func (a *smsAuthenticator) Authenticate(ctx context.Context, req *LoginRequest) (*model.User, error) {
	if req.Phonenumber == "" {
		return nil, errors.New("手机号不能为空")
	}
	if err := a.service.verifyCaptcha(ctx, captcha.CaptchaTypeSMS, req.Uuid, req.Code); err != nil {
		return nil, err
	}
	user, err := (&model.User{}).FindByPhonenumber(a.service.db.WithContext(ctx), req.Phonenumber)
	return validateLoginUser(user, err, "手机号或验证码错误")
}

type wechatAuthenticator struct{ service *authService }

func (a *wechatAuthenticator) Authenticate(ctx context.Context, req *LoginRequest) (*model.User, error) {
	if req.WxCode == "" {
		return nil, errors.New("微信 code 不能为空")
	}
	if a.service.wechat == nil {
		return nil, errors.New("微信登录未启用")
	}
	session, err := a.service.wechat.Code2Session(ctx, req.WxCode)
	if err != nil {
		return nil, err
	}
	if session.OpenID == "" {
		return nil, errors.New("微信 OpenID 为空")
	}
	users := &model.User{}
	user, err := users.FindByOpenId(a.service.db.WithContext(ctx), session.OpenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询微信用户失败: %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.User{
			UserName:  wechatUsername(session.OpenID),
			NickName:  "微信用户",
			UserType:  miniProgramUserType,
			OrgID:     0,
			Status:    0,
			OpenId:    session.OpenID,
			UnionId:   session.UnionID,
			Sex:       2,
			LoginIp:   req.LoginIP,
			LoginDate: time.Now().Unix(),
		}
		if err := a.service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := users.Create(tx, user); err != nil {
				return err
			}
			var defaultRole model.Role
			if err := tx.Where("role_key = ? AND status = 0", "user").First(&defaultRole).Error; err != nil {
				return fmt.Errorf("查找默认角色失败: %w", err)
			}
			return tx.Create(&model.MUserRole{UserId: user.ID, RoleId: defaultRole.ID}).Error
		}); err != nil {
			return nil, fmt.Errorf("创建微信用户失败: %w", err)
		}
	} else if session.UnionID != "" && user.UnionId != session.UnionID {
		if err := a.service.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Update("union_id", session.UnionID).Error; err != nil {
			return nil, fmt.Errorf("更新微信用户失败: %w", err)
		}
		user.UnionId = session.UnionID
	}
	if !user.CanLogin() {
		return nil, errors.New("用户已被停用")
	}
	return user, nil
}

func (s *authService) verifyCaptcha(ctx context.Context, captchaType captcha.CaptchaType, captchaID, code string) error {
	if captchaID == "" || code == "" {
		return errors.New("验证码 ID 和验证码不能为空")
	}
	return s.captcha.Verify(ctx, captchaType, map[string]interface{}{"captchaID": captchaID, "code": code})
}

func validateLoginUser(user *model.User, err error, credentialError string) (*model.User, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(credentialError)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if !user.CanLogin() {
		return nil, errors.New("用户已被停用")
	}
	return user, nil
}

func (s *authService) recordLogin(ctx context.Context, req *LoginRequest, username, clientID string, status int32, message string) {
	if s.logs == nil {
		return
	}
	browser, osName := parseUserAgent(req.UserAgent)
	entry := &sysv1.RecordLoginRequest{
		Id: idgen.MustNextID(), UserName: username, IpAddress: req.LoginIP,
		Browser: browser, Os: osName, Status: status, Message: message,
		LoginTimeUnixMilli: time.Now().UnixMilli(), ClientId: clientID,
	}
	if err := s.logs.RecordLogin(ctx, entry); err != nil {
		s.logger.Error("record login log failed", zap.Error(err))
	}
}

func resolveLoginAccount(req *LoginRequest) string {
	for _, value := range []string{req.Username, req.Email, req.Phonenumber} {
		if value != "" {
			return value
		}
	}
	return "-"
}

func newLoginUserInfo(user *model.User) *UserInfo {
	return &UserInfo{
		UserId:      user.ID,
		Username:    user.UserName,
		Nickname:    user.NickName,
		Phonenumber: user.Phonenumber,
		Email:       user.Email,
		Avatar:      user.Avatar,
		OrgID:       user.OrgID,
		UserType:    user.UserType,
		OpenID:      user.OpenId,
		UnionID:     user.UnionId,
	}
}

func wechatUsername(openID string) string {
	digest := sha256.Sum256([]byte(openID))
	return "wx_" + hex.EncodeToString(digest[:16])
}

func parseUserAgent(userAgent string) (browser, operatingSystem string) {
	lower := strings.ToLower(userAgent)
	switch {
	case strings.Contains(lower, "edg"):
		browser = "Edge"
	case strings.Contains(lower, "chrome"):
		browser = "Chrome"
	case strings.Contains(lower, "safari"):
		browser = "Safari"
	case strings.Contains(lower, "firefox"):
		browser = "Firefox"
	default:
		browser = "Unknown"
	}
	switch {
	case strings.Contains(lower, "windows"):
		operatingSystem = "Windows"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macos"):
		operatingSystem = "macOS"
	case strings.Contains(lower, "android"):
		operatingSystem = "Android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ios"):
		operatingSystem = "iOS"
	case strings.Contains(lower, "linux"):
		operatingSystem = "Linux"
	default:
		operatingSystem = "Unknown"
	}
	return browser, operatingSystem
}
