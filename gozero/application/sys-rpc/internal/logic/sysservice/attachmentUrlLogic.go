package sysservicelogic

import (
	"context"
	"time"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttachmentUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentUrlLogic {
	return &AttachmentUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AttachmentUrlLogic) AttachmentUrl(in *pb.AttachmentUrlQueryReq) (*pb.AttachmentUrlResp, error) {
	row, err := getAttachmentByID(l.ctx, l.svcCtx, in.AttachmentId)
	if err != nil {
		return nil, err
	}
	expires := in.Expires
	if expires <= 0 {
		expires = 3600
	}
	url := ""
	if row.AccessUrl.Valid && row.AccessUrl.String != "" {
		url = row.AccessUrl.String
	} else {
		generatedURL, err := getAttachmentAccessURL(l.ctx, l.svcCtx, row, time.Duration(expires)*time.Second)
		if err == nil {
			url = generatedURL
		}
	}
	return &pb.AttachmentUrlResp{AttachmentId: in.AttachmentId, Url: url, Expires: expires}, nil
}
