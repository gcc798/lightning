package biz

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type AttachmentUsecase struct{ res *data.Resources }

func NewAttachmentUsecase(res *data.Resources) *AttachmentUsecase {
	return &AttachmentUsecase{res: res}
}

func (uc *AttachmentUsecase) Upload(ctx context.Context, req *v1.UploadFileRequest) (*v1.AttachmentItem, error) {
	return uc.res.SaveAttachment(ctx, req)
}

func (uc *AttachmentUsecase) Bind(ctx context.Context, req *v1.BindAttachmentRequest) (*v1.MessageReply, error) {
	if err := uc.res.BindAttachment(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *AttachmentUsecase) Get(ctx context.Context, id int64) (*v1.AttachmentItem, error) {
	item, err := uc.res.GetAttachment(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "附件不存在")
}

func (uc *AttachmentUsecase) List(ctx context.Context, businessType, businessID string) (*v1.AttachmentListReply, error) {
	items, err := uc.res.ListAttachments(ctx, businessType, businessID)
	if err != nil {
		return nil, err
	}
	return &v1.AttachmentListReply{Items: items}, nil
}

func (uc *AttachmentUsecase) Page(ctx context.Context, req *v1.PageAttachmentsRequest) (*v1.PageAttachmentsReply, error) {
	return uc.res.PageAttachments(ctx, req)
}

func (uc *AttachmentUsecase) Download(ctx context.Context, id int64) (*v1.AttachmentDownloadReply, error) {
	return uc.res.DownloadAttachment(ctx, id)
}

func (uc *AttachmentUsecase) URL(ctx context.Context, id int64, expires int64) (*v1.AttachmentURLReply, error) {
	reply, err := uc.res.GetAttachmentURL(ctx, id, expires)
	if err != nil {
		return nil, err
	}
	return requireFound(reply, "附件不存在")
}

func (uc *AttachmentUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteAttachments(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}
