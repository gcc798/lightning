package svc

import (
	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	_ "github.com/lib/pq"
)

type SMSProvider interface {
	SendSMS(phone, code string) error
}

type EmailProvider interface {
	SendEmail(email, code string) error
}

type consoleSMSProvider struct{}

func (p *consoleSMSProvider) SendSMS(phone, code string) error {
	logx.Infof("[验证码] 短信验证码发送至 %s: %s", phone, code)
	return nil
}

type consoleEmailProvider struct{}

func (p *consoleEmailProvider) SendEmail(email, code string) error {
	logx.Infof("[验证码] 邮箱验证码发送至 %s: %s", email, code)
	return nil
}

type ServiceContext struct {
	Config         config.Config
	DB             sqlx.SqlConn
	Redis          *redis.Client
	SMSProvider    SMSProvider
	EmailProvider  EmailProvider
	StorageManager *StorageManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewSqlConn("postgres", c.Postgres.Dsn)
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.CacheRedis.Addr,
		Password: c.CacheRedis.Password,
		DB:       c.CacheRedis.Db,
	})

	return &ServiceContext{
		Config:         c,
		DB:             db,
		Redis:          rdb,
		SMSProvider:    &consoleSMSProvider{},
		EmailProvider:  &consoleEmailProvider{},
		StorageManager: NewStorageManager(db),
	}
}
