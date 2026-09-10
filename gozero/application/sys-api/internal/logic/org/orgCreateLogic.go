// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package org

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgCreateLogic {
	return &OrgCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgCreateLogic) OrgCreate(req *types.OrgCreateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.OrgCreate(l.ctx, &sysservice.OrgCreateReq{
		ParentId: req.ParentId,
		OrgName:  req.OrgName,
		OrgCode:  req.OrgCode,
		OrgType:  req.OrgType,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   int32(req.Status),
		Sort:     req.Sort,
		Remark:   req.Remark,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: nil}, nil
}
