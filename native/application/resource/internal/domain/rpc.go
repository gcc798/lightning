package resource

import (
	"context"

	resourcev1 "github.com/gcc798/lightning/internal/api/resource/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API struct{ attachments AttachmentService }

func NewAPI(attachments AttachmentService) *API { return &API{attachments: attachments} }

func (a *API) CleanExpired(ctx context.Context) (*resourcev1.CleanExpiredResponse, error) {
	cleaned, failed, err := a.attachments.CleanExpired(ctx)
	return &resourcev1.CleanExpiredResponse{Cleaned: cleaned, Failed: failed}, err
}

type GRPCServer struct {
	resourcev1.UnimplementedResourceServiceServer
	api resourcev1.API
}

func NewGRPCServer(api resourcev1.API) *GRPCServer { return &GRPCServer{api: api} }

func (s *GRPCServer) CleanExpired(ctx context.Context, _ *resourcev1.CleanExpiredRequest) (*resourcev1.CleanExpiredResponse, error) {
	result, err := s.api.CleanExpired(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return result, nil
}
