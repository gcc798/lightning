package biz

import (
	"context"
	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/data"
)

type StorageEnvUsecase struct{ repo *data.StorageEnvRepo }

func NewStorageEnvUsecase(repo *data.StorageEnvRepo) *StorageEnvUsecase {
	return &StorageEnvUsecase{repo: repo}
}
func (uc *StorageEnvUsecase) Create(ctx context.Context, item *v1.StorageEnvItem) (*v1.StorageEnvItem, error) {
	return uc.repo.Create(ctx, item)
}
func (uc *StorageEnvUsecase) Page(ctx context.Context, req *v1.PageStorageEnvRequest) (*v1.PageStorageEnvReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *StorageEnvUsecase) Default(ctx context.Context) (*v1.StorageEnvItem, error) {
	return uc.repo.Default(ctx)
}
func (uc *StorageEnvUsecase) SetDefault(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.SetDefault(ctx, id)
}
func (uc *StorageEnvUsecase) Update(ctx context.Context, req *v1.UpdateStorageEnvRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *StorageEnvUsecase) Get(ctx context.Context, id int64) (*v1.StorageEnvItem, error) {
	return uc.repo.Get(ctx, id)
}
func (uc *StorageEnvUsecase) Test(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Test(ctx, id)
}
func (uc *StorageEnvUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
