package service

import (
	"context"

	pb "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type MenuServiceService struct {
	pb.UnimplementedMenuServiceServer
}

func NewMenuServiceService() *MenuServiceService {
	return &MenuServiceService{}
}

func (s *MenuServiceService) GetUserMenuTree(ctx context.Context, req *pb.GetUserMenuTreeRequest) (*pb.MenuTreeReply, error) {
	return &pb.MenuTreeReply{}, nil
}
func (s *MenuServiceService) GetMenuTree(ctx context.Context, req *pb.GetMenuTreeRequest) (*pb.MenuTreeReply, error) {
	return &pb.MenuTreeReply{}, nil
}
func (s *MenuServiceService) GetMenuList(ctx context.Context, req *pb.GetMenuListRequest) (*pb.MenuListReply, error) {
	return &pb.MenuListReply{}, nil
}
func (s *MenuServiceService) GetMenuById(ctx context.Context, req *pb.MenuIdRequest) (*pb.MenuItem, error) {
	return &pb.MenuItem{}, nil
}
func (s *MenuServiceService) CreateMenu(ctx context.Context, req *pb.MenuItem) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *MenuServiceService) UpdateMenu(ctx context.Context, req *pb.UpdateMenuRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *MenuServiceService) DeleteMenu(ctx context.Context, req *pb.MenuIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
