package sysservicelogic

import (
	"context"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgDeleteLogic {
	return &OrgDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgDeleteLogic) OrgDelete(in *pb.IdReq) (*pb.Ack, error) {
	if _, err := getOrgByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	var childCount int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &childCount, `select count(1) from public.s_org where parent_id = $1`, in.Id); err != nil {
		return nil, err
	}
	if childCount > 0 {
		return nil, errors.New("存在子组织，无法删除")
	}
	var userCount int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &userCount, `select count(1) from public.s_user where org_id = $1`, in.Id); err != nil {
		return nil, err
	}
	if userCount > 0 {
		return nil, errors.New("组织下存在用户，无法删除")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_org where id = $1`, in.Id); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
