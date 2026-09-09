package service

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/biz"
)

type OrgServiceService struct {
	v1.UnimplementedOrgServiceServer
	uc *biz.OrgUsecase
}

func NewOrgServiceService(uc *biz.OrgUsecase) *OrgServiceService { return &OrgServiceService{uc: uc} }
func (s *OrgServiceService) CreateOrg(ctx context.Context, req *v1.CreateOrgRequest) (*v1.OrgIdReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *OrgServiceService) UpdateOrg(ctx context.Context, req *v1.UpdateOrgRequest) (*v1.MessageReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *OrgServiceService) DeleteOrg(ctx context.Context, req *v1.OrgIdRequest) (*v1.MessageReply, error) {
	return s.uc.Delete(ctx, req.GetId())
}
func (s *OrgServiceService) BatchDeleteOrg(ctx context.Context, req *v1.BatchDeleteOrgsRequest) (*v1.MessageReply, error) {
	return s.uc.BatchDelete(ctx, req.GetIds())
}
func (s *OrgServiceService) GetOrgById(ctx context.Context, req *v1.OrgIdRequest) (*v1.OrgItem, error) {
	return s.uc.GetByID(ctx, req.GetId())
}
func (s *OrgServiceService) GetOrgTree(ctx context.Context, req *v1.GetOrgTreeRequest) (*v1.OrgTreeReply, error) {
	return s.uc.Tree(ctx)
}
func (s *OrgServiceService) PageOrg(ctx context.Context, req *v1.PageOrgsRequest) (*v1.PageOrgsReply, error) {
	return s.uc.Page(ctx, req)
}
