package service

import (
	"context"

	pb "github.com/gcc798/lightning/kratos/api/system/v1"
)

type DictServiceService struct {
	pb.UnimplementedDictServiceServer
}

func NewDictServiceService() *DictServiceService {
	return &DictServiceService{}
}

func (s *DictServiceService) CreateDict(ctx context.Context, req *pb.DictItem) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *DictServiceService) PageDict(ctx context.Context, req *pb.PageDictRequest) (*pb.PageDictReply, error) {
	return &pb.PageDictReply{}, nil
}
func (s *DictServiceService) BatchDeleteDict(ctx context.Context, req *pb.DictBatchIdsRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *DictServiceService) GetDictByType(ctx context.Context, req *pb.GetDictByTypeRequest) (*pb.DictListReply, error) {
	return &pb.DictListReply{}, nil
}
func (s *DictServiceService) GetDictLabel(ctx context.Context, req *pb.GetDictLabelRequest) (*pb.DictLabelReply, error) {
	return &pb.DictLabelReply{}, nil
}
func (s *DictServiceService) UpdateDict(ctx context.Context, req *pb.UpdateDictRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
func (s *DictServiceService) GetDictById(ctx context.Context, req *pb.DictIdRequest) (*pb.DictItem, error) {
	return &pb.DictItem{}, nil
}
func (s *DictServiceService) DeleteDict(ctx context.Context, req *pb.DictIdRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}
