package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/biz"
)

type ConfigServiceService struct {
	v1.UnimplementedConfigServiceServer
	uc *biz.ConfigUsecase
}

func NewConfigServiceService(uc *biz.ConfigUsecase) *ConfigServiceService {
	return &ConfigServiceService{uc: uc}
}
func (s *ConfigServiceService) CreateConfig(ctx context.Context, req *v1.ConfigItem) (*v1.MessageReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *ConfigServiceService) PageConfig(ctx context.Context, req *v1.PageConfigRequest) (*v1.PageConfigReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *ConfigServiceService) BatchDeleteConfig(ctx context.Context, req *v1.BatchIdsRequest) (*v1.MessageReply, error) {
	return s.uc.BatchDelete(ctx, req.GetIds())
}
func (s *ConfigServiceService) GetConfigByCode(ctx context.Context, req *v1.GetConfigByCodeRequest) (*v1.ConfigListReply, error) {
	return s.uc.ByCode(ctx, req.GetCode())
}
func (s *ConfigServiceService) GetConfigDataByCode(ctx context.Context, req *v1.GetConfigByCodeRequest) (*v1.ConfigDataReply, error) {
	return s.uc.Data(ctx, req.GetCode())
}
func (s *ConfigServiceService) UpdateConfig(ctx context.Context, req *v1.UpdateConfigRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *ConfigServiceService) GetConfigById(ctx context.Context, req *v1.IdRequest) (*v1.ConfigItem, error) {
	return s.uc.Get(ctx, req.GetId())
}
func (s *ConfigServiceService) DeleteConfig(ctx context.Context, req *v1.IdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
