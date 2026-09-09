package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/biz"
)

type OperLogServiceService struct {
	v1.UnimplementedOperLogServiceServer
	uc *biz.OperLogUsecase
}

func NewOperLogServiceService(uc *biz.OperLogUsecase) *OperLogServiceService {
	return &OperLogServiceService{uc: uc}
}
func (s *OperLogServiceService) CreateOperLog(ctx context.Context, req *v1.CreateOperLogRequest) (*v1.MessageReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *OperLogServiceService) PageOperLog(ctx context.Context, req *v1.PageOperLogRequest) (*v1.PageLogReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *OperLogServiceService) BatchDeleteOperLog(ctx context.Context, req *v1.LogBatchIdsRequest) (*v1.MessageReply, error) {
	return s.uc.BatchDelete(ctx, req.GetIds())
}
func (s *OperLogServiceService) CleanOperLog(ctx context.Context, req *v1.CleanLogRequest) (*v1.LogCleanReply, error) {
	return s.uc.Clean(ctx, req.GetDays())
}
func (s *OperLogServiceService) UpdateOperLog(ctx context.Context, req *v1.UpdateOperLogRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *OperLogServiceService) GetOperLogById(ctx context.Context, req *v1.LogIdRequest) (*v1.LogItem, error) {
	return s.uc.Get(ctx, req.GetId())
}
func (s *OperLogServiceService) DeleteOperLog(ctx context.Context, req *v1.LogIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
