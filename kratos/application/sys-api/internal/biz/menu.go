package biz

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/data"
)

type MenuUsecase struct {
	repo *data.MenuRepo
}

func NewMenuUsecase(repo *data.MenuRepo) *MenuUsecase {
	return &MenuUsecase{repo: repo}
}

func (uc *MenuUsecase) GetUserMenuTree(ctx context.Context) (*v1.MenuTreeReply, error) {
	return uc.repo.GetUserMenuTree(ctx)
}

func (uc *MenuUsecase) GetMenuTree(ctx context.Context) (*v1.MenuTreeReply, error) {
	return uc.repo.GetMenuTree(ctx)
}

func (uc *MenuUsecase) GetMenuList(ctx context.Context) (*v1.MenuListReply, error) {
	return uc.repo.GetMenuList(ctx)
}

func (uc *MenuUsecase) GetMenuByID(ctx context.Context, id int64) (*v1.MenuItem, error) {
	return uc.repo.GetMenuByID(ctx, id)
}

func (uc *MenuUsecase) CreateMenu(ctx context.Context, item *v1.MenuItem) (*v1.MessageReply, error) {
	return uc.repo.CreateMenu(ctx, item)
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) (*v1.MessageReply, error) {
	return uc.repo.UpdateMenu(ctx, req)
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.DeleteMenu(ctx, id)
}
