package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type LoginLogUsecase struct{ repo *data.LoginLogRepo }
type OperLogUsecase struct{ repo *data.OperLogRepo }

func NewLoginLogUsecase(repo *data.LoginLogRepo) *LoginLogUsecase {
	return &LoginLogUsecase{repo: repo}
}
func NewOperLogUsecase(repo *data.OperLogRepo) *OperLogUsecase { return &OperLogUsecase{repo: repo} }

func (uc *LoginLogUsecase) Create(ctx context.Context, item *v1.CreateLoginLogRequest) (*v1.MessageReply, error) {
	return uc.repo.Create(ctx, item)
}
func (uc *LoginLogUsecase) Page(ctx context.Context, req *v1.PageLoginLogRequest) (*v1.PageLogReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *LoginLogUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return uc.repo.BatchDelete(ctx, ids)
}
func (uc *LoginLogUsecase) Clean(ctx context.Context, days int32) (*v1.LogCleanReply, error) {
	return uc.repo.Clean(ctx, days)
}
func (uc *LoginLogUsecase) Update(ctx context.Context, req *v1.UpdateLoginLogRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *LoginLogUsecase) Get(ctx context.Context, id int64) (*v1.LogItem, error) {
	return uc.repo.Get(ctx, id)
}
func (uc *LoginLogUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}

func (uc *OperLogUsecase) Create(ctx context.Context, item *v1.CreateOperLogRequest) (*v1.MessageReply, error) {
	return uc.repo.Create(ctx, item)
}
func (uc *OperLogUsecase) Page(ctx context.Context, req *v1.PageOperLogRequest) (*v1.PageLogReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *OperLogUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return uc.repo.BatchDelete(ctx, ids)
}
func (uc *OperLogUsecase) Clean(ctx context.Context, days int32) (*v1.LogCleanReply, error) {
	return uc.repo.Clean(ctx, days)
}
func (uc *OperLogUsecase) Update(ctx context.Context, req *v1.UpdateOperLogRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *OperLogUsecase) Get(ctx context.Context, id int64) (*v1.LogItem, error) {
	return uc.repo.Get(ctx, id)
}
func (uc *OperLogUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
