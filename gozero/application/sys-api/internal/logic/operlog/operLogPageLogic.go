package operlog

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperLogPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogPageLogic {
	return &OperLogPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *OperLogPageLogic) OperLogPage(req *types.OperLogPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.OperLogPage(l.ctx, &sysservice.OperLogPageReq{
		PageNum:      req.PageNum,
		PageSize:     req.PageSize,
		Title:        req.Title,
		OperName:     req.OperName,
		BusinessType: req.BusinessType,
		Status:       req.Status,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size)}, nil
}
