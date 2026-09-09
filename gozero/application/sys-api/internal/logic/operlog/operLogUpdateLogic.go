package operlog

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperLogUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogUpdateLogic {
	return &OperLogUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *OperLogUpdateLogic) OperLogUpdate(req *types.OperLogUpdateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.OperLogUpdate(l.ctx, &sysservice.OperLogUpdateReq{
		Id:            req.Id,
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
