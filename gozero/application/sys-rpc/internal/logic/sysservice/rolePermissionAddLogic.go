package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRolePermissionAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionAddLogic {
	return &RolePermissionAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RolePermissionAddLogic) RolePermissionAdd(in *pb.RolePermissionReq) (*pb.Ack, error) {
	if in.RoleKey == "" || in.Resource == "" || in.Action == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	if err := l.svcCtx.Redis.SAdd(l.ctx, "casbin:role:"+in.RoleKey+":permissions", in.Resource+"::"+in.Action).Err(); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
