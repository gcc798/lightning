package service

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/biz"
)

type DictServiceService struct {
	v1.UnimplementedDictServiceServer
	uc *biz.DictUsecase
}

func NewDictServiceService(uc *biz.DictUsecase) *DictServiceService {
	return &DictServiceService{uc: uc}
}
func (s *DictServiceService) CreateDict(ctx context.Context, req *v1.DictItem) (*v1.MessageReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *DictServiceService) PageDict(ctx context.Context, req *v1.PageDictRequest) (*v1.PageDictReply, error) {
	return s.uc.Page(ctx, req)
}
func (s *DictServiceService) BatchDeleteDict(ctx context.Context, req *v1.DictBatchIdsRequest) (*v1.MessageReply, error) {
	return s.uc.BatchDelete(ctx, req.GetIds())
}
func (s *DictServiceService) GetDictByType(ctx context.Context, req *v1.GetDictByTypeRequest) (*v1.DictListReply, error) {
	return s.uc.ByType(ctx, req.GetDictType(), req.ParentId)
}
func (s *DictServiceService) GetDictLabel(ctx context.Context, req *v1.GetDictLabelRequest) (*v1.DictLabelReply, error) {
	return s.uc.Label(ctx, req.GetDictType(), req.GetDictValue())
}
func (s *DictServiceService) UpdateDict(ctx context.Context, req *v1.UpdateDictRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *DictServiceService) GetDictById(ctx context.Context, req *v1.DictIdRequest) (*v1.DictItem, error) {
	return s.uc.Get(ctx, req.GetId())
}
func (s *DictServiceService) DeleteDict(ctx context.Context, req *v1.DictIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
