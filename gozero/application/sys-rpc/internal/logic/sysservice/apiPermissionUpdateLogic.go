package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiPermissionUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiPermissionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionUpdateLogic {
	return &ApiPermissionUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiPermissionUpdateLogic) ApiPermissionUpdate(in *pb.ApiPermissionSaveReq) (*pb.Ack, error) {
	if in.Id <= 0 {
		return nil, fmt.Errorf("权限ID不能为空")
	}
	if _, err := getApiPermissionByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	in.Method = strings.ToUpper(in.Method)
	in.Action = normalizeApiPermissionAction(in.Code, in.Action)
	if err := validateApiPermission(l.ctx, l.svcCtx, in); err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		update public.s_api_permission
		set parent_id = $2, module = $3, code = $4, name = $5, node_type = $6, action = $7, method = $8, path = $9, sort = $10, status = $11, remark = $12, update_by = nullif($13, 0), updated_time = now()
		where id = $1
	`,
		in.Id, in.ParentId, in.Module, in.Code, in.Name, in.NodeType, in.Action,
		sql.NullString{String: in.Method, Valid: in.Method != ""},
		sql.NullString{String: in.Path, Valid: in.Path != ""},
		in.Sort, in.Status,
		sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		in.UserId,
	); err != nil {
		return nil, err
	}
	roleIDs, userIDs, err := findAffectedPermissionSubjects(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		l.Errorf("查询受影响主体失败: %v", err)
	}
	for _, roleID := range roleIDs {
		if err := syncRolePermissionRedis(l.ctx, l.svcCtx, roleID); err != nil {
			l.Errorf("同步角色权限缓存失败 roleID=%d: %v", roleID, err)
		}
	}
	for _, userID := range userIDs {
		if err := syncUserPermissionRedis(l.ctx, l.svcCtx, userID); err != nil {
			l.Errorf("同步用户权限缓存失败 userID=%d: %v", userID, err)
		}
	}
	return &pb.Ack{Msg: "ok"}, nil
}
