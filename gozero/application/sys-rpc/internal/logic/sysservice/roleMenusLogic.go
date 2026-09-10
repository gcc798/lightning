package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenusLogic {
	return &RoleMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleMenusLogic) RoleMenus(in *pb.RoleMenusReq) (*pb.MenuIdsResp, error) {
	rows := make([]struct {
		MenuId int64 `db:"menu_id"`
	}, 0)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, `
		select menu_id
		from public.m_role_menu
		where role_id = $1
		order by menu_id asc
	`, in.RoleId); err != nil {
		return nil, err
	}
	menuIds := make([]int64, 0, len(rows))
	for _, row := range rows {
		menuIds = append(menuIds, row.MenuId)
	}
	return &pb.MenuIdsResp{MenuIds: menuIds}, nil
}
