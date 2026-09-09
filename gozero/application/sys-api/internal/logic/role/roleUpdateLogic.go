package role

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RoleUpdateLogic) RoleUpdate(req *types.RoleUpdateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RoleUpdate(l.ctx, &sysservice.RoleUpdateReq{
		RoleId:    req.RoleId,
		RoleName:  req.RoleName,
		Sort:      req.Sort,
		Status:    int32(req.Status),
		DataScope: int32(req.DataScope),
		Remark:    req.Remark,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
