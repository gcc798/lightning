package biz

import (
	"context"
	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/data"
)

type AttachmentUsecase struct{ repo *data.AttachmentRepo }

func NewAttachmentUsecase(repo *data.AttachmentRepo) *AttachmentUsecase {
	return &AttachmentUsecase{repo: repo}
}
func (uc *AttachmentUsecase) Upload(ctx context.Context, req *v1.UploadFileRequest) (*v1.AttachmentItem, error) {
	return uc.repo.Upload(ctx, req)
}
func (uc *AttachmentUsecase) Bind(ctx context.Context, req *v1.BindAttachmentRequest) (*v1.MessageReply, error) {
	return uc.repo.Bind(ctx, req)
}
func (uc *AttachmentUsecase) Get(ctx context.Context, id int64) (*v1.AttachmentItem, error) {
	return uc.repo.Get(ctx, id)
}
func (uc *AttachmentUsecase) List(ctx context.Context, bt, bid string) (*v1.AttachmentListReply, error) {
	return uc.repo.List(ctx, bt, bid)
}
func (uc *AttachmentUsecase) Page(ctx context.Context, req *v1.PageAttachmentsRequest) (*v1.PageAttachmentsReply, error) {
	return uc.repo.Page(ctx, req)
}
func (uc *AttachmentUsecase) Download(ctx context.Context, id int64) (*v1.AttachmentDownloadReply, error) {
	return uc.repo.Download(ctx, id)
}
func (uc *AttachmentUsecase) URL(ctx context.Context, id int64, expires int64) (*v1.AttachmentURLReply, error) {
	return uc.repo.URL(ctx, id, expires)
}
func (uc *AttachmentUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return uc.repo.Delete(ctx, id)
}
