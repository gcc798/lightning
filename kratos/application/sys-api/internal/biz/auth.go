package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type AuthUsecase struct {
	repo *data.AuthRepo
}

func NewAuthUsecase(repo *data.AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (uc *AuthUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return uc.repo.Login(ctx, req)
}

func (uc *AuthUsecase) Logout(ctx context.Context) (*v1.MessageReply, error) {
	return uc.repo.Logout(ctx)
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	return uc.repo.RefreshToken(ctx, req)
}
