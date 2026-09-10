package email

import (
	"fmt"
	"time"

	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const (
	// EmailCodeLength 定义业务常量。
	EmailCodeLength = 6
	// EmailCodeExpire 定义业务常量。
	EmailCodeExpire = 5 * time.Minute
	// EmailCodeRedisKey 定义业务常量。
	EmailCodeRedisKey = "global:captcha_codes:email:"
)

// Config 定义业务数据结构。
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Manager 定义业务数据结构。
type Manager struct {
	config Config
	redis  *redis.Client
	logger logging.Logger
}

// NewManager 创建组件实例。
func NewManager(config Config, redis *redis.Client, logger logging.Logger) (*Manager, error) {
	return &Manager{
		config: config,
		redis:  redis,
		logger: logger,
	}, nil
}

// Send 执行业务逻辑。
func (m *Manager) Send(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.config.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(m.config.Host, m.config.Port, m.config.Username, m.config.Password)
	if err := dialer.DialAndSend(msg); err != nil {
		m.logger.Error("failed to send email", zap.Error(err))
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	m.logger.Info("email sent successfully", zap.String("to", to))
	return nil
}

// SendWithTemplate 执行业务逻辑。
func (m *Manager) SendWithTemplate(to string, code string, template string) error {
	subject := "验证码"
	body := fmt.Sprintf(template, code)
	return m.Send(to, subject, body)
}
