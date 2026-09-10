package service

import (
	"context"

	pb "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceService() *UserServiceService {
	return &UserServiceService{}
}

func (s *UserServiceService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *UserServiceService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *UserServiceService) DeleteUser(ctx context.Context, req *pb.UserIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *UserServiceService) BatchDeleteUser(ctx context.Context, req *pb.BatchDeleteUsersRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *UserServiceService) GetUserById(ctx context.Context, req *pb.UserIdRequest) (*pb.UserItem, error) {
	return &pb.UserItem{}, nil
}
func (s *UserServiceService) ImportUsers(ctx context.Context, req *pb.BatchImportUsersRequest) (*pb.ImportUsersReply, error) {
	return &pb.ImportUsersReply{}, nil
}
func (s *UserServiceService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *UserServiceService) PageUser(ctx context.Context, req *pb.PageUsersRequest) (*pb.PageUsersReply, error) {
	return &pb.PageUsersReply{}, nil
}
func (s *UserServiceService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
