package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiPermissionsLogic {
	return &RoleApiPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiPermissionsLogic) RoleApiPermissions(in *pb.IdReq) (*pb.ApiPermissionIdsResp, error) {
	ids, err := apiPermissionIDsByOwner(l.ctx, l.svcCtx, "m_role_api_permission", "role_id", in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ApiPermissionIdsResp{PermissionIds: ids}, nil
}
