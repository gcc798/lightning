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

type OrgBatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgBatchDeleteLogic {
	return &OrgBatchDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgBatchDeleteLogic) OrgBatchDelete(req *types.BatchIdsReq) (resp *types.CommonResp, err error) {
	if len(req.Ids) == 0 {
		return &types.CommonResp{Code: 400, Msg: "请至少选择一条数据"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.OrgBatchDelete(l.ctx, &sysservice.BatchIdsReq{Ids: req.Ids}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功"}, nil
}
