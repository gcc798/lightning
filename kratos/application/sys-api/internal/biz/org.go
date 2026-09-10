package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type OrgUsecase struct{ repo *data.OrgRepo }

func NewOrgUsecase(repo *data.OrgRepo) *OrgUsecase { return &OrgUsecase{repo: repo} }
func (uc *OrgUsecase) Create(ctx context.Context, req *v1.CreateOrgRequest) (*v1.OrgIdReply, error) {
	return uc.repo.Create(ctx, req)
}
func (uc *OrgUsecase) Update(ctx context.Context, req *v1.UpdateOrgRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *OrgUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
func (uc *OrgUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return uc.repo.BatchDelete(ctx, ids)
}
func (uc *OrgUsecase) GetByID(ctx context.Context, id int64) (*v1.OrgItem, error) {
	return uc.repo.GetByID(ctx, id)
}
func (uc *OrgUsecase) Tree(ctx context.Context) (*v1.OrgTreeReply, error) { return uc.repo.Tree(ctx) }
func (uc *OrgUsecase) Page(ctx context.Context, req *v1.PageOrgsRequest) (*v1.PageOrgsReply, error) {
	return uc.repo.Page(ctx, req)
}
