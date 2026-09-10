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

type AttachmentUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentUrlLogic {
	return &AttachmentUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentUrlLogic) AttachmentUrl(req *types.AttachmentUrlQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.AttachmentUrl(l.ctx, &sysservice.AttachmentUrlQueryReq{
		AttachmentId: req.AttachmentId,
		Expires:      int64(req.Expires),
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data}, nil
}
