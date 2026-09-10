package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleCreateLogic) RoleCreate(in *pb.RoleCreateReq) (*pb.Ack, error) {
	if in.RoleKey == "" || in.RoleName == "" {
		return nil, errors.New("角色标识和角色名称不能为空")
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_role where role_key = $1`, in.RoleKey); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("角色标识已存在")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_role (role_key, role_name, sort, status, data_scope, is_system, remark, created_time, updated_time) values ($1, $2, $3, $4, $5, false, $6, now(), now())`,
		in.RoleKey, in.RoleName, in.Sort, in.Status, in.DataScope, sql.NullString{String: in.Remark, Valid: in.Remark != ""}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
