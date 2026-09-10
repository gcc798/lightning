// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"context"
	"encoding/json"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentBindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentBindLogic {
	return &AttachmentBindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentBindLogic) AttachmentBind(req *types.AttachmentBindReq) (resp *types.CommonResp, err error) {
	metadataJSON := ""
	if len(req.Metadata) > 0 {
		buf, err := json.Marshal(req.Metadata)
		if err != nil {
			return &types.CommonResp{Code: 400, Msg: "metadata 格式错误"}, nil
		}
		metadataJSON = string(buf)
	}
	if _, err := l.svcCtx.SysRpcClient.AttachmentBind(l.ctx, &sysservice.AttachmentBindReq{
		AttachmentId:  req.AttachmentId,
		BusinessType:  req.BusinessType,
		BusinessId:    req.BusinessId,
		BusinessField: req.BusinessField,
		IsPublic:      req.IsPublic,
		MetadataJson:  metadataJSON,
		ExpireTime:    req.ExpireTime,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
