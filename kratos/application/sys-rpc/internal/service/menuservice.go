package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/biz"
)

type MenuServiceService struct {
	v1.UnimplementedMenuServiceServer

	uc *biz.MenuUsecase
}

func NewMenuServiceService(uc *biz.MenuUsecase) *MenuServiceService {
	return &MenuServiceService{uc: uc}
}

func (s *MenuServiceService) GetUserMenuTree(ctx context.Context, req *v1.GetUserMenuTreeRequest) (*v1.MenuTreeReply, error) {
	items, err := s.uc.GetUserMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MenuTreeReply{Items: items}, nil
}

func (s *MenuServiceService) GetMenuTree(ctx context.Context, req *v1.GetMenuTreeRequest) (*v1.MenuTreeReply, error) {
	items, err := s.uc.GetMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MenuTreeReply{Items: items}, nil
}

func (s *MenuServiceService) GetMenuList(ctx context.Context, req *v1.GetMenuListRequest) (*v1.MenuListReply, error) {
	items, err := s.uc.GetMenuList(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MenuListReply{Items: items}, nil
}

func (s *MenuServiceService) GetMenuById(ctx context.Context, req *v1.MenuIdRequest) (*v1.MenuItem, error) {
	return s.uc.GetMenuByID(ctx, req.GetId())
}

func (s *MenuServiceService) CreateMenu(ctx context.Context, req *v1.MenuItem) (*v1.MessageReply, error) {
	return s.uc.CreateMenu(ctx, req)
}

func (s *MenuServiceService) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) (*v1.MessageReply, error) {
	return s.uc.UpdateMenu(ctx, req)
}

func (s *MenuServiceService) DeleteMenu(ctx context.Context, req *v1.MenuIdRequest) (*v1.MessageReply, error) {
	return s.uc.DeleteMenu(ctx, req.GetId())
}
