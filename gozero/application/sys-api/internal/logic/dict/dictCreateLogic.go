package dict

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictCreateLogic {
	return &DictCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *DictCreateLogic) DictCreate(req *types.DictCreateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.DictCreate(l.ctx, &sysservice.DictCreateReq{
		ParentId:  req.ParentId,
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		Sort:      req.Sort,
		IsDefault: req.IsDefault,
		Status:    int32(req.Status),
		Remark:    req.Remark,
		CreateBy:  req.CreateBy,
		UpdateBy:  req.UpdateBy,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
