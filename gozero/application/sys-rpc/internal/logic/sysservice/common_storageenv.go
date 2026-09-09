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

type storageEnvRow struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Code        string         `db:"code"`
	StorageType sql.NullString `db:"storage_type"`
	IsDefault   bool           `db:"is_default"`
	Status      int64          `db:"status"`
	Config      string         `db:"config"`
	Remark      sql.NullString `db:"remark"`
	CreateBy    sql.NullInt64  `db:"create_by"`
	CreatedTime sql.NullTime   `db:"created_time"`
	UpdateBy    sql.NullInt64  `db:"update_by"`
	UpdatedTime sql.NullTime   `db:"updated_time"`
}

func getStorageEnvByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*storageEnvRow, error) {
	var row storageEnvRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, name, code, storage_type, is_default, status, config, remark, create_by, created_time, update_by, updated_time
		from public.s_storage_env
		where id = $1
		limit 1
	`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("存储环境不存在")
		}
		return nil, err
	}
	return &row, nil
}

func getStorageEnvByCode(ctx context.Context, svcCtx *svc.ServiceContext, code string) (*storageEnvRow, error) {
	var row storageEnvRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, name, code, storage_type, is_default, status, config, remark, create_by, created_time, update_by, updated_time
		from public.s_storage_env
		where code = $1
		limit 1
	`, code)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("存储环境不存在")
		}
		return nil, err
	}
	return &row, nil
}

func getDefaultStorageEnv(ctx context.Context, svcCtx *svc.ServiceContext) (*storageEnvRow, error) {
	var row storageEnvRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, name, code, storage_type, is_default, status, config, remark, create_by, created_time, update_by, updated_time
		from public.s_storage_env
		where is_default = true and status = 0
		order by id desc
		limit 1
	`)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("默认存储环境不存在")
		}
		return nil, err
	}
	return &row, nil
}

func storageEnvCodeExists(ctx context.Context, svcCtx *svc.ServiceContext, code string, excludeID int64) (bool, error) {
	query := `select count(1) from public.s_storage_env where code = $1`
	args := []interface{}{code}
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

func toStorageEnvPB(row storageEnvRow) *pb.StorageEnv {
	return &pb.StorageEnv{
		Id:          row.Id,
		Name:        row.Name,
		Code:        row.Code,
		StorageType: nullString(row.StorageType),
		IsDefault:   row.IsDefault,
		Status:      int32(row.Status),
		ConfigJson:  row.Config,
		Remark:      nullString(row.Remark),
		CreateBy:    nullInt64(row.CreateBy),
		CreatedAt:   nullTime(row.CreatedTime),
		UpdateBy:    nullInt64(row.UpdateBy),
		UpdatedAt:   nullTime(row.UpdatedTime),
	}
}

func toStorageEnvList(rows []storageEnvRow) []*pb.StorageEnv {
	list := make([]*pb.StorageEnv, 0, len(rows))
	for _, row := range rows {
		list = append(list, toStorageEnvPB(row))
	}
	return list
}

func buildStorageInt64In(ids []int64, start int) (string, []interface{}) {
	parts := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		parts = append(parts, fmt.Sprintf("$%d", start+i))
		args = append(args, id)
	}
	return strings.Join(parts, ", "), args
}
