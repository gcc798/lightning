package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentDownloadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttachmentDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDownloadLogic {
	return &AttachmentDownloadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AttachmentDownloadLogic) AttachmentDownload(in *pb.AttachmentDownloadReq) (*pb.AttachmentDownloadResp, error) {
	row, err := getAttachmentByID(l.ctx, l.svcCtx, in.AttachmentId)
	if err != nil {
		return nil, err
	}
	content, contentType, err := readAttachmentContent(l.ctx, l.svcCtx, row)
	if err != nil {
		return nil, err
	}
	return &pb.AttachmentDownloadResp{
		FileName:    row.FileName,
		ContentType: contentType,
		Content:     content,
	}, nil
}
