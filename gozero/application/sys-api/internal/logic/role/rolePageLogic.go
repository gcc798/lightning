// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRolePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePageLogic {
	return &RolePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RolePageLogic) RolePage(req *types.RolePageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.RolePage(l.ctx, &sysservice.RolePageReq{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		RoleName: req.RoleName,
		Status:   int32(req.Status),
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size),
	}, nil
}
