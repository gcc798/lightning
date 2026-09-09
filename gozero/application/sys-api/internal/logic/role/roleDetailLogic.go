// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDetailLogic {
	return &RoleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDetailLogic) RoleDetail(req *types.RoleIdPathReq) (resp *types.CommonResp, err error) {
	row, err := l.svcCtx.SysRpcClient.RoleDetail(l.ctx, &sysservice.IdReq{Id: req.RoleId})
	if err != nil {
		return &types.CommonResp{Code: 404, Msg: err.Error()}, nil
	}
	return &types.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: row,
	}, nil
}
