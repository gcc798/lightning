package service

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/biz"
)

type RoleServiceService struct {
	v1.UnimplementedRoleServiceServer
	uc *biz.RoleUsecase
}

func NewRoleServiceService(uc *biz.RoleUsecase) *RoleServiceService {
	return &RoleServiceService{uc: uc}
}
func (s *RoleServiceService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleItem, error) {
	return s.uc.Create(ctx, req)
}
func (s *RoleServiceService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *RoleServiceService) DeleteRole(ctx context.Context, req *v1.RoleIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetRoleId())
}
func (s *RoleServiceService) GetRole(ctx context.Context, req *v1.RoleIdRequest) (*v1.RoleItem, error) {
	return s.uc.GetByID(ctx, req.GetRoleId())
}
func (s *RoleServiceService) PageRole(ctx context.Context, req *v1.PageRoleRequest) (*v1.PageRoleReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *RoleServiceService) AssignRoleToUser(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	return s.uc.Assign(ctx, req)
}
func (s *RoleServiceService) RemoveRoleFromUser(ctx context.Context, req *v1.AssignRoleToUserRequest) (*v1.MessageReply, error) {
	return s.uc.Remove(ctx, req)
}
func (s *RoleServiceService) GetUserRoles(ctx context.Context, req *v1.GetUserRolesRequest) (*v1.GetUserRolesReply, error) {
	return s.uc.GetUserRoles(ctx, req.GetUserId())
}
func (s *RoleServiceService) AddRolePermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	return s.uc.AddPermission(ctx, req)
}
func (s *RoleServiceService) DeleteRolePermission(ctx context.Context, req *v1.RolePermissionRequest) (*v1.MessageReply, error) {
	return s.uc.DeletePermission(ctx, req)
}
func (s *RoleServiceService) GetRolePermissions(ctx context.Context, req *v1.GetRolePermissionsRequest) (*v1.GetRolePermissionsReply, error) {
	return s.uc.GetPermissions(ctx, req)
}
