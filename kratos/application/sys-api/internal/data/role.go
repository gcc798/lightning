package data

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type RoleRepo struct {
	client v1.RoleServiceClient
}

func NewRoleRepo(clients *RPCClientSet) *RoleRepo {
	return &RoleRepo{client: clients.Role}
}
func (r *RoleRepo) Create(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleItem, error) {
	return r.client.CreateRole(ctx, req)
}
func (r *RoleRepo) Update(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.MessageReply, error) {
	return r.client.UpdateRole(ctx, req)
}
func (r *RoleRepo) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.DeleteRole(ctx, &v1.RoleIdRequest{RoleId: id})
}
func (r *RoleRepo) GetByID(ctx context.Context, id int64) (*v1.RoleItem, error) {
	return r.client.GetRole(ctx, &v1.RoleIdRequest{RoleId: id})
}
func (r *RoleRepo) Page(ctx context.Context, req *v1.PageRoleRequest) (*v1.PageRoleReply, error) {
	return r.client.PageRole(ctx, req)
}
func (r *RoleRepo) Assign(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	return r.client.AssignRoleToUser(ctx, req)
}
func (r *RoleRepo) Remove(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	return r.client.RemoveRoleFromUser(ctx, req)
}
func (r *RoleRepo) GetUserRoles(ctx context.Context, userID int64) (*v1.GetUserRolesReply, error) {
	return r.client.GetUserRoles(ctx, &v1.GetUserRolesRequest{UserId: userID})
}
func (r *RoleRepo) AddPermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	return r.client.AddRolePermission(ctx, req)
}
func (r *RoleRepo) DeletePermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	return r.client.DeleteRolePermission(ctx, req)
}
func (r *RoleRepo) GetPermissions(ctx context.Context, roleKey string) (*v1.GetRolePermissionsReply, error) {
	return r.client.GetRolePermissions(ctx, &v1.GetRolePermissionsRequest{RoleKey: roleKey})
}
