package service

import (
	"context"

	pb "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type RoleServiceService struct {
	pb.UnimplementedRoleServiceServer
}

func NewRoleServiceService() *RoleServiceService {
	return &RoleServiceService{}
}

func (s *RoleServiceService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.RoleItem, error) {
	return &pb.RoleItem{}, nil
}
func (s *RoleServiceService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *RoleServiceService) DeleteRole(ctx context.Context, req *pb.RoleIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *RoleServiceService) GetRole(ctx context.Context, req *pb.RoleIdRequest) (*pb.RoleItem, error) {
	return &pb.RoleItem{}, nil
}
func (s *RoleServiceService) PageRole(ctx context.Context, req *pb.PageRoleRequest) (*pb.PageRoleReply, error) {
	return &pb.PageRoleReply{}, nil
}
func (s *RoleServiceService) AssignRoleToUser(ctx context.Context, req *pb.AssignRoleToUserRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *RoleServiceService) RemoveRoleFromUser(ctx context.Context, req *pb.AssignRoleToUserRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *RoleServiceService) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesReply, error) {
	return &pb.GetUserRolesReply{}, nil
}
func (s *RoleServiceService) AddRolePermission(ctx context.Context, req *pb.RolePermissionRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *RoleServiceService) DeleteRolePermission(ctx context.Context, req *pb.RolePermissionRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *RoleServiceService) GetRolePermissions(ctx context.Context, req *pb.GetRolePermissionsRequest) (*pb.GetRolePermissionsReply, error) {
	return &pb.GetRolePermissionsReply{}, nil
}
