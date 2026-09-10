package data

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type CaptchaRepo struct {
	client v1.CaptchaServiceClient
}

func NewCaptchaRepo(clients *RPCClientSet) *CaptchaRepo {
	return &CaptchaRepo{client: clients.Captcha}
}

func (r *CaptchaRepo) Image(ctx context.Context) (*v1.CaptchaReply, error) {
	return r.client.GenerateImageCaptcha(ctx, &v1.GenerateImageCaptchaRequest{})
}
func (r *CaptchaRepo) SMS(ctx context.Context, phone string) (*v1.CaptchaReply, error) {
	return r.client.SendSMSCaptcha(ctx, &v1.SendSMSCaptchaRequest{Phonenumber: phone, Phone: phone})
}
func (r *CaptchaRepo) Email(ctx context.Context, email string) (*v1.CaptchaReply, error) {
	return r.client.SendEmailCaptcha(ctx, &v1.SendEmailCaptchaRequest{Email: email})
}
func (r *CaptchaRepo) Enabled(ctx context.Context) (*v1.GetEnabledTypesReply, error) {
	return r.client.GetEnabledTypes(ctx, &v1.GetEnabledTypesRequest{})
}
