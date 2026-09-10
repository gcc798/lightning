package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserApiPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserApiPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserApiPermissionsLogic {
	return &UserApiPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserApiPermissionsLogic) UserApiPermissions(in *pb.IdReq) (*pb.ApiPermissionIdsResp, error) {
	ids, err := apiPermissionIDsByOwner(l.ctx, l.svcCtx, "m_user_api_permission", "user_id", in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ApiPermissionIdsResp{PermissionIds: ids}, nil
}
