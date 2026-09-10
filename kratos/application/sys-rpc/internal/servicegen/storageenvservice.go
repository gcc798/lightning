package service

import (
	"context"

	pb "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type StorageEnvServiceService struct {
	pb.UnimplementedStorageEnvServiceServer
}

func NewStorageEnvServiceService() *StorageEnvServiceService {
	return &StorageEnvServiceService{}
}

func (s *StorageEnvServiceService) CreateStorageEnv(ctx context.Context, req *pb.StorageEnvItem) (*pb.StorageEnvItem, error) {
	return &pb.StorageEnvItem{}, nil
}
func (s *StorageEnvServiceService) PageStorageEnv(ctx context.Context, req *pb.PageStorageEnvRequest) (*pb.PageStorageEnvReply, error) {
	return &pb.PageStorageEnvReply{}, nil
}
func (s *StorageEnvServiceService) GetDefaultStorageEnv(ctx context.Context, req *pb.StorageEmpty) (*pb.StorageEnvItem, error) {
	return &pb.StorageEnvItem{}, nil
}
func (s *StorageEnvServiceService) SetDefaultStorageEnv(ctx context.Context, req *pb.SetDefaultStorageEnvRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *StorageEnvServiceService) UpdateStorageEnv(ctx context.Context, req *pb.UpdateStorageEnvRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *StorageEnvServiceService) GetStorageEnv(ctx context.Context, req *pb.StorageIdRequest) (*pb.StorageEnvItem, error) {
	return &pb.StorageEnvItem{}, nil
}
func (s *StorageEnvServiceService) TestStorageEnvConnection(ctx context.Context, req *pb.StorageIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *StorageEnvServiceService) DeleteStorageEnv(ctx context.Context, req *pb.StorageIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
