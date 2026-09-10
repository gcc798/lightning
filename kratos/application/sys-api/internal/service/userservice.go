package service

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/biz"
)

type UserServiceService struct {
	v1.UnimplementedUserServiceServer
	uc *biz.UserUsecase
}

func NewUserServiceService(uc *biz.UserUsecase) *UserServiceService {
	return &UserServiceService{uc: uc}
}
func (s *UserServiceService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.MessageReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *UserServiceService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *UserServiceService) DeleteUser(ctx context.Context, req *v1.UserIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
func (s *UserServiceService) BatchDeleteUser(ctx context.Context, req *v1.BatchDeleteUsersRequest) (*v1.MessageReply, error) {
	return s.uc.BatchDelete(ctx, req.GetIds())
}
func (s *UserServiceService) GetUserById(ctx context.Context, req *v1.UserIdRequest) (*v1.UserItem, error) {
	return s.uc.GetByID(ctx, req.GetId())
}
func (s *UserServiceService) ImportUsers(ctx context.Context, req *v1.BatchImportUsersRequest) (*v1.ImportUsersReply, error) {
	return s.uc.Import(ctx, req.GetUsers())
}
func (s *UserServiceService) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) (*v1.MessageReply, error) {
	return s.uc.ResetPassword(ctx, req)
}
func (s *UserServiceService) PageUser(ctx context.Context, req *v1.PageUsersRequest) (*v1.PageUsersReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *UserServiceService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.MessageReply, error) {
	return s.uc.ChangePassword(ctx, req)
}
