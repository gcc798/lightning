package sysservicelogic

import (
	"context"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleRemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleRemoveLogic {
	return &RoleRemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleRemoveLogic) RoleRemove(in *pb.RemoveRoleReq) (*pb.Ack, error) {
	if _, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId); err != nil {
		return nil, err
	}
	result, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.m_user_role where user_id = $1 and role_id = $2`, in.UserId, in.RoleId)
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		return nil, errors.New("用户角色关系不存在")
	}
	return &pb.Ack{Msg: "ok"}, nil
}
