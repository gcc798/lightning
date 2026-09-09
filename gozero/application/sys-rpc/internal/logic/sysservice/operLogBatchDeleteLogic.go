package sysservicelogic

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogBatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogBatchDeleteLogic {
	return &OperLogBatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogBatchDeleteLogic) OperLogBatchDelete(in *pb.BatchIdsReq) (*pb.Ack, error) {
	if len(in.Ids) == 0 {
		return nil, errors.New("ids 不能为空")
	}
	placeholders, args := operLogIn(in.Ids, 1)
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, fmt.Sprintf(`delete from public.s_oper_log where id in (%s)`, placeholders), args...); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
