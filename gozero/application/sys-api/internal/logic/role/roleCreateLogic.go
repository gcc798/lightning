package role

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RoleCreateLogic) RoleCreate(req *types.RoleCreateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RoleCreate(l.ctx, &sysservice.RoleCreateReq{
		RoleKey:   req.RoleKey,
		RoleName:  req.RoleName,
		Sort:      req.Sort,
		Status:    int32(req.Status),
		DataScope: int32(req.DataScope),
		Remark:    req.Remark,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: nil}, nil
}
