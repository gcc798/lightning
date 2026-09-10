package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgDetailLogic {
	return &OrgDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgDetailLogic) OrgDetail(in *pb.IdReq) (*pb.Org, error) {
	row, err := getOrgByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return toOrgPB(*row), nil
}
