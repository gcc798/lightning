// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuUpdateLogic) MenuUpdate(req *types.MenuUpdateReq) (resp *types.CommonResp, err error) {
	userID := commonutil.UserIDFromContext(l.ctx)
	if _, err := l.svcCtx.SysRpcClient.MenuUpdate(l.ctx, &sysservice.MenuReq{
		Id:        req.Id,
		MenuName:  req.MenuName,
		ParentId:  req.ParentId,
		Sort:      req.Sort,
		Path:      req.Path,
		Component: req.Component,
		Query:     req.Query,
		IsFrame:   int64(req.IsFrame),
		IsCache:   int64(req.IsCache),
		MenuType:  int64(req.MenuType),
		Visible:   int64(req.Visible),
		Status:    int64(req.Status),
		Perms:     req.Perms,
		Icon:      req.Icon,
		Remark:    req.Remark,
		UpdateBy:  userID,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
