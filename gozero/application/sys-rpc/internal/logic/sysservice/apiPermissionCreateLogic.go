package sysservicelogic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiPermissionCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiPermissionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionCreateLogic {
	return &ApiPermissionCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiPermissionCreateLogic) ApiPermissionCreate(in *pb.ApiPermissionSaveReq) (*pb.ApiPermission, error) {
	in.Method = strings.ToUpper(in.Method)
	in.Action = normalizeApiPermissionAction(in.Code, in.Action)
	if err := validateApiPermission(l.ctx, l.svcCtx, in); err != nil {
		return nil, err
	}
	var id int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &id, `
		insert into public.s_api_permission (parent_id, module, code, name, node_type, action, method, path, sort, status, remark, create_by, update_by, created_time, updated_time)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, nullif($12, 0), nullif($12, 0), now(), now())
		returning id
	`,
		in.ParentId, in.Module, in.Code, in.Name, in.NodeType, in.Action,
		sql.NullString{String: in.Method, Valid: in.Method != ""},
		sql.NullString{String: in.Path, Valid: in.Path != ""},
		in.Sort, in.Status,
		sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		in.UserId,
	); err != nil {
		return nil, err
	}
	row, err := getApiPermissionByID(l.ctx, l.svcCtx, id)
	if err != nil {
		return nil, err
	}
	return toApiPermissionPB(*row), nil
}
