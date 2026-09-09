// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package org

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgTreeLogic {
	return &OrgTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgTreeLogic) OrgTree() (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.OrgTree(l.ctx, &sysservice.Empty{})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: data.Records,
	}, nil
}
