// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentPageLogic {
	return &AttachmentPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentPageLogic) AttachmentPage(req *types.AttachmentPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.AttachmentPage(l.ctx, &sysservice.AttachmentPageReq{
		PageNum:      req.PageNum,
		PageSize:     req.PageSize,
		FileName:     req.FileName,
		FileType:     req.FileType,
		BusinessType: req.BusinessType,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size)}, nil
}
