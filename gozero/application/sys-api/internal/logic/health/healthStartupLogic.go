// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthStartupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthStartupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthStartupLogic {
	return &HealthStartupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthStartupLogic) HealthStartup() (resp *types.CommonResp, err error) {
	return NewHealthLogic(l.ctx, l.svcCtx).Health()
}
