// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangePasswordLogic {
	return &UserChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChangePasswordLogic) UserChangePassword(req *types.UserChangePasswordReq) (resp *types.CommonResp, err error) {
	userID := commonutil.UserIDFromContext(l.ctx)
	if userID <= 0 {
		return &types.CommonResp{Code: 401, Msg: "未登录"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.UserChangePassword(l.ctx, &sysservice.UserChangePasswordReq{
		UserId:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
