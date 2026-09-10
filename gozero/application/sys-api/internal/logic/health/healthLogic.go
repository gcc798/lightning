// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthLogic) Health() (resp *types.CommonResp, err error) {
	pingResp, pingErr := l.svcCtx.SysRpcClient.Ping(l.ctx, &sysservice.PingReq{})
	if pingErr != nil {
		return &types.CommonResp{
			Code: 500,
			Msg:  pingErr.Error(),
		}, nil
	}

	return &types.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: map[string]interface{}{
			"status": "ok",
			"rpc":    pingResp.Message,
		},
	}, nil
}
