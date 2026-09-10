package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiPermissionsAssignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiPermissionsAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiPermissionsAssignLogic {
	return &RoleApiPermissionsAssignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiPermissionsAssignLogic) RoleApiPermissionsAssign(in *pb.RoleApiPermissionsReq) (*pb.Ack, error) {
	role, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId)
	if err != nil {
		return nil, err
	}
	if err := replaceRoleApiPermissions(l.ctx, l.svcCtx, in.RoleId, role.RoleKey, in.PermissionIds, in.OperatorId); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
