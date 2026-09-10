package service

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/biz"
)

type StorageEnvServiceService struct {
	v1.UnimplementedStorageEnvServiceServer
	uc *biz.StorageEnvUsecase
}

func NewStorageEnvServiceService(uc *biz.StorageEnvUsecase) *StorageEnvServiceService {
	return &StorageEnvServiceService{uc: uc}
}
func (s *StorageEnvServiceService) CreateStorageEnv(ctx context.Context, req *v1.StorageEnvItem) (*v1.StorageEnvItem, error) {
	return s.uc.Create(ctx, req)
}
func (s *StorageEnvServiceService) PageStorageEnv(ctx context.Context, req *v1.PageStorageEnvRequest) (*v1.PageStorageEnvReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *StorageEnvServiceService) GetDefaultStorageEnv(ctx context.Context, req *v1.StorageEmpty) (*v1.StorageEnvItem, error) {
	return s.uc.Default(ctx)
}
func (s *StorageEnvServiceService) SetDefaultStorageEnv(ctx context.Context, req *v1.SetDefaultStorageEnvRequest) (*v1.MessageReply, error) {
	return s.uc.SetDefault(ctx, req.GetId())
}
func (s *StorageEnvServiceService) UpdateStorageEnv(ctx context.Context, req *v1.UpdateStorageEnvRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *StorageEnvServiceService) GetStorageEnv(ctx context.Context, req *v1.StorageIdRequest) (*v1.StorageEnvItem, error) {
	return s.uc.Get(ctx, req.GetId())
}
func (s *StorageEnvServiceService) TestStorageEnvConnection(ctx context.Context, req *v1.StorageIdRequest) (*v1.MessageReply, error) {
	return s.uc.Test(ctx, req.GetId())
}
func (s *StorageEnvServiceService) DeleteStorageEnv(ctx context.Context, req *v1.StorageIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
