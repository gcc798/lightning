package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/conf"
)

const weChatCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

type WeChatManager struct {
	appID  string
	secret string
	client *http.Client
}

type WeChatCode2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func NewWeChatManager(cfg *conf.WeChat) *WeChatManager {
	appID := strings.TrimSpace(cfg.GetAppId())
	secret := strings.TrimSpace(cfg.GetSecret())
	if appID == "" || secret == "" {
		return nil
	}
	timeout := 10 * time.Second
	if parsed, err := time.ParseDuration(strings.TrimSpace(cfg.GetTimeout())); err == nil && parsed > 0 {
		timeout = parsed
	}
	return &WeChatManager{
		appID:  appID,
		secret: secret,
		client: &http.Client{Timeout: timeout},
	}
}

func (m *WeChatManager) Code2Session(ctx context.Context, wxCode string) (*WeChatCode2SessionResponse, error) {
	if m == nil {
		return nil, errors.New("微信配置缺失")
	}
	wxCode = strings.TrimSpace(wxCode)
	if wxCode == "" {
		return nil, errors.New("微信code不能为空")
	}
	query := url.Values{}
	query.Set("appid", m.appID)
	query.Set("secret", m.secret)
	query.Set("js_code", wxCode)
	query.Set("grant_type", "authorization_code")
	endpoint := weChatCode2SessionURL + "?" + query.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	response, err := m.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("调用微信接口失败: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信响应失败: %w", err)
	}
	var result WeChatCode2SessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析微信响应失败: %w", err)
	}
	if result.ErrCode != 0 {
		message := strings.TrimSpace(result.ErrMsg)
		if message == "" {
			message = "微信授权失败"
		}
		return nil, fmt.Errorf("%s", message)
	}
	if strings.TrimSpace(result.OpenID) == "" {
		return nil, errors.New("获取微信OpenID失败")
	}
	return &result, nil
}
