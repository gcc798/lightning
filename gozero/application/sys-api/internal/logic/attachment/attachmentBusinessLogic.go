// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentBusinessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentBusinessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentBusinessLogic {
	return &AttachmentBusinessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentBusinessLogic) AttachmentBusiness(req *types.AttachmentBusinessQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.AttachmentBusiness(l.ctx, &sysservice.AttachmentBusinessQueryReq{BusinessType: req.BusinessType, BusinessId: req.BusinessId})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.Records}, nil
}
