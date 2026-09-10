package sysservicelogic

import (
	"context"
	"errors"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgBatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgBatchDeleteLogic {
	return &OrgBatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgBatchDeleteLogic) OrgBatchDelete(in *pb.BatchIdsReq) (*pb.Ack, error) {
	if len(in.Ids) == 0 {
		return nil, errors.New("请至少选择一条数据")
	}
	failures := make([]string, 0)
	for _, id := range in.Ids {
		if _, err := NewOrgDeleteLogic(l.ctx, l.svcCtx).OrgDelete(&pb.IdReq{Id: id}); err != nil {
			failures = append(failures, err.Error())
		}
	}
	if len(failures) > 0 {
		return nil, errors.New(strings.Join(failures, "; "))
	}
	return &pb.Ack{Msg: "ok"}, nil
}
