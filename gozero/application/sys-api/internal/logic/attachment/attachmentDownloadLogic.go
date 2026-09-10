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

type AttachmentDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDownloadLogic {
	return &AttachmentDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDownloadLogic) AttachmentDownload(req *types.AttachmentIdPathReq) (*sysservice.AttachmentDownloadResp, error) {
	return l.svcCtx.SysRpcClient.AttachmentDownload(l.ctx, &sysservice.AttachmentDownloadReq{AttachmentId: req.AttachmentId})
}
