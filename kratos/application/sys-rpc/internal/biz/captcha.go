package biz

import (
	"context"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
	"google.golang.org/protobuf/types/known/structpb"
)

type CaptchaUsecase struct {
	res *data.Resources
}

func NewCaptchaUsecase(res *data.Resources) *CaptchaUsecase {
	return &CaptchaUsecase{res: res}
}

func (uc *CaptchaUsecase) Image(ctx context.Context) (*v1.CaptchaReply, error) {
	captchaID, image, expireAt, err := uc.res.CreateImageCaptcha(ctx)
	if err != nil {
		return nil, err
	}
	return captchaReply(captchaID, "image", map[string]any{"image": image}, expireAt)
}

func (uc *CaptchaUsecase) SMS(ctx context.Context, phone string) (*v1.CaptchaReply, error) {
	captchaID, expireAt, err := uc.res.CreateSMSCaptcha(ctx, phone)
	if err != nil {
		return nil, err
	}
	return captchaReply(captchaID, "sms", map[string]any{"target": phone}, expireAt)
}

func (uc *CaptchaUsecase) Email(ctx context.Context, email string) (*v1.CaptchaReply, error) {
	captchaID, expireAt, err := uc.res.CreateEmailCaptcha(ctx, email)
	if err != nil {
		return nil, err
	}
	return captchaReply(captchaID, "email", map[string]any{"target": email}, expireAt)
}

func (uc *CaptchaUsecase) Enabled(context.Context) (*v1.GetEnabledTypesReply, error) {
	return &v1.GetEnabledTypesReply{Items: []string{"image", "sms", "email"}}, nil
}

func captchaReply(id, captchaType string, payload map[string]any, expireAt time.Time) (*v1.CaptchaReply, error) {
	data, err := structpb.NewStruct(payload)
	if err != nil {
		return nil, err
	}
	return &v1.CaptchaReply{
		Id:       id,
		Type:     captchaType,
		Data:     data,
		ExpireAt: expireAt.Format(time.RFC3339),
	}, nil
}
