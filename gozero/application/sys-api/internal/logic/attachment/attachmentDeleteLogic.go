// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDeleteLogic {
	return &AttachmentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDeleteLogic) AttachmentDelete(req *types.AttachmentIdPathReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.AttachmentDelete(l.ctx, &sysservice.IdReq{Id: req.AttachmentId}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
