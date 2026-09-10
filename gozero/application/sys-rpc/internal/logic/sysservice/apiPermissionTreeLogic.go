package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiPermissionTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiPermissionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionTreeLogic {
	return &ApiPermissionTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiPermissionTreeLogic) ApiPermissionTree(in *pb.Empty) (*pb.ApiPermissionListResp, error) {
	rows, err := listApiPermissions(l.ctx, l.svcCtx)
	if err != nil {
		return nil, err
	}
	return &pb.ApiPermissionListResp{Records: buildApiPermissionTree(rows, 0)}, nil
}
