package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type OrgServiceService struct {
	pb.UnimplementedOrgServiceServer
}

func NewOrgServiceService() *OrgServiceService {
	return &OrgServiceService{}
}

func (s *OrgServiceService) CreateOrg(ctx context.Context, req *pb.CreateOrgRequest) (*pb.OrgIdReply, error) {
	return &pb.OrgIdReply{}, nil
}
func (s *OrgServiceService) UpdateOrg(ctx context.Context, req *pb.UpdateOrgRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *OrgServiceService) DeleteOrg(ctx context.Context, req *pb.OrgIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *OrgServiceService) BatchDeleteOrg(ctx context.Context, req *pb.BatchDeleteOrgsRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *OrgServiceService) GetOrgById(ctx context.Context, req *pb.OrgIdRequest) (*pb.OrgItem, error) {
	return &pb.OrgItem{}, nil
}
func (s *OrgServiceService) GetOrgTree(ctx context.Context, req *pb.GetOrgTreeRequest) (*pb.OrgTreeReply, error) {
	return &pb.OrgTreeReply{}, nil
}
func (s *OrgServiceService) PageOrg(ctx context.Context, req *pb.PageOrgsRequest) (*pb.PageOrgsReply, error) {
	return &pb.PageOrgsReply{}, nil
}
