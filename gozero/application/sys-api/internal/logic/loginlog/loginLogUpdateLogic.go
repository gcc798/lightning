package loginlog

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogUpdateLogic {
	return &LoginLogUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *LoginLogUpdateLogic) LoginLogUpdate(req *types.LoginLogUpdateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.LoginLogUpdate(l.ctx, &sysservice.LoginLogUpdateReq{
		Id:            req.Id,
		UserName:      req.UserName,
		Ipaddr:        req.Ipaddr,
		LoginLocation: req.LoginLocation,
		Browser:       req.Browser,
		Os:            req.Os,
		Status:        int32(req.Status),
		Msg:           req.Msg,
		ClientId:      req.ClientId,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
