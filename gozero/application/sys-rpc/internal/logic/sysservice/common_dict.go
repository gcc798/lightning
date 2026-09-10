package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"
	gzsqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type dictRow struct {
	Id          int64          `db:"id"`
	ParentId    int64          `db:"parent_id"`
	DictType    sql.NullString `db:"dict_type"`
	DictLabel   sql.NullString `db:"dict_label"`
	DictValue   sql.NullString `db:"dict_value"`
	Sort        int64          `db:"sort"`
	IsDefault   bool           `db:"is_default"`
	Status      int64          `db:"status"`
	Remark      sql.NullString `db:"remark"`
	CreateBy    sql.NullInt64  `db:"create_by"`
	UpdateBy    sql.NullInt64  `db:"update_by"`
	CreatedTime sql.NullTime   `db:"created_time"`
	UpdatedTime sql.NullTime   `db:"updated_time"`
}

func parseDictID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func getDictByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*dictRow, error) {
	var row dictRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, parent_id, dict_type, dict_label, dict_value, sort, is_default, status, remark, create_by, update_by, created_time, updated_time
		from public.s_dict_data
		where id = $1
		limit 1
	`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("字典不存在")
		}
		return nil, err
	}
	return &row, nil
}

func dictValueExists(ctx context.Context, svcCtx *svc.ServiceContext, dictType, dictValue string, excludeID int64) (bool, error) {
	query := `select count(1) from public.s_dict_data where dict_type = $1 and dict_value = $2`
	args := []interface{}{dictType, dictValue}
	if excludeID > 0 {
		query += ` and id <> $3`
		args = append(args, excludeID)
	}
	var count int64
	if err := svcCtx.DB.QueryRowCtx(ctx, &count, query, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

func dictDescendantExists(ctx context.Context, svcCtx *svc.ServiceContext, id, parentID int64) (bool, error) {
	var count int64
	err := svcCtx.DB.QueryRowCtx(ctx, &count, `
		with recursive dict_tree as (
			select id, parent_id from public.s_dict_data where id = $1
			union all
			select d.id, d.parent_id from public.s_dict_data d
			inner join dict_tree dt on d.parent_id = dt.id
		)
		select count(1) from dict_tree where id = $2
	`, id, parentID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func buildInt64In(ids []int64, start int) (string, []interface{}) {
	parts := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		parts = append(parts, fmt.Sprintf("$%d", start+i))
		args = append(args, id)
	}
	return strings.Join(parts, ", "), args
}

func toDictPB(row dictRow) *pb.Dict {
	return &pb.Dict{
		Id:          row.Id,
		ParentId:    row.ParentId,
		DictType:    nullString(row.DictType),
		DictLabel:   nullString(row.DictLabel),
		DictValue:   nullString(row.DictValue),
		Sort:        row.Sort,
		IsDefault:   row.IsDefault,
		Status:      int32(row.Status),
		Remark:      nullString(row.Remark),
		CreateBy:    nullInt64(row.CreateBy),
		UpdateBy:    nullInt64(row.UpdateBy),
		CreatedTime: nullTime(row.CreatedTime),
		UpdatedTime: nullTime(row.UpdatedTime),
	}
}

func toDictList(rows []dictRow) []*pb.Dict {
	list := make([]*pb.Dict, 0, len(rows))
	for _, row := range rows {
		list = append(list, toDictPB(row))
	}
	return list
}
