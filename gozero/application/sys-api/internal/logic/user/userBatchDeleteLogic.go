package user

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserBatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBatchDeleteLogic {
	return &UserBatchDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *UserBatchDeleteLogic) UserBatchDelete(req *types.BatchIdsReq) (resp *types.CommonResp, err error) {
	if len(req.Ids) == 0 {
		return &types.CommonResp{Code: 400, Msg: "请至少选择一条数据"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.UserBatchDelete(l.ctx, &sysservice.BatchIdsReq{Ids: req.Ids}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功"}, nil
}
