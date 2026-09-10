package sysservicelogic

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogBatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogBatchDeleteLogic {
	return &LoginLogBatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogBatchDeleteLogic) LoginLogBatchDelete(in *pb.BatchIdsReq) (*pb.Ack, error) {
	if len(in.Ids) == 0 {
		return nil, errors.New("ids 不能为空")
	}
	placeholders, args := loginLogIn(in.Ids, 1)
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, fmt.Sprintf(`delete from public.s_login_log where id in (%s)`, placeholders), args...); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
