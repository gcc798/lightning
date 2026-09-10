package sms

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// SmsCodeLength 定义业务常量。
	SmsCodeLength = 6
	// SmsCodeExpire 定义业务常量。
	SmsCodeExpire = 2 * time.Minute
	// SmsCodeRedisKey 定义业务常量。
	SmsCodeRedisKey = "global:captcha_codes:"
)

// Config 定义业务数据结构。
type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// Manager 定义业务数据结构。
type Manager struct {
	config Config
	client *dysmsapi.Client
	redis  *redis.Client
	logger logging.Logger
}

// NewManager 创建组件实例。
func NewManager(config Config, redis *redis.Client, logger logging.Logger) (*Manager, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create sms client: %v", err)
	}
	return &Manager{config: config, client: client, redis: redis, logger: logger}, nil
}

// SendVerificationCode 执行业务逻辑。
func (s *Manager) SendVerificationCode(ctx context.Context, phonenumber string) (string, error) {
	if phonenumber == "" {
		return "", fmt.Errorf("手机号不能为空")
	}
	code := generateCode(SmsCodeLength)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phonenumber
	request.SignName = s.config.SignName
	request.TemplateCode = s.config.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)
	response, err := s.client.SendSms(request)
	if err != nil {
		s.logger.Error("failed to send sms", zap.Error(err))
		return "", fmt.Errorf("发送短信失败")
	}
	if response.Code != "OK" {
		s.logger.Error("sms send failed", zap.String("code", response.Code), zap.String("message", response.Message))
		return "", fmt.Errorf("发送短信失败: %s", response.Message)
	}
	s.logger.Info("sms sent successfully", zap.String("phonenumber", phonenumber), zap.String("bizId", response.BizId))
	key := SmsCodeRedisKey + phonenumber
	if err := s.redis.Set(ctx, key, code, SmsCodeExpire).Err(); err != nil {
		s.logger.Error("failed to save sms code to redis", zap.Error(err))
		return "", fmt.Errorf("保存验证码失败")
	}
	return code, nil
}

// Send 发送自定义短信
func (s *Manager) Send(phonenumber string, templateCode string, params map[string]string) error {
	if phonenumber == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if templateCode == "" {
		return fmt.Errorf("模板编码不能为空")
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phonenumber
	request.SignName = s.config.SignName
	request.TemplateCode = templateCode

	// 构建模板参数
	if len(params) > 0 {
		templateParam := "{"
		first := true
		for k, v := range params {
			if !first {
				templateParam += ","
			}
			templateParam += fmt.Sprintf(`"%s":"%s"`, k, v)
			first = false
		}
		templateParam += "}"
		request.TemplateParam = templateParam
	}

	response, err := s.client.SendSms(request)
	if err != nil {
		s.logger.Error("failed to send sms", zap.Error(err))
		return fmt.Errorf("发送短信失败")
	}
	if response.Code != "OK" {
		s.logger.Error("sms send failed", zap.String("code", response.Code), zap.String("message", response.Message))
		return fmt.Errorf("发送短信失败: %s", response.Message)
	}

	s.logger.Info("sms sent successfully", zap.String("phonenumber", phonenumber), zap.String("bizId", response.BizId))
	return nil
}

func generateCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}
