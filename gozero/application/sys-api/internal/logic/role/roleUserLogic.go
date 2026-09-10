package role

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUserLogic {
	return &RoleUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RoleUserLogic) RoleUser(req *types.UserRoleQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.RoleUser(l.ctx, &sysservice.UserRoleQueryReq{UserId: req.UserId})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.Records}, nil
}
