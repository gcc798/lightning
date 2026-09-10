package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/data"
)

type UserUsecase struct {
	res *data.Resources
}

func NewUserUsecase(res *data.Resources) *UserUsecase { return &UserUsecase{res: res} }

func (uc *UserUsecase) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.MessageReply, error) {
	if err := uc.res.CreateUser(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *UserUsecase) Update(ctx context.Context, req *v1.UpdateUserRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateUser(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *UserUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteUsers(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *UserUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteUsers(ctx, ids...); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *UserUsecase) Import(ctx context.Context, users []*v1.CreateUserRequest) (*v1.ImportUsersReply, error) {
	return uc.res.ImportUsers(ctx, users)
}

func (uc *UserUsecase) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) (*v1.MessageReply, error) {
	if err := uc.res.ResetPassword(ctx, req.GetUserId(), req.GetNewPassword()); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.MessageReply, error) {
	if err := uc.res.ChangePassword(ctx, currentUserID(ctx), req.GetOldPassword(), req.GetNewPassword()); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *UserUsecase) GetByID(ctx context.Context, id int64) (*v1.UserItem, error) {
	item, err := uc.res.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "用户不存在")
}

func (uc *UserUsecase) Page(ctx context.Context, req *v1.PageUsersRequest) (*v1.PageUsersReply, error) {
	return uc.res.PageUsers(ctx, req)
}
