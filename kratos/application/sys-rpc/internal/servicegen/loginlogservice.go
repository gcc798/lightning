package service

import (
	"context"

	pb "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type LoginLogServiceService struct {
	pb.UnimplementedLoginLogServiceServer
}

func NewLoginLogServiceService() *LoginLogServiceService {
	return &LoginLogServiceService{}
}

func (s *LoginLogServiceService) CreateLoginLog(ctx context.Context, req *pb.CreateLoginLogRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *LoginLogServiceService) PageLoginLog(ctx context.Context, req *pb.PageLoginLogRequest) (*pb.PageLogReply, error) {
	return &pb.PageLogReply{}, nil
}
func (s *LoginLogServiceService) BatchDeleteLoginLog(ctx context.Context, req *pb.LogBatchIdsRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *LoginLogServiceService) CleanLoginLog(ctx context.Context, req *pb.CleanLogRequest) (*pb.LogCleanReply, error) {
	return &pb.LogCleanReply{}, nil
}
func (s *LoginLogServiceService) UpdateLoginLog(ctx context.Context, req *pb.UpdateLoginLogRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *LoginLogServiceService) GetLoginLogById(ctx context.Context, req *pb.LogIdRequest) (*pb.LogItem, error) {
	return &pb.LogItem{}, nil
}
func (s *LoginLogServiceService) DeleteLoginLog(ctx context.Context, req *pb.LogIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
