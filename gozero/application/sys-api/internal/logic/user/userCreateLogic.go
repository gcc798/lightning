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

type UserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateLogic) UserCreate(req *types.UserCreateReq) (resp *types.CommonResp, err error) {
	userID := commonutil.UserIDFromContext(l.ctx)
	if _, err := l.svcCtx.SysRpcClient.UserCreate(l.ctx, &sysservice.UserCreateReq{
		UserName:    req.UserName,
		NickName:    req.NickName,
		Password:    req.Password,
		UserType:    int32(req.UserType),
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Sex:         int32(req.Sex),
		Avatar:      req.Avatar,
		Status:      int32(req.Status),
		Remark:      req.Remark,
		CreateBy:    userID,
		UpdateBy:    userID,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: map[string]interface{}{"userId": "ok"}}, nil
}
