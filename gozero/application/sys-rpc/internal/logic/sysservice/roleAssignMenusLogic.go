package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAssignMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAssignMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAssignMenusLogic {
	return &RoleAssignMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAssignMenusLogic) RoleAssignMenus(in *pb.RoleMenusAssignReq) (*pb.Ack, error) {
	if in.RoleId <= 0 {
		return nil, fmt.Errorf("角色ID不能为空")
	}
	if _, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId); err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.m_role_menu where role_id = $1`, in.RoleId); err != nil {
		return nil, err
	}
	for _, menuId := range in.MenuIds {
		if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.m_role_menu (role_id, menu_id, created_time, updated_time) values ($1, $2, now(), now())`, in.RoleId, menuId); err != nil {
			return nil, err
		}
	}
	return &pb.Ack{Msg: "ok"}, nil
}
