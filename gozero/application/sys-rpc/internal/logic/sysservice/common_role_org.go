package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"
	gzsqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type roleRow struct {
	Id          int64          `db:"id"`
	RoleKey     string         `db:"role_key"`
	RoleName    string         `db:"role_name"`
	Sort        int64          `db:"sort"`
	Status      int64          `db:"status"`
	DataScope   int64          `db:"data_scope"`
	IsSystem    bool           `db:"is_system"`
	Remark      sql.NullString `db:"remark"`
	CreateBy    sql.NullInt64  `db:"create_by"`
	CreatedTime sql.NullTime   `db:"created_time"`
}

type orgRow struct {
	Id          int64          `db:"id"`
	ParentId    int64          `db:"parent_id"`
	Ancestors   sql.NullString `db:"ancestors"`
	OrgName     string         `db:"org_name"`
	OrgCode     sql.NullString `db:"org_code"`
	OrgType     string         `db:"org_type"`
	Leader      sql.NullString `db:"leader"`
	Phone       sql.NullString `db:"phone"`
	Email       sql.NullString `db:"email"`
	Status      int64          `db:"status"`
	Sort        int64          `db:"sort"`
	Remark      sql.NullString `db:"remark"`
	CreatedTime sql.NullTime   `db:"created_time"`
	UpdatedTime sql.NullTime   `db:"updated_time"`
}

func getRoleByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*roleRow, error) {
	var row roleRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, role_key, role_name, sort, status, data_scope, is_system, remark, create_by, created_time
		from public.s_role
		where id = $1
		limit 1
	`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, err
	}
	return &row, nil
}

func toRolePB(row roleRow) *pb.Role {
	return &pb.Role{
		RoleId:     row.Id,
		RoleKey:    row.RoleKey,
		RoleName:   row.RoleName,
		Sort:       row.Sort,
		Status:     int32(row.Status),
		DataScope:  int32(row.DataScope),
		IsSystem:   row.IsSystem,
		Remark:     nullString(row.Remark),
		CreateBy:   nullInt64(row.CreateBy),
		CreateTime: nullTime(row.CreatedTime),
	}
}

func toRoleList(rows []roleRow) []*pb.Role {
	list := make([]*pb.Role, 0, len(rows))
	for _, row := range rows {
		list = append(list, toRolePB(row))
	}
	return list
}

func getOrgByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*orgRow, error) {
	var row orgRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, parent_id, ancestors, org_name, org_code, org_type, leader, phone, email, status, sort, remark, created_time, updated_time
		from public.s_org
		where id = $1
		limit 1
	`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("组织不存在")
		}
		return nil, err
	}
	return &row, nil
}

func toOrgPB(row orgRow) *pb.Org {
	return &pb.Org{
		Id:          row.Id,
		OrgId:       row.Id,
		ParentId:    row.ParentId,
		OrgName:     row.OrgName,
		OrgCode:     nullString(row.OrgCode),
		OrgType:     row.OrgType,
		Leader:      nullString(row.Leader),
		Phone:       nullString(row.Phone),
		Email:       nullString(row.Email),
		Status:      int32(row.Status),
		Sort:        row.Sort,
		Remark:      nullString(row.Remark),
		CreatedTime: nullTime(row.CreatedTime),
		UpdatedTime: nullTime(row.UpdatedTime),
	}
}

func toOrgList(rows []orgRow) []*pb.Org {
	list := make([]*pb.Org, 0, len(rows))
	for _, row := range rows {
		list = append(list, toOrgPB(row))
	}
	return list
}

func buildOrgTree(rows []orgRow, parentID int64) []*pb.Org {
	tree := make([]*pb.Org, 0)
	for _, row := range rows {
		if row.ParentId != parentID {
			continue
		}
		node := toOrgPB(row)
		node.Children = buildOrgTree(rows, row.Id)
		tree = append(tree, node)
	}
	return tree
}

func buildOrgAncestors(ctx context.Context, svcCtx *svc.ServiceContext, parentID int64) (string, error) {
	if parentID == 0 {
		return "0", nil
	}
	parent, err := getOrgByID(ctx, svcCtx, parentID)
	if err != nil {
		return "", fmt.Errorf("父组织不存在")
	}
	base := strings.TrimSpace(nullString(parent.Ancestors))
	if base == "" {
		base = "0"
	}
	return base + "," + fmt.Sprint(parentID), nil
}
