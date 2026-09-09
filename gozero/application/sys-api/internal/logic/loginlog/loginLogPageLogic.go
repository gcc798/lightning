package loginlog

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogPageLogic {
	return &LoginLogPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *LoginLogPageLogic) LoginLogPage(req *types.LoginLogPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.LoginLogPage(l.ctx, &sysservice.LoginLogPageReq{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		UserName:  req.UserName,
		Ipaddr:    req.Ipaddr,
		Status:    int32(req.Status),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size)}, nil
}
