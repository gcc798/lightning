package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttachmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDeleteLogic {
	return &AttachmentDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AttachmentDeleteLogic) AttachmentDelete(in *pb.IdReq) (*pb.Ack, error) {
	row, err := getAttachmentByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.biz_attachment where id = $1`, in.Id); err != nil {
		return nil, err
	}
	if row.FileKey != "" {
		s, err := l.svcCtx.StorageManager.GetStorage(l.ctx, row.EnvId)
		if err == nil {
			_ = s.Delete(l.ctx, row.FileKey)
		}
	}
	return &pb.Ack{Msg: "ok"}, nil
}
