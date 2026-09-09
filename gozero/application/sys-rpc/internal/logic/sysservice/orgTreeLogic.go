package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgTreeLogic {
	return &OrgTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgTreeLogic) OrgTree(in *pb.Empty) (*pb.OrgTreeResp, error) {
	var rows []orgRow
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, `
		select id, parent_id, ancestors, org_name, org_code, org_type, leader, phone, email, status, sort, remark, created_time, updated_time
		from public.s_org
		order by sort asc, id asc
	`); err != nil {
		return nil, err
	}
	return &pb.OrgTreeResp{Records: buildOrgTree(rows, 0)}, nil
}
