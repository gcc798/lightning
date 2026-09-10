package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/data"
)

type RoleUsecase struct {
	res *data.Resources
}

func NewRoleUsecase(res *data.Resources) *RoleUsecase { return &RoleUsecase{res: res} }

func (uc *RoleUsecase) Create(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleItem, error) {
	return uc.res.CreateRole(ctx, req)
}

func (uc *RoleUsecase) Update(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateRole(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *RoleUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteRole(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *RoleUsecase) GetByID(ctx context.Context, id int64) (*v1.RoleItem, error) {
	item, err := uc.res.GetRole(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "角色不存在")
}

func (uc *RoleUsecase) Page(ctx context.Context, req *v1.PageRoleRequest) (*v1.PageRoleReply, error) {
	return uc.res.PageRoles(ctx, req)
}

func (uc *RoleUsecase) Assign(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	if err := uc.res.AddUserRole(ctx, req.GetUserId(), req.GetRoleId()); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *RoleUsecase) Remove(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	if err := uc.res.RemoveUserRole(ctx, req.GetUserId(), req.GetRoleId()); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *RoleUsecase) GetUserRoles(ctx context.Context, userID int64) (*v1.GetUserRolesReply, error) {
	items, err := uc.res.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserRolesReply{Items: items}, nil
}

func (uc *RoleUsecase) AddPermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	if err := uc.res.AddRolePermission(ctx, req.GetRoleKey(), req.GetResource(), req.GetAction()); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *RoleUsecase) DeletePermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	if err := uc.res.DeleteRolePermission(ctx, req.GetRoleKey(), req.GetResource(), req.GetAction()); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *RoleUsecase) GetPermissions(ctx context.Context, req *v1.GetRolePermissionsRequest) (*v1.GetRolePermissionsReply, error) {
	items, err := uc.res.GetRolePermissions(ctx, req.GetRoleKey())
	if err != nil {
		return nil, err
	}
	return &v1.GetRolePermissionsReply{Items: items}, nil
}
