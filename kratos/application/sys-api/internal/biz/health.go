package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type HealthUsecase struct {
	repo *data.HealthRepo
}

func NewHealthUsecase(repo *data.HealthRepo) *HealthUsecase {
	return &HealthUsecase{repo: repo}
}

func (uc *HealthUsecase) Ping(ctx context.Context, name string) (*v1.PingReply, error) {
	return uc.repo.Ping(ctx, name)
}
