package biz

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

type RoleUsecase struct{ repo *data.RoleRepo }

func NewRoleUsecase(repo *data.RoleRepo) *RoleUsecase { return &RoleUsecase{repo: repo} }
func (uc *RoleUsecase) Create(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleItem, error) {
	return uc.repo.Create(ctx, req)
}
func (uc *RoleUsecase) Update(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.MessageReply, error) {
	return uc.repo.Update(ctx, req)
}
func (uc *RoleUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
func (uc *RoleUsecase) GetByID(ctx context.Context, id int64) (*v1.RoleItem, error) {
	return uc.repo.GetByID(ctx, id)
}
func (uc *RoleUsecase) Page(ctx context.Context, req *v1.PageRoleRequest) (*v1.PageRoleReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *RoleUsecase) Assign(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	return uc.repo.Assign(ctx, req)
}
func (uc *RoleUsecase) Remove(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	return uc.repo.Remove(ctx, req)
}
func (uc *RoleUsecase) GetUserRoles(ctx context.Context, userID int64) (*v1.GetUserRolesReply, error) {
	return uc.repo.GetUserRoles(ctx, userID)
}
func (uc *RoleUsecase) AddPermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	return uc.repo.AddPermission(ctx, req)
}
func (uc *RoleUsecase) DeletePermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	return uc.repo.DeletePermission(ctx, req)
}
func (uc *RoleUsecase) GetPermissions(ctx context.Context, roleKey string) (*v1.GetRolePermissionsReply, error) {
	return uc.repo.GetPermissions(ctx, roleKey)
}
