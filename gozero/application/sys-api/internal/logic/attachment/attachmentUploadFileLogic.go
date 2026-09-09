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

type AttachmentUploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentUploadFileLogic {
	return &AttachmentUploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentUploadFileLogic) AttachmentUploadFile(req *types.AttachmentUploadReq, fileName, contentType string, content []byte) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.AttachmentUploadFile(l.ctx, &sysservice.AttachmentUploadReq{
		EnvCode:     req.EnvCode,
		FileName:    fileName,
		ContentType: contentType,
		Content:     content,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data}, nil
}
