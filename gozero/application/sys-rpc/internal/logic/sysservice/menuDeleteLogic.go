package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuDeleteLogic {
	return &MenuDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuDeleteLogic) MenuDelete(in *pb.IdReq) (*pb.Ack, error) {
	if in.Id <= 0 {
		return nil, fmt.Errorf("菜单ID不能为空")
	}
	if _, err := getMenuByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_menu where parent_id = $1`, in.Id); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("存在子菜单，无法删除")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_menu where id = $1`, in.Id); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
