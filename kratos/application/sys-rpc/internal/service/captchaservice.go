package service

import (
	"context"
	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/biz"
	"strings"
)

type CaptchaServiceService struct {
	v1.UnimplementedCaptchaServiceServer
	uc *biz.CaptchaUsecase
}

func NewCaptchaServiceService(uc *biz.CaptchaUsecase) *CaptchaServiceService {
	return &CaptchaServiceService{uc: uc}
}
func (s *CaptchaServiceService) GenerateImageCaptcha(ctx context.Context, req *v1.GenerateImageCaptchaRequest) (*v1.CaptchaReply, error) {
	return s.uc.Image(ctx)
}
func (s *CaptchaServiceService) SendSMSCaptcha(ctx context.Context, req *v1.SendSMSCaptchaRequest) (*v1.CaptchaReply, error) {
	phone := strings.TrimSpace(req.GetPhonenumber())
	if phone == "" {
		phone = strings.TrimSpace(req.GetPhone())
	}
	return s.uc.SMS(ctx, phone)
}
func (s *CaptchaServiceService) SendEmailCaptcha(ctx context.Context, req *v1.SendEmailCaptchaRequest) (*v1.CaptchaReply, error) {
	return s.uc.Email(ctx, req.GetEmail())
}
func (s *CaptchaServiceService) GetEnabledTypes(ctx context.Context, req *v1.GetEnabledTypesRequest) (*v1.GetEnabledTypesReply, error) {
	return s.uc.Enabled(ctx)
}
