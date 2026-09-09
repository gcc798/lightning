package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/biz"
)

type AuthServiceService struct {
	v1.UnimplementedAuthServiceServer

	uc *biz.AuthUsecase
}

func NewAuthServiceService(uc *biz.AuthUsecase) *AuthServiceService {
	return &AuthServiceService{uc: uc}
}

func (s *AuthServiceService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return s.uc.Login(ctx, req)
}

func (s *AuthServiceService) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.MessageReply, error) {
	return s.uc.Logout(ctx)
}

func (s *AuthServiceService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	return s.uc.RefreshToken(ctx, req)
}
