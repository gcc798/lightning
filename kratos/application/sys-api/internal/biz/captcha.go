package biz

import (
	"context"
	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type CaptchaUsecase struct{ repo *data.CaptchaRepo }

func NewCaptchaUsecase(repo *data.CaptchaRepo) *CaptchaUsecase { return &CaptchaUsecase{repo: repo} }
func (uc *CaptchaUsecase) Image(ctx context.Context) (*v1.CaptchaReply, error) {
	return uc.repo.Image(ctx)
}
func (uc *CaptchaUsecase) SMS(ctx context.Context, phone string) (*v1.CaptchaReply, error) {
	return uc.repo.SMS(ctx, phone)
}
func (uc *CaptchaUsecase) Email(ctx context.Context, email string) (*v1.CaptchaReply, error) {
	return uc.repo.Email(ctx, email)
}
func (uc *CaptchaUsecase) Enabled(ctx context.Context) (*v1.GetEnabledTypesReply, error) {
	return uc.repo.Enabled(ctx)
}
