// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthReadyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthReadyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthReadyLogic {
	return &HealthReadyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthReadyLogic) HealthReady() (resp *types.CommonResp, err error) {
	return NewHealthLogic(l.ctx, l.svcCtx).Health()
}
