package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	gzsqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type configRow struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Code        string         `db:"code"`
	Data        sql.NullString `db:"data"`
	Remark      sql.NullString `db:"remark"`
	CreateBy    sql.NullInt64  `db:"create_by"`
	CreatedTime sql.NullTime   `db:"created_time"`
	UpdateBy    sql.NullInt64  `db:"update_by"`
	UpdatedTime sql.NullTime   `db:"updated_time"`
}

func getConfigByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*configRow, error) {
	var row configRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, name, code, data, remark, create_by, created_time, update_by, updated_time
		from public.s_config
		where id = $1
		limit 1
	`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("配置不存在")
		}
		return nil, err
	}
	return &row, nil
}

func configNameExists(ctx context.Context, svcCtx *svc.ServiceContext, name string, excludeID int64) (bool, error) {
	query := `select count(1) from public.s_config where name = $1`
	args := []interface{}{name}
	if excludeID > 0 {
		query += ` and id <> $2`
		args = append(args, excludeID)
	}
	var count int64
	if err := svcCtx.DB.QueryRowCtx(ctx, &count, query, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

func buildConfigInt64In(ids []int64, start int) (string, []interface{}) {
	parts := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		parts = append(parts, fmt.Sprintf("$%d", start+i))
		args = append(args, id)
	}
	return strings.Join(parts, ", "), args
}

func toConfigPB(row configRow) *pb.Config {
	return &pb.Config{
		Id:          row.Id,
		Name:        row.Name,
		Code:        row.Code,
		DataJson:    nullString(row.Data),
		Remark:      nullString(row.Remark),
		CreateBy:    nullInt64(row.CreateBy),
		CreatedTime: nullTime(row.CreatedTime),
		UpdateBy:    nullInt64(row.UpdateBy),
		UpdatedTime: nullTime(row.UpdatedTime),
	}
}

func toConfigList(rows []configRow) []*pb.Config {
	list := make([]*pb.Config, 0, len(rows))
	for _, row := range rows {
		list = append(list, toConfigPB(row))
	}
	return list
}
