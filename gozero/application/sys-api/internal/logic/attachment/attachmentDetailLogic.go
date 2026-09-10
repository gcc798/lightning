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

type AttachmentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDetailLogic {
	return &AttachmentDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDetailLogic) AttachmentDetail(req *types.AttachmentIdPathReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.AttachmentDetail(l.ctx, &sysservice.IdReq{Id: req.AttachmentId})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data}, nil
}
