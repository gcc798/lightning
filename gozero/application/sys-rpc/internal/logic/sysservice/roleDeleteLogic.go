package sysservicelogic

import (
	"context"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDeleteLogic {
	return &RoleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleDeleteLogic) RoleDelete(in *pb.IdReq) (*pb.Ack, error) {
	row, err := getRoleByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	if row.IsSystem {
		return nil, errors.New("系统角色不可删除")
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.m_user_role where role_id = $1`, in.Id); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该角色已被用户使用，无法删除")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.m_role_menu where role_id = $1`, in.Id); err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_role where id = $1`, in.Id); err != nil {
		return nil, err
	}
	if err := l.svcCtx.Redis.Del(l.ctx, "casbin:role:"+row.RoleKey+":permissions").Err(); err != nil {
		l.Errorf("删除角色权限缓存失败 roleKey=%s: %v", row.RoleKey, err)
	}
	return &pb.Ack{Msg: "ok"}, nil
}
