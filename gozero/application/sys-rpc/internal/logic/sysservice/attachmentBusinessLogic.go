package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentBusinessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttachmentBusinessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentBusinessLogic {
	return &AttachmentBusinessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AttachmentBusinessLogic) AttachmentBusiness(in *pb.AttachmentBusinessQueryReq) (*pb.AttachmentListResp, error) {
	var rows []attachmentRow
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, `
		select id, env_id, file_name, file_key, file_size, file_type, file_ext, business_type, business_id, business_field, is_public, access_url, metadata, status, expire_time, create_by, create_time, update_time
		from public.biz_attachment
		where business_type = $1 and business_id = $2 and status = 0
		order by create_time desc
	`, in.BusinessType, in.BusinessId); err != nil {
		return nil, err
	}
	return &pb.AttachmentListResp{Records: toAttachmentList(rows)}, nil
}
