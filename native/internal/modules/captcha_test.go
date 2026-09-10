package modules

import (
	"context"
	"sync"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/platform/captcha"
	"github.com/gcc798/microservice-kit/internal/platform/redislock"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type mapSource struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func (s *mapSource) Load(_ context.Context, code string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]byte(nil), s.data[code]...), nil
}

type moduleTestContainer struct {
	redis   *redis.Client
	store   *runtimeconfig.Store
	modules map[string]Module
}

type noopLogger struct{ value *zap.Logger }

func newNoopLogger() logging.Logger            { return &noopLogger{value: zap.NewNop()} }
func (l *noopLogger) Get() *zap.Logger         { return l.value }
func (*noopLogger) Debug(string, ...zap.Field) {}
func (*noopLogger) Info(string, ...zap.Field)  {}
func (*noopLogger) Warn(string, ...zap.Field)  {}
func (*noopLogger) Error(string, ...zap.Field) {}
func (*noopLogger) Fatal(string, ...zap.Field) {}
func (l *noopLogger) With(fields ...zap.Field) logging.Logger {
	return &noopLogger{value: l.value.With(fields...)}
}

func (*moduleTestContainer) GetDB() *gorm.DB                          { return nil }
func (c *moduleTestContainer) GetRedis() *redis.Client                { return c.redis }
func (*moduleTestContainer) GetLogger() logging.Logger                { return newNoopLogger() }
func (c *moduleTestContainer) GetRuntimeConfig() *runtimeconfig.Store { return c.store }
func (c *moduleTestContainer) GetModule(name string) Module           { return c.modules[name] }

func TestCaptchaModuleReadsChangedRedisConfiguration(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	source := &mapSource{data: map[string][]byte{
		runtimeconfig.CodeSMS:     []byte(`{"enabled":false,"accessKeyId":"","accessKeySecret":"","signName":"","templateCode":""}`),
		runtimeconfig.CodeEmail:   []byte(`{"enabled":false,"host":"","port":0,"username":"","password":"","from":""}`),
		runtimeconfig.CodeCaptcha: []byte(`{"image":{"enabled":false,"length":4,"width":120,"height":40,"expire":300},"sms":{"enabled":false,"length":6,"expire":300,"template":"SMS_CODE_TEMPLATE","provider":"aliyun"},"email":{"enabled":false,"length":6,"expire":300,"template":"验证码：%s"}}`),
	}}
	cont := &moduleTestContainer{redis: client, modules: make(map[string]Module)}
	cont.store = runtimeconfig.NewStore(client, source, redislock.New(client))
	smsModule, emailModule, captchaModule := NewSMSModule(), NewEmailModule(), NewCaptchaModule()
	cont.modules[SMSName], cont.modules[EmailName], cont.modules[CaptchaName] = smsModule, emailModule, captchaModule
	ctx := context.Background()
	if err := smsModule.Init(ctx, cont); err != nil {
		t.Fatal(err)
	}
	if err := emailModule.Init(ctx, cont); err != nil {
		t.Fatal(err)
	}
	if err := captchaModule.Init(ctx, cont); err != nil {
		t.Fatal(err)
	}
	types, err := captchaModule.GetEnabledTypes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 0 {
		t.Fatalf("enabled types = %v, want none", types)
	}

	updated := []byte(`{"image":{"enabled":true,"length":4,"width":120,"height":40,"expire":300},"sms":{"enabled":false,"length":6,"expire":300,"template":"SMS_CODE_TEMPLATE","provider":"aliyun"},"email":{"enabled":false,"length":6,"expire":300,"template":"验证码：%s"}}`)
	if err := client.Set(ctx, runtimeconfig.CacheKey(runtimeconfig.CodeCaptcha), updated, 0).Err(); err != nil {
		t.Fatal(err)
	}
	types, err = captchaModule.GetEnabledTypes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 1 || types[0] != captcha.CaptchaTypeImage {
		t.Fatalf("enabled types = %v, want image", types)
	}
}
