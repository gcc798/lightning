package sysservicelogic

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"
	"github.com/google/uuid"
	base64Captcha "github.com/mojocn/base64Captcha"
)

const (
	captchaImageKeyPrefix = "captcha:image:"
	captchaSmsKeyPrefix   = "captcha:sms:"
	captchaEmailKeyPrefix = "captcha:email:"
	captchaRatePrefix     = "captcha:ratelimit:"
)

func captchaEnabledTypes(svcCtx *svc.ServiceContext) []string {
	types := make([]string, 0, 3)
	if svcCtx.Config.Captcha.Image.Enabled {
		types = append(types, "image")
	}
	if svcCtx.Config.Captcha.Sms.Enabled {
		types = append(types, "sms")
	}
	if svcCtx.Config.Captcha.Email.Enabled {
		types = append(types, "email")
	}
	return types
}

func verifySmsCaptcha(ctx context.Context, svcCtx *svc.ServiceContext, uuidValue, phone, code string) error {
	if uuidValue == "" || phone == "" || code == "" {
		return fmt.Errorf("手机号、验证码ID和验证码不能为空")
	}
	key := captchaSmsKeyPrefix + uuidValue
	fields, err := svcCtx.Redis.HGetAll(ctx, key).Result()
	if err != nil || len(fields) == 0 {
		return fmt.Errorf("验证码已过期或不存在")
	}
	_ = svcCtx.Redis.Del(ctx, key).Err()
	if fields["phone"] != phone {
		return fmt.Errorf("手机号与验证码不匹配")
	}
	if fields["code"] != code {
		return fmt.Errorf("验证码错误")
	}
	return nil
}

func verifyEmailCaptcha(ctx context.Context, svcCtx *svc.ServiceContext, uuidValue, email, code string) error {
	if uuidValue == "" || email == "" || code == "" {
		return fmt.Errorf("邮箱、验证码ID和验证码不能为空")
	}
	key := captchaEmailKeyPrefix + uuidValue
	fields, err := svcCtx.Redis.HGetAll(ctx, key).Result()
	if err != nil || len(fields) == 0 {
		return fmt.Errorf("验证码已过期或不存在")
	}
	_ = svcCtx.Redis.Del(ctx, key).Err()
	if fields["email"] != email {
		return fmt.Errorf("邮箱与验证码不匹配")
	}
	if fields["code"] != code {
		return fmt.Errorf("验证码错误")
	}
	return nil
}

func verifyImageCaptcha(ctx context.Context, svcCtx *svc.ServiceContext, uuidValue, code string) error {
	if !svcCtx.Config.Captcha.Image.Enabled {
		return nil
	}
	if uuidValue == "" || code == "" {
		return fmt.Errorf("验证码不能为空")
	}
	key := captchaImageKeyPrefix + uuidValue
	saved, err := svcCtx.Redis.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}
	_ = svcCtx.Redis.Del(ctx, key).Err()
	if !strings.EqualFold(saved, code) {
		return fmt.Errorf("验证码错误")
	}
	return nil
}

func newCaptchaData(captchaType, id string, data interface{}, expireAt time.Time) (*pb.CaptchaDataResp, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &pb.CaptchaDataResp{
		Id:       id,
		Type:     captchaType,
		DataJson: string(jsonData),
		ExpireAt: expireAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func generateImageCaptcha(ctx context.Context, svcCtx *svc.ServiceContext) (*pb.CaptchaDataResp, error) {
	if !svcCtx.Config.Captcha.Image.Enabled {
		return nil, fmt.Errorf("验证码类型 image 未启用")
	}
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	_, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("生成验证码失败: %w", err)
	}
	id := uuid.NewString()
	expire := 5 * time.Minute
	if err := svcCtx.Redis.Set(ctx, captchaImageKeyPrefix+id, answer, expire).Err(); err != nil {
		return nil, err
	}
	return newCaptchaData("image", id, b64s, time.Now().Add(expire))
}

func generateSmsCaptcha(ctx context.Context, svcCtx *svc.ServiceContext, phone string) (*pb.CaptchaDataResp, error) {
	if !svcCtx.Config.Captcha.Sms.Enabled {
		return nil, fmt.Errorf("验证码类型 sms 未启用")
	}
	if phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}
	if !isValidPhone(phone) {
		return nil, fmt.Errorf("手机号格式错误")
	}
	rateKey := captchaRatePrefix + "sms:" + phone
	exists, err := svcCtx.Redis.Exists(ctx, rateKey).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("发送过于频繁，请稍后再试")
	}
	id := uuid.NewString()
	code := randomDigits(6)
	expire := 5 * time.Minute
	if err := svcCtx.Redis.HSet(ctx, captchaSmsKeyPrefix+id, map[string]interface{}{"code": code, "phone": phone}).Err(); err != nil {
		return nil, err
	}
	if err := svcCtx.Redis.Expire(ctx, captchaSmsKeyPrefix+id, expire).Err(); err != nil {
		return nil, err
	}
	_ = svcCtx.Redis.Set(ctx, rateKey, "1", time.Minute).Err()
	if err := svcCtx.SMSProvider.SendSMS(phone, code); err != nil {
		return nil, fmt.Errorf("发送短信验证码失败: %w", err)
	}
	return newCaptchaData("sms", id, map[string]interface{}{"phone": maskPhone(phone)}, time.Now().Add(expire))
}

func generateEmailCaptcha(ctx context.Context, svcCtx *svc.ServiceContext, email string) (*pb.CaptchaDataResp, error) {
	if !svcCtx.Config.Captcha.Email.Enabled {
		return nil, fmt.Errorf("验证码类型 email 未启用")
	}
	if email == "" {
		return nil, fmt.Errorf("邮箱不能为空")
	}
	if !isValidEmail(email) {
		return nil, fmt.Errorf("邮箱格式错误")
	}
	rateKey := captchaRatePrefix + "email:" + email
	exists, err := svcCtx.Redis.Exists(ctx, rateKey).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("发送过于频繁，请稍后再试")
	}
	id := uuid.NewString()
	code := randomDigits(6)
	expire := 5 * time.Minute
	if err := svcCtx.Redis.HSet(ctx, captchaEmailKeyPrefix+id, map[string]interface{}{"code": code, "email": email}).Err(); err != nil {
		return nil, err
	}
	if err := svcCtx.Redis.Expire(ctx, captchaEmailKeyPrefix+id, expire).Err(); err != nil {
		return nil, err
	}
	_ = svcCtx.Redis.Set(ctx, rateKey, "1", time.Minute).Err()
	if err := svcCtx.EmailProvider.SendEmail(email, code); err != nil {
		return nil, fmt.Errorf("发送邮箱验证码失败: %w", err)
	}
	return newCaptchaData("email", id, map[string]interface{}{"email": maskEmail(email)}, time.Now().Add(expire))
}

func randomDigits(length int) string {
	const digits = "0123456789"
	return randomFromCharset(length, digits)
}

func randomFromCharset(length int, charset string) string {
	out := make([]byte, length)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			out[i] = charset[0]
			continue
		}
		out[i] = charset[n.Int64()]
	}
	return string(out)
}

func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) < 2 {
		return email
	}
	return parts[0][:1] + "***@" + parts[1]
}

func isValidPhone(phone string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

func isValidEmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}
