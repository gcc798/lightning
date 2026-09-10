package biz

import (
	"context"
	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type DictUsecase struct{ repo *data.DictRepo }

func NewDictUsecase(repo *data.DictRepo) *DictUsecase { return &DictUsecase{repo: repo} }
func (uc *DictUsecase) Create(ctx context.Context, item *v1.DictItem) (*v1.MessageReply, error) {
	return uc.repo.Create(ctx, item)
}
func (uc *DictUsecase) Page(ctx context.Context, req *v1.PageDictRequest) (*v1.PageDictReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *DictUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return uc.repo.BatchDelete(ctx, ids)
}
func (uc *DictUsecase) ByType(ctx context.Context, t string, parentID *int64) (*v1.DictListReply, error) {
	return uc.repo.ByType(ctx, t, parentID)
}
func (uc *DictUsecase) Label(ctx context.Context, t, v string) (*v1.DictLabelReply, error) {
	return uc.repo.Label(ctx, t, v)
}
func (uc *DictUsecase) Update(ctx context.Context, req *v1.UpdateDictRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *DictUsecase) Get(ctx context.Context, id int64) (*v1.DictItem, error) {
	return uc.repo.Get(ctx, id)
}
func (uc *DictUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
