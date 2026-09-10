package data

import (
	"context"
	"errors"
	"sort"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
)

func menuEntityToItem(item *entpkg.Menu) *v1.MenuItem {
	if item == nil {
		return nil
	}
	return &v1.MenuItem{
		Id:          item.ID,
		MenuName:    item.MenuName,
		ParentId:    item.ParentID,
		Sort:        item.Sort,
		Path:        item.Path,
		Component:   item.Component,
		Query:       item.Query,
		IsFrame:     item.IsFrame,
		IsCache:     item.IsCache,
		MenuType:    item.MenuType,
		Visible:     item.Visible,
		Status:      item.Status,
		Perms:       item.Perms,
		Icon:        item.Icon,
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func (r *Resources) activeMenus(ctx context.Context) ([]*entpkg.Menu, error) {
	items, err := r.Ent.Menu.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.Menu, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func menuCanHaveChild(parentType int32, childType int32) bool {
	if parentType == 0 && (childType == 0 || childType == 1) {
		return true
	}
	if parentType == 1 && childType == 2 {
		return true
	}
	return false
}

func ensureMenuNameUnique(items []*entpkg.Menu, menuID int64, parentID int64, menuName string) error {
	for _, item := range items {
		if menuID > 0 && item.ID == menuID {
			continue
		}
		if item.ParentID == parentID && item.MenuName == menuName {
			return errors.New("menu name already exists in current level")
		}
	}
	return nil
}

func (r *Resources) CreateMenu(ctx context.Context, req *v1.MenuItem) error {
	if req.GetId() != 0 {
		return errors.New("menu id must be empty when creating")
	}
	items, err := r.activeMenus(ctx)
	if err != nil {
		return err
	}
	if err := ensureMenuNameUnique(items, 0, req.GetParentId(), req.GetMenuName()); err != nil {
		return err
	}
	if req.GetParentId() != 0 {
		var parent *entpkg.Menu
		for _, item := range items {
			if item.ID == req.GetParentId() {
				parent = item
				break
			}
		}
		if parent == nil {
			return errors.New("parent menu not found")
		}
		if !menuCanHaveChild(parent.MenuType, req.GetMenuType()) {
			return errors.New("invalid menu hierarchy")
		}
	}
	now := time.Now()
	createBy := operatorID(ctx, req.GetCreateBy())
	updateBy := operatorID(ctx, req.GetUpdateBy())
	_, err = r.Ent.Menu.Create().
		SetID(nextID()).
		SetMenuName(req.GetMenuName()).
		SetParentID(req.GetParentId()).
		SetSort(req.GetSort()).
		SetPath(req.GetPath()).
		SetComponent(req.GetComponent()).
		SetQuery(req.GetQuery()).
		SetIsFrame(req.GetIsFrame()).
		SetIsCache(req.GetIsCache()).
		SetMenuType(req.GetMenuType()).
		SetVisible(req.GetVisible()).
		SetStatus(req.GetStatus()).
		SetPerms(req.GetPerms()).
		SetIcon(req.GetIcon()).
		SetRemark(req.GetRemark()).
		SetCreateBy(createBy).
		SetUpdateBy(updateBy).
		SetCreatedTime(now).
		SetUpdatedTime(now).
		Save(ctx)
	return err
}

func (r *Resources) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) error {
	item, err := r.Ent.Menu.Get(ctx, req.GetId())
	if err != nil {
		return err
	}
	if item.DeletedAt != nil {
		return errors.New("menu not found")
	}
	if req.GetParentId() == req.GetId() {
		return errors.New("menu can not be its own parent")
	}
	items, err := r.activeMenus(ctx)
	if err != nil {
		return err
	}
	if err := ensureMenuNameUnique(items, req.GetId(), req.GetParentId(), req.GetMenuName()); err != nil {
		return err
	}
	if req.GetParentId() != 0 {
		var parent *entpkg.Menu
		for _, current := range items {
			if current.ID == req.GetParentId() {
				parent = current
				break
			}
		}
		if parent == nil {
			return errors.New("parent menu not found")
		}
		if !menuCanHaveChild(parent.MenuType, req.GetMenuType()) {
			return errors.New("invalid menu hierarchy")
		}
	}
	_, err = r.Ent.Menu.UpdateOneID(req.GetId()).
		SetMenuName(req.GetMenuName()).
		SetParentID(req.GetParentId()).
		SetSort(req.GetSort()).
		SetPath(req.GetPath()).
		SetComponent(req.GetComponent()).
		SetQuery(req.GetQuery()).
		SetIsFrame(req.GetIsFrame()).
		SetIsCache(req.GetIsCache()).
		SetMenuType(req.GetMenuType()).
		SetVisible(req.GetVisible()).
		SetStatus(req.GetStatus()).
		SetPerms(req.GetPerms()).
		SetIcon(req.GetIcon()).
		SetRemark(req.GetRemark()).
		SetUpdateBy(operatorID(ctx, req.GetUpdateBy())).
		SetUpdatedTime(time.Now()).
		Save(ctx)
	return err
}

func (r *Resources) GetMenu(ctx context.Context, id int64) (*v1.MenuItem, error) {
	item, err := r.Ent.Menu.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return menuEntityToItem(item), nil
}

func (r *Resources) ListMenus(ctx context.Context) ([]*v1.MenuItem, error) {
	items, err := r.activeMenus(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.MenuItem, 0, len(items))
	for _, item := range items {
		out = append(out, menuEntityToItem(item))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].GetSort() == out[j].GetSort() {
			return out[i].GetId() < out[j].GetId()
		}
		return out[i].GetSort() < out[j].GetSort()
	})
	return out, nil
}

func (r *Resources) DeleteMenu(ctx context.Context, id int64) error {
	operator := currentOperatorID(ctx)
	return r.withTx(ctx, func(tx *entpkg.Tx) error {
		menu, err := tx.Menu.Get(ctx, id)
		if err != nil {
			return err
		}
		if menu.DeletedAt != nil {
			return errors.New("menu not found")
		}
		menus, err := tx.Menu.Query().All(ctx)
		if err != nil {
			return err
		}
		for _, item := range menus {
			if item.DeletedAt == nil && item.ParentID == id {
				return errors.New("menu has children")
			}
		}
		now := time.Now()
		if _, err := tx.Menu.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
		roleMenus, err := tx.RoleMenu.Query().All(ctx)
		if err != nil {
			return err
		}
		for _, item := range roleMenus {
			if item.DeletedAt != nil {
				continue
			}
			if item.MenuID == id {
				if _, err := tx.RoleMenu.UpdateOneID(item.ID).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
					return err
				}
			}
		}
		return nil
	})
}

func (r *Resources) UserMenus(ctx context.Context, userID int64) ([]*v1.MenuItem, error) {
	if userID <= 0 {
		return []*v1.MenuItem{}, nil
	}
	userRoles, err := r.activeUserRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleIDs := make(map[int64]struct{})
	for _, item := range userRoles {
		if item.UserID == userID {
			roleIDs[item.RoleID] = struct{}{}
		}
	}
	if len(roleIDs) == 0 {
		return []*v1.MenuItem{}, nil
	}
	menus, err := r.activeMenus(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := r.activeRoles(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range roles {
		if _, ok := roleIDs[item.ID]; ok && item.RoleKey == "super_admin" {
			out := make([]*v1.MenuItem, 0, len(menus))
			for _, menu := range menus {
				if menu.Status == 0 {
					out = append(out, menuEntityToItem(menu))
				}
			}
			return out, nil
		}
	}
	roleMenus, err := r.activeRoleMenus(ctx)
	if err != nil {
		return nil, err
	}
	menuIDs := make(map[int64]struct{})
	for _, item := range roleMenus {
		if _, ok := roleIDs[item.RoleID]; ok {
			menuIDs[item.MenuID] = struct{}{}
		}
	}
	out := make([]*v1.MenuItem, 0)
	for _, item := range menus {
		if _, ok := menuIDs[item.ID]; ok && item.Status == 0 {
			out = append(out, menuEntityToItem(item))
		}
	}
	return out, nil
}
