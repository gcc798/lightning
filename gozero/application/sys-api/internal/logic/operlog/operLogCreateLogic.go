package operlog

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCreateLogic {
	return &OperLogCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *OperLogCreateLogic) OperLogCreate(req *types.OperLogReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.OperLogCreate(l.ctx, &sysservice.OperLogReq{
		Title:         req.Title,
		BusinessType:  req.BusinessType,
		Method:        req.Method,
		RequestMethod: req.RequestMethod,
		DeviceType:    req.DeviceType,
		OperName:      req.OperName,
		OperUrl:       req.OperUrl,
		OperIp:        req.OperIp,
		OperLocation:  req.OperLocation,
		OperParam:     req.OperParam,
		JsonResult:    req.JsonResult,
		Status:        req.Status,
		ErrorMsg:      req.ErrorMsg,
		CostTime:      req.CostTime,
		UserAgent:     req.UserAgent,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
