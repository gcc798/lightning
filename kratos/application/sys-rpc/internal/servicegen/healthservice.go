package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type HealthServiceService struct {
	pb.UnimplementedHealthServiceServer
}

func NewHealthServiceService() *HealthServiceService {
	return &HealthServiceService{}
}

func (s *HealthServiceService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{}, nil
}
