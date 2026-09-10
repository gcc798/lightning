package data

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
)

func roleEntityToItem(item *entpkg.Role) *v1.RoleItem {
	if item == nil {
		return nil
	}
	return &v1.RoleItem{
		RoleId:      item.ID,
		RoleKey:     item.RoleKey,
		RoleName:    item.RoleName,
		Sort:        item.Sort,
		Status:      item.Status,
		DataScope:   item.DataScope,
		IsSystem:    item.IsSystem,
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func (r *Resources) activeRoles(ctx context.Context) ([]*entpkg.Role, error) {
	items, err := r.Ent.Role.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.Role, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *Resources) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleItem, error) {
	items, err := r.activeRoles(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.RoleKey == req.GetRoleKey() {
			return nil, errors.New("role key already exists")
		}
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	item, err := r.Ent.Role.Create().
		SetID(nextID()).
		SetRoleKey(req.GetRoleKey()).
		SetRoleName(req.GetRoleName()).
		SetSort(req.GetSort()).
		SetStatus(req.GetStatus()).
		SetDataScope(req.GetDataScope()).
		SetRemark(req.GetRemark()).
		SetCreateBy(operator).
		SetUpdateBy(operator).
		SetCreatedTime(now).
		SetUpdatedTime(now).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return roleEntityToItem(item), nil
}

func (r *Resources) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) error {
	item, err := r.Ent.Role.Get(ctx, req.GetRoleId())
	if err != nil {
		return err
	}
	if item.DeletedAt != nil {
		return errors.New("role not found")
	}
	operator := currentOperatorID(ctx)
	_, err = r.Ent.Role.UpdateOneID(req.GetRoleId()).
		SetRoleName(req.GetRoleName()).
		SetSort(req.GetSort()).
		SetStatus(req.GetStatus()).
		SetDataScope(req.GetDataScope()).
		SetRemark(req.GetRemark()).
		SetUpdateBy(operator).
		SetUpdatedTime(time.Now()).
		Save(ctx)
	return err
}

func (r *Resources) DeleteRole(ctx context.Context, id int64) error {
	now := time.Now()
	operator := currentOperatorID(ctx)
	return r.withTx(ctx, func(tx *entpkg.Tx) error {
		role, err := tx.Role.Get(ctx, id)
		if err != nil {
			return err
		}
		if role.DeletedAt != nil {
			return errors.New("role not found")
		}
		if role.IsSystem {
			return errors.New("system role can not be deleted")
		}
		userRoles, err := tx.UserRole.Query().All(ctx)
		if err != nil {
			return err
		}
		for _, item := range userRoles {
			if item.DeletedAt == nil && item.RoleID == id {
				return errors.New("role is in use")
			}
		}
		if _, err := tx.Role.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
		roleMenus, err := tx.RoleMenu.Query().All(ctx)
		if err != nil {
			return err
		}
		for _, item := range roleMenus {
			if item.DeletedAt == nil && item.RoleID == id {
				if _, err := tx.RoleMenu.UpdateOneID(item.ID).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
					return err
				}
			}
		}
		return nil
	})
}

func (r *Resources) GetRole(ctx context.Context, id int64) (*v1.RoleItem, error) {
	item, err := r.Ent.Role.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return roleEntityToItem(item), nil
}

func (r *Resources) PageRoles(ctx context.Context, req *v1.PageRoleRequest) (*v1.PageRoleReply, error) {
	items, err := r.activeRoles(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.RoleItem, 0, len(items))
	for _, item := range items {
		if req.GetRoleName() != "" && !strings.Contains(item.RoleName, req.GetRoleName()) {
			continue
		}
		if req.Status != nil && req.GetStatus() >= 0 && item.Status != req.GetStatus() {
			continue
		}
		filtered = append(filtered, roleEntityToItem(item))
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].GetSort() == filtered[j].GetSort() {
			return filtered[i].GetRoleId() < filtered[j].GetRoleId()
		}
		return filtered[i].GetSort() < filtered[j].GetSort()
	})
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageRoleReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}
