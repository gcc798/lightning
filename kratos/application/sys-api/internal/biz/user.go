package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type UserUsecase struct{ repo *data.UserRepo }

func NewUserUsecase(repo *data.UserRepo) *UserUsecase { return &UserUsecase{repo: repo} }
func (uc *UserUsecase) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.MessageReply, error) {
	return uc.repo.Create(ctx, req)
}
func (uc *UserUsecase) Update(ctx context.Context, req *v1.UpdateUserRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *UserUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
func (uc *UserUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return uc.repo.BatchDelete(ctx, ids)
}
func (uc *UserUsecase) GetByID(ctx context.Context, id int64) (*v1.UserItem, error) {
	return uc.repo.GetByID(ctx, id)
}
func (uc *UserUsecase) Import(ctx context.Context, users []*v1.CreateUserRequest) (*v1.ImportUsersReply, error) {
	return uc.repo.Import(ctx, users)
}
func (uc *UserUsecase) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) (*v1.MessageReply, error) {
	return uc.repo.ResetPassword(ctx, req)
}
func (uc *UserUsecase) Page(ctx context.Context, req *v1.PageUsersRequest) (*v1.PageUsersReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *UserUsecase) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.MessageReply, error) {
	return uc.repo.ChangePassword(ctx, req)
}
