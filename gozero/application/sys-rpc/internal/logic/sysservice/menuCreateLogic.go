package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuCreateLogic {
	return &MenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuCreateLogic) MenuCreate(in *pb.MenuReq) (*pb.Ack, error) {
	if in.MenuName == "" {
		return nil, fmt.Errorf("菜单名称不能为空")
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
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_menu where parent_id = $1 and menu_name = $2`, in.ParentId, in.MenuName); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("同级菜单名称已存在")
	}
	_, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		insert into public.s_menu (menu_name, parent_id, sort, path, component, query, is_frame, is_cache, menu_type, visible, status, perms, icon, remark, create_by, update_by, created_time, updated_time)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, nullif($15, 0), nullif($16, 0), now(), now())
	`,
		in.MenuName, in.ParentId, in.Sort,
		sql.NullString{String: in.Path, Valid: in.Path != ""},
		sql.NullString{String: in.Component, Valid: in.Component != ""},
		sql.NullString{String: in.Query, Valid: in.Query != ""},
		in.IsFrame, in.IsCache, in.MenuType, in.Visible, in.Status,
		sql.NullString{String: in.Perms, Valid: in.Perms != ""},
		sql.NullString{String: in.Icon, Valid: in.Icon != ""},
		sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		in.CreateBy, in.UpdateBy,
	)
	if err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
