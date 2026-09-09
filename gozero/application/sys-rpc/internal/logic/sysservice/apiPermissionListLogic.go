package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiPermissionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiPermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionListLogic {
	return &ApiPermissionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiPermissionListLogic) ApiPermissionList(in *pb.Empty) (*pb.ApiPermissionListResp, error) {
	rows, err := listApiPermissions(l.ctx, l.svcCtx)
	if err != nil {
		return nil, err
	}
	return &pb.ApiPermissionListResp{Records: toApiPermissionList(rows)}, nil
}
