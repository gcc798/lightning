package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type ConfigServiceService struct {
	pb.UnimplementedConfigServiceServer
}

func NewConfigServiceService() *ConfigServiceService {
	return &ConfigServiceService{}
}

func (s *ConfigServiceService) CreateConfig(ctx context.Context, req *pb.ConfigItem) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *ConfigServiceService) PageConfig(ctx context.Context, req *pb.PageConfigRequest) (*pb.PageConfigReply, error) {
	return &pb.PageConfigReply{}, nil
}
func (s *ConfigServiceService) BatchDeleteConfig(ctx context.Context, req *pb.BatchIdsRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *ConfigServiceService) GetConfigByCode(ctx context.Context, req *pb.GetConfigByCodeRequest) (*pb.ConfigListReply, error) {
	return &pb.ConfigListReply{}, nil
}
func (s *ConfigServiceService) GetConfigDataByCode(ctx context.Context, req *pb.GetConfigByCodeRequest) (*pb.ConfigDataReply, error) {
	return &pb.ConfigDataReply{}, nil
}
func (s *ConfigServiceService) UpdateConfig(ctx context.Context, req *pb.UpdateConfigRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *ConfigServiceService) GetConfigById(ctx context.Context, req *pb.IdRequest) (*pb.ConfigItem, error) {
	return &pb.ConfigItem{}, nil
}
func (s *ConfigServiceService) DeleteConfig(ctx context.Context, req *pb.IdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
