package data

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type MenuRepo struct {
	client v1.MenuServiceClient
}

func NewMenuRepo(clients *RPCClientSet) *MenuRepo {
	return &MenuRepo{client: clients.Menu}
}

func (r *MenuRepo) GetUserMenuTree(ctx context.Context) (*v1.MenuTreeReply, error) {
	return r.client.GetUserMenuTree(ctx, &v1.GetUserMenuTreeRequest{})
}

func (r *MenuRepo) GetMenuTree(ctx context.Context) (*v1.MenuTreeReply, error) {
	return r.client.GetMenuTree(ctx, &v1.GetMenuTreeRequest{})
}

func (r *MenuRepo) GetMenuList(ctx context.Context) (*v1.MenuListReply, error) {
	return r.client.GetMenuList(ctx, &v1.GetMenuListRequest{})
}

func (r *MenuRepo) GetMenuByID(ctx context.Context, id int64) (*v1.MenuItem, error) {
	return r.client.GetMenuById(ctx, &v1.MenuIdRequest{Id: id})
}

func (r *MenuRepo) CreateMenu(ctx context.Context, item *v1.MenuItem) (*v1.MessageReply, error) {
	return r.client.CreateMenu(ctx, item)
}

func (r *MenuRepo) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) (*v1.MessageReply, error) {
	return r.client.UpdateMenu(ctx, req)
}

func (r *MenuRepo) DeleteMenu(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.DeleteMenu(ctx, &v1.MenuIdRequest{Id: id})
}
