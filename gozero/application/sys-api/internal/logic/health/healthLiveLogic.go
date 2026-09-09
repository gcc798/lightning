// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthLiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthLiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLiveLogic {
	return &HealthLiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthLiveLogic) HealthLive() (resp *types.CommonResp, err error) {
	return NewHealthLogic(l.ctx, l.svcCtx).Health()
}
