package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type CaptchaServiceService struct {
	pb.UnimplementedCaptchaServiceServer
}

func NewCaptchaServiceService() *CaptchaServiceService {
	return &CaptchaServiceService{}
}

func (s *CaptchaServiceService) GenerateImageCaptcha(ctx context.Context, req *pb.GenerateImageCaptchaRequest) (*pb.CaptchaReply, error) {
	return &pb.CaptchaReply{}, nil
}
func (s *CaptchaServiceService) SendSMSCaptcha(ctx context.Context, req *pb.SendSMSCaptchaRequest) (*pb.CaptchaReply, error) {
	return &pb.CaptchaReply{}, nil
}
func (s *CaptchaServiceService) SendEmailCaptcha(ctx context.Context, req *pb.SendEmailCaptchaRequest) (*pb.CaptchaReply, error) {
	return &pb.CaptchaReply{}, nil
}
func (s *CaptchaServiceService) GetEnabledTypes(ctx context.Context, req *pb.GetEnabledTypesRequest) (*pb.GetEnabledTypesReply, error) {
	return &pb.GetEnabledTypesReply{}, nil
}
