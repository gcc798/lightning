package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRolePermissionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionDeleteLogic {
	return &RolePermissionDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RolePermissionDeleteLogic) RolePermissionDelete(in *pb.RolePermissionReq) (*pb.Ack, error) {
	if in.RoleKey == "" || in.Resource == "" || in.Action == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	if err := l.svcCtx.Redis.SRem(l.ctx, "casbin:role:"+in.RoleKey+":permissions", in.Resource+"::"+in.Action).Err(); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
