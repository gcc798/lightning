package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUsersLogic {
	return &RoleUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleUsersLogic) RoleUsers(in *pb.RoleUsersReq) (*pb.UserListResp, error) {
	if _, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId); err != nil {
		return nil, err
	}
	rows, err := listUsersByRole(l.ctx, l.svcCtx, in.RoleId)
	if err != nil {
		return nil, err
	}
	return &pb.UserListResp{Records: toUserListPB(rows)}, nil
}
