package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuUpdateLogic) MenuUpdate(in *pb.MenuReq) (*pb.Ack, error) {
	if in.Id <= 0 {
		return nil, fmt.Errorf("菜单ID不能为空")
	}
	if in.MenuName == "" {
		return nil, fmt.Errorf("菜单名称不能为空")
	}
	if _, err := getMenuByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	if in.ParentId == in.Id {
		return nil, fmt.Errorf("不能将自己设置为父菜单")
	}
	if in.ParentId > 0 {
		parent, err := getMenuByID(l.ctx, l.svcCtx, in.ParentId)
		if err != nil {
			return nil, fmt.Errorf("父菜单不存在")
		}
		if parent.MenuType == 0 && in.MenuType == 2 {
			return nil, fmt.Errorf("目录下不能直接创建按钮")
		}
		if parent.MenuType == 1 && in.MenuType != 2 {
			return nil, fmt.Errorf("菜单下只能创建按钮")
		}
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_menu where id <> $1 and parent_id = $2 and menu_name = $3`, in.Id, in.ParentId, in.MenuName); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("同级菜单名称已存在")
	}
	_, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		update public.s_menu
		set menu_name = $2, parent_id = $3, sort = $4, path = $5, component = $6, query = $7, is_frame = $8, is_cache = $9, menu_type = $10, visible = $11, status = $12, perms = $13, icon = $14, remark = $15, update_by = nullif($16, 0), updated_time = now()
		where id = $1
	`,
		in.Id, in.MenuName, in.ParentId, in.Sort,
		sql.NullString{String: in.Path, Valid: in.Path != ""},
		sql.NullString{String: in.Component, Valid: in.Component != ""},
		sql.NullString{String: in.Query, Valid: in.Query != ""},
		in.IsFrame, in.IsCache, in.MenuType, in.Visible, in.Status,
		sql.NullString{String: in.Perms, Valid: in.Perms != ""},
		sql.NullString{String: in.Icon, Valid: in.Icon != ""},
		sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		in.UpdateBy,
	)
	if err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
