package biz

import (
	"context"
	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type ConfigUsecase struct{ repo *data.ConfigRepo }

func NewConfigUsecase(repo *data.ConfigRepo) *ConfigUsecase { return &ConfigUsecase{repo: repo} }
func (uc *ConfigUsecase) Create(ctx context.Context, item *v1.ConfigItem) (*v1.MessageReply, error) {
	return uc.repo.Create(ctx, item)
}
func (uc *ConfigUsecase) Page(ctx context.Context, req *v1.PageConfigRequest) (*v1.PageConfigReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *ConfigUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return uc.repo.BatchDelete(ctx, ids)
}
func (uc *ConfigUsecase) ByCode(ctx context.Context, code string) (*v1.ConfigListReply, error) {
	return uc.repo.ByCode(ctx, code)
}
func (uc *ConfigUsecase) Data(ctx context.Context, code string) (*v1.ConfigDataReply, error) {
	return uc.repo.Data(ctx, code)
}
func (uc *ConfigUsecase) Update(ctx context.Context, req *v1.UpdateConfigRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *ConfigUsecase) Get(ctx context.Context, id int64) (*v1.ConfigItem, error) {
	return uc.repo.Get(ctx, id)
}
func (uc *ConfigUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
