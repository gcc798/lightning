package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// Code2SessionURL 定义业务常量。
	Code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"
	// AccessTokenURL 定义业务常量。
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// SendTemplateURL 定义业务常量。
	SendTemplateURL = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=%s"
	// AccessTokenCacheKey 定义业务常量。
	AccessTokenCacheKey = "wechat:access_token"
	// AccessTokenTTL 定义业务常量。
	AccessTokenTTL = 115 * time.Minute // 微信token有效期2小时，提前5分钟刷新
)

// Config 定义业务数据结构。
type Config struct {
	Enabled bool   `json:"enabled"`
	AppID   string `json:"appid"`
	Secret  string `json:"secret"`
}

// Manager 定义业务数据结构。
type Manager struct {
	config Config
	logger logging.Logger
	redis  *redis.Client
	client *http.Client
}

// NewManager 创建组件实例。
func NewManager(config Config, logger logging.Logger, redis *redis.Client) *Manager {
	return &Manager{
		config: config,
		logger: logger,
		redis:  redis,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Code2SessionResponse 定义业务数据结构。
type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Code2Session 执行业务逻辑。
func (s *Manager) Code2Session(ctx context.Context, wxCode string) (*Code2SessionResponse, error) {
	if !s.config.Enabled {
		return nil, fmt.Errorf("微信小程序登录未启用")
	}
	if strings.TrimSpace(s.config.AppID) == "" || strings.TrimSpace(s.config.Secret) == "" {
		return nil, fmt.Errorf("微信小程序配置缺失")
	}
	if wxCode == "" {
		return nil, fmt.Errorf("微信code不能为空")
	}
	query := url.Values{}
	query.Set("appid", s.config.AppID)
	query.Set("secret", s.config.Secret)
	query.Set("js_code", wxCode)
	query.Set("grant_type", "authorization_code")
	endpoint := Code2SessionURL + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("创建微信请求失败")
	}
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("failed to call wechat code2session", zap.Error(err))
		return nil, fmt.Errorf("调用微信接口失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("failed to read wechat response", zap.Error(err))
		return nil, fmt.Errorf("读取微信响应失败")
	}
	var result Code2SessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		s.logger.Error("failed to unmarshal wechat response", zap.Error(err))
		return nil, fmt.Errorf("解析微信响应失败")
	}
	if result.ErrCode != 0 {
		s.logger.Error("wechat code2session failed", zap.Int("errcode", result.ErrCode), zap.String("errmsg", result.ErrMsg))
		return nil, fmt.Errorf("微信授权失败: %s", result.ErrMsg)
	}
	return &result, nil
}

// AccessTokenResponse 获取Access Token响应
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// GetAccessToken 获取Access Token（优先从Redis缓存读取）
func (s *Manager) GetAccessToken(ctx context.Context) (string, error) {
	// 1. 尝试从Redis获取缓存的token
	if s.redis != nil {
		cachedToken, err := s.redis.Get(ctx, AccessTokenCacheKey).Result()
		if err == nil && cachedToken != "" {
			s.logger.Debug("access token fetched from cache")
			return cachedToken, nil
		}
	}

	// 2. 缓存未命中，从微信API获取新token
	s.logger.Info("fetching new access token from WeChat API")
	token, expiresIn, err := s.fetchAccessToken()
	if err != nil {
		s.logger.Error("failed to fetch access token", zap.Error(err))
		return "", err
	}

	// 3. 将新token缓存到Redis
	if s.redis != nil {
		ttl := time.Duration(expiresIn-300) * time.Second // 提前5分钟过期
		if ttl <= 0 {
			ttl = AccessTokenTTL
		}
		if err := s.redis.Set(ctx, AccessTokenCacheKey, token, ttl).Err(); err != nil {
			s.logger.Warn("failed to cache access token", zap.Error(err))
		}
		s.logger.Info("access token fetched and cached successfully",
			zap.String("token", token[:10]+"..."),
			zap.Duration("ttl", ttl))
	}

	return token, nil
}

// fetchAccessToken 从微信API获取Access Token
func (s *Manager) fetchAccessToken() (string, int64, error) {
	url := fmt.Sprintf(AccessTokenURL, s.config.AppID, s.config.Secret)
	resp, err := s.client.Get(url)
	if err != nil {
		return "", 0, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("读取响应失败: %w", err)
	}

	var result AccessTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", 0, fmt.Errorf("微信API错误: [%d] %s", result.ErrCode, result.ErrMsg)
	}

	if result.AccessToken == "" {
		return "", 0, fmt.Errorf("获取到的access_token为空")
	}

	return result.AccessToken, result.ExpiresIn, nil
}

// TemplateData 模板消息数据
type TemplateData struct {
	Value string `json:"value"`
}

// TemplateMessage 模板消息结构
type TemplateMessage struct {
	ToUser           string                  `json:"touser"`                      // 接收者OpenID
	TemplateID       string                  `json:"template_id"`                 // 模板ID
	Page             string                  `json:"page,omitempty"`              // 点击后跳转页面（可选）
	Data             map[string]TemplateData `json:"data"`                        // 模板数据
	MiniprogramState string                  `json:"miniprogram_state,omitempty"` // 小程序状态（developer/trial/formal）
}

// TemplateMessageResponse 模板消息发送响应
type TemplateMessageResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgID   int64  `json:"msgid"`
}

// SendTemplateMessage 发送模板消息
func (s *Manager) SendTemplateMessage(ctx context.Context, msg *TemplateMessage) error {
	// 1. 获取Access Token
	accessToken, err := s.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("获取access_token失败: %w", err)
	}

	// 2. 构造请求
	url := fmt.Sprintf(SendTemplateURL, accessToken)
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 4. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result TemplateMessageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 5. 检查结果
	if result.ErrCode != 0 {
		// 如果是access_token过期错误，清除缓存并重试一次
		if result.ErrCode == 40001 || result.ErrCode == 42001 {
			s.logger.Warn("access_token expired, clearing cache and retrying",
				zap.Int("errcode", result.ErrCode),
				zap.String("errmsg", result.ErrMsg))
			if s.redis != nil {
				_ = s.redis.Del(ctx, AccessTokenCacheKey).Err()
			}
			// 重试一次
			return s.sendTemplateMessageWithRetry(ctx, msg)
		}
		return fmt.Errorf("发送失败: [%d] %s", result.ErrCode, result.ErrMsg)
	}

	s.logger.Info("template message sent successfully",
		zap.String("toUser", msg.ToUser),
		zap.String("templateId", msg.TemplateID),
		zap.Int64("msgId", result.MsgID))

	return nil
}

// sendTemplateMessageWithRetry 重试发送模板消息（用于token过期场景）
func (s *Manager) sendTemplateMessageWithRetry(ctx context.Context, msg *TemplateMessage) error {
	// 重新获取token并发送
	accessToken, err := s.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("重新获取access_token失败: %w", err)
	}

	url := fmt.Sprintf(SendTemplateURL, accessToken)
	payload, _ := json.Marshal(msg)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("创建重试请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("重试HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result TemplateMessageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析重试响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("重试发送失败: [%d] %s", result.ErrCode, result.ErrMsg)
	}

	s.logger.Info("template message sent successfully after retry",
		zap.String("toUser", msg.ToUser),
		zap.Int64("msgId", result.MsgID))

	return nil
}

// SendDeviceControlNotification 发送设备控制结果通知
func (s *Manager) SendDeviceControlNotification(ctx context.Context, deviceName, state string, deviceId int64, openIDs []string, templateID string) {
	if len(openIDs) == 0 {
		s.logger.Debug("no openIDs to notify, skip sending WeChat notification")
		return
	}

	// 构造模板消息
	msg := &TemplateMessage{
		TemplateID:       templateID,
		Page:             fmt.Sprintf("pages/device/detail?id=%d", deviceId),
		MiniprogramState: "formal", // 正式版
		Data: map[string]TemplateData{
			"thing1": {Value: deviceName},                               // 设备名称
			"thing2": {Value: state},                                    // 状态
			"time3":  {Value: time.Now().Format("2006-01-02 15:04:05")}, // 时间
		},
	}

	// 批量发送给所有用户
	for _, openID := range openIDs {
		if openID == "" {
			continue
		}
		msg.ToUser = openID
		if err := s.SendTemplateMessage(ctx, msg); err != nil {
			s.logger.Error("failed to send template message",
				zap.String("openID", openID),
				zap.String("deviceName", deviceName),
				zap.Error(err))
		} else {
			s.logger.Info("device control notification sent",
				zap.String("openID", openID),
				zap.String("deviceName", deviceName),
				zap.String("state", state))
		}
	}
}
