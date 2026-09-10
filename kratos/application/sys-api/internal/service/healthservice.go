package service

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/biz"
)

type HealthServiceService struct {
	v1.UnimplementedHealthServiceServer

	uc *biz.HealthUsecase
}

func NewHealthServiceService(uc *biz.HealthUsecase) *HealthServiceService {
	return &HealthServiceService{uc: uc}
}

func (s *HealthServiceService) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingReply, error) {
	return s.uc.Ping(ctx, req.GetName())
}
