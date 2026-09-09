package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/biz"
)

type LoginLogServiceService struct {
	v1.UnimplementedLoginLogServiceServer
	uc *biz.LoginLogUsecase
}

func NewLoginLogServiceService(uc *biz.LoginLogUsecase) *LoginLogServiceService {
	return &LoginLogServiceService{uc: uc}
}
func (s *LoginLogServiceService) CreateLoginLog(ctx context.Context, req *v1.CreateLoginLogRequest) (*v1.MessageReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *LoginLogServiceService) PageLoginLog(ctx context.Context, req *v1.PageLoginLogRequest) (*v1.PageLogReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *LoginLogServiceService) BatchDeleteLoginLog(ctx context.Context, req *v1.LogBatchIdsRequest) (*v1.MessageReply, error) {
	return s.uc.BatchDelete(ctx, req.GetIds())
}
func (s *LoginLogServiceService) CleanLoginLog(ctx context.Context, req *v1.CleanLogRequest) (*v1.LogCleanReply, error) {
	return s.uc.Clean(ctx, req.GetDays())
}
func (s *LoginLogServiceService) UpdateLoginLog(ctx context.Context, req *v1.UpdateLoginLogRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *LoginLogServiceService) GetLoginLogById(ctx context.Context, req *v1.LogIdRequest) (*v1.LogItem, error) {
	return s.uc.Get(ctx, req.GetId())
}
func (s *LoginLogServiceService) DeleteLoginLog(ctx context.Context, req *v1.LogIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
