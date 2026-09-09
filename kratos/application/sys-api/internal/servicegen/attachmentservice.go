package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type AttachmentServiceService struct {
	pb.UnimplementedAttachmentServiceServer
}

func NewAttachmentServiceService() *AttachmentServiceService {
	return &AttachmentServiceService{}
}

func (s *AttachmentServiceService) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.AttachmentItem, error) {
	return &pb.AttachmentItem{}, nil
}
func (s *AttachmentServiceService) BindAttachmentToBusiness(ctx context.Context, req *pb.BindAttachmentRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *AttachmentServiceService) GetAttachment(ctx context.Context, req *pb.AttachmentIdRequest) (*pb.AttachmentItem, error) {
	return &pb.AttachmentItem{}, nil
}
func (s *AttachmentServiceService) ListAttachmentsByBusiness(ctx context.Context, req *pb.ListAttachmentsByBusinessRequest) (*pb.AttachmentListReply, error) {
	return &pb.AttachmentListReply{}, nil
}
func (s *AttachmentServiceService) PageAttachments(ctx context.Context, req *pb.PageAttachmentsRequest) (*pb.PageAttachmentsReply, error) {
	return &pb.PageAttachmentsReply{}, nil
}
func (s *AttachmentServiceService) DownloadAttachment(ctx context.Context, req *pb.AttachmentIdRequest) (*pb.AttachmentDownloadReply, error) {
	return &pb.AttachmentDownloadReply{}, nil
}
func (s *AttachmentServiceService) GetAttachmentURL(ctx context.Context, req *pb.AttachmentURLRequest) (*pb.AttachmentURLReply, error) {
	return &pb.AttachmentURLReply{}, nil
}
func (s *AttachmentServiceService) DeleteAttachment(ctx context.Context, req *pb.AttachmentIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
