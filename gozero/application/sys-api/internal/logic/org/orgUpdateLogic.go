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

type OrgUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgUpdateLogic {
	return &OrgUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgUpdateLogic) OrgUpdate(req *types.OrgUpdateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.OrgUpdate(l.ctx, &sysservice.OrgUpdateReq{
		Id:       req.Id,
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
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
