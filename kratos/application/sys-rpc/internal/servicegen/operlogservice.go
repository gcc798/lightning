package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type OperLogServiceService struct {
	pb.UnimplementedOperLogServiceServer
}

func NewOperLogServiceService() *OperLogServiceService {
	return &OperLogServiceService{}
}

func (s *OperLogServiceService) CreateOperLog(ctx context.Context, req *pb.CreateOperLogRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *OperLogServiceService) PageOperLog(ctx context.Context, req *pb.PageOperLogRequest) (*pb.PageLogReply, error) {
	return &pb.PageLogReply{}, nil
}
func (s *OperLogServiceService) BatchDeleteOperLog(ctx context.Context, req *pb.LogBatchIdsRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *OperLogServiceService) CleanOperLog(ctx context.Context, req *pb.CleanLogRequest) (*pb.LogCleanReply, error) {
	return &pb.LogCleanReply{}, nil
}
func (s *OperLogServiceService) UpdateOperLog(ctx context.Context, req *pb.UpdateOperLogRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *OperLogServiceService) GetOperLogById(ctx context.Context, req *pb.LogIdRequest) (*pb.LogItem, error) {
	return &pb.LogItem{}, nil
}
func (s *OperLogServiceService) DeleteOperLog(ctx context.Context, req *pb.LogIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
