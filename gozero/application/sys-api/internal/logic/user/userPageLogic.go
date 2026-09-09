package user

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPageLogic {
	return &UserPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *UserPageLogic) UserPage(req *types.UserPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.UserPage(l.ctx, &sysservice.UserPageReq{PageNum: req.PageNum, PageSize: req.PageSize, Username: req.Username, Phonenumber: req.Phonenumber, Status: int32(req.Status)})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size)}, nil
}
