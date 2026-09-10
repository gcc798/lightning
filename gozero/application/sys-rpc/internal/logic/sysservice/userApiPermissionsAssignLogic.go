package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserApiPermissionsAssignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserApiPermissionsAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserApiPermissionsAssignLogic {
	return &UserApiPermissionsAssignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserApiPermissionsAssignLogic) UserApiPermissionsAssign(in *pb.UserApiPermissionsReq) (*pb.Ack, error) {
	if _, err := getUserByID(l.ctx, l.svcCtx, in.UserId); err != nil {
		return nil, err
	}
	if err := replaceUserApiPermissions(l.ctx, l.svcCtx, in.UserId, in.PermissionIds, in.OperatorId); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
