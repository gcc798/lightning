// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package org

import (
	"context"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgPageLogic {
	return &OrgPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgPageLogic) OrgPage(req *types.OrgPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.OrgPage(l.ctx, &sysservice.OrgPageReq{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		OrgName:  req.OrgName,
		OrgCode:  req.OrgCode,
		Status:   int32(req.Status),
		ParentId: req.ParentId,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size),
	}, nil
}
