package biz

import (
	"context"
	"sort"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type MenuUsecase struct {
	res *data.Resources
}

func NewMenuUsecase(res *data.Resources) *MenuUsecase {
	return &MenuUsecase{res: res}
}

func (uc *MenuUsecase) GetUserMenuTree(ctx context.Context) ([]*v1.MenuItem, error) {
	items, err := uc.res.UserMenus(ctx, currentUserID(ctx))
	if err != nil {
		return nil, err
	}
	return buildMenuTree(items), nil
}

func (uc *MenuUsecase) GetMenuTree(ctx context.Context) ([]*v1.MenuItem, error) {
	items, err := uc.res.ListMenus(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(items), nil
}

func (uc *MenuUsecase) GetMenuList(ctx context.Context) ([]*v1.MenuItem, error) {
	return uc.res.ListMenus(ctx)
}

func (uc *MenuUsecase) GetMenuByID(ctx context.Context, id int64) (*v1.MenuItem, error) {
	item, err := uc.res.GetMenu(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "菜单不存在")
}

func (uc *MenuUsecase) CreateMenu(ctx context.Context, req *v1.MenuItem) (*v1.MessageReply, error) {
	if err := uc.res.CreateMenu(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateMenu(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteMenu(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func buildMenuTree(items []*v1.MenuItem) []*v1.MenuItem {
	byParent := make(map[int64][]*v1.MenuItem)
	for _, item := range items {
		item.Children = nil
		byParent[item.GetParentId()] = append(byParent[item.GetParentId()], item)
	}

	var build func(parentID int64) []*v1.MenuItem
	build = func(parentID int64) []*v1.MenuItem {
		nodes := byParent[parentID]
		sort.Slice(nodes, func(i, j int) bool {
			if nodes[i].GetSort() == nodes[j].GetSort() {
				return nodes[i].GetId() < nodes[j].GetId()
			}
			return nodes[i].GetSort() < nodes[j].GetSort()
		})
		for _, node := range nodes {
			node.Children = build(node.GetId())
		}
		return nodes
	}
	return build(0)
}
