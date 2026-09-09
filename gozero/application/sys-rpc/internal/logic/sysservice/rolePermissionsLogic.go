package sysservicelogic

import (
	"context"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionsLogic {
	return &RolePermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RolePermissionsLogic) RolePermissions(in *pb.RolePermissionsQueryReq) (*pb.RolePermissionsResp, error) {
	rows, err := l.svcCtx.Redis.SMembers(l.ctx, "casbin:role:"+in.RoleKey+":permissions").Result()
	if err != nil {
		return nil, err
	}
	resp := make([]*pb.RolePermission, 0, len(rows))
	for _, row := range rows {
		parts := strings.SplitN(row, "::", 2)
		perm := &pb.RolePermission{RoleKey: in.RoleKey}
		if len(parts) > 0 {
			perm.Resource = parts[0]
		}
		if len(parts) > 1 {
			perm.Action = parts[1]
		}
		resp = append(resp, perm)
	}
	return &pb.RolePermissionsResp{Records: resp}, nil
}
