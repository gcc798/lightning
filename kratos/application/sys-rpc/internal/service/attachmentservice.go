package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/biz"
)

type AttachmentServiceService struct {
	v1.UnimplementedAttachmentServiceServer
	uc *biz.AttachmentUsecase
}

func NewAttachmentServiceService(uc *biz.AttachmentUsecase) *AttachmentServiceService {
	return &AttachmentServiceService{uc: uc}
}
func (s *AttachmentServiceService) UploadFile(ctx context.Context, req *v1.UploadFileRequest) (*v1.AttachmentItem, error) {
	return s.uc.Upload(ctx, req)
}
func (s *AttachmentServiceService) BindAttachmentToBusiness(ctx context.Context, req *v1.BindAttachmentRequest) (*v1.MessageReply, error) {
	return s.uc.Bind(ctx, req)
}
func (s *AttachmentServiceService) GetAttachment(ctx context.Context, req *v1.AttachmentIdRequest) (*v1.AttachmentItem, error) {
	return s.uc.Get(ctx, req.GetAttachmentId())
}
func (s *AttachmentServiceService) ListAttachmentsByBusiness(ctx context.Context, req *v1.ListAttachmentsByBusinessRequest) (*v1.AttachmentListReply, error) {
	return s.uc.List(ctx, req.GetBusinessType(), req.GetBusinessId())
}
func (s *AttachmentServiceService) PageAttachments(ctx context.Context, req *v1.PageAttachmentsRequest) (*v1.PageAttachmentsReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *AttachmentServiceService) DownloadAttachment(ctx context.Context, req *v1.AttachmentIdRequest) (*v1.AttachmentDownloadReply, error) {
	return s.uc.Download(ctx, req.GetAttachmentId())
}
func (s *AttachmentServiceService) GetAttachmentURL(ctx context.Context, req *v1.AttachmentURLRequest) (*v1.AttachmentURLReply, error) {
	return s.uc.URL(ctx, req.GetAttachmentId(), req.GetExpires())
}
func (s *AttachmentServiceService) DeleteAttachment(ctx context.Context, req *v1.AttachmentIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetAttachmentId())
}
