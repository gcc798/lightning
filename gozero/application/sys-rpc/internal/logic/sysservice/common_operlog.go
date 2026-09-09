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

type operLogRow struct {
	Id            int64          `db:"id"`
	Title         sql.NullString `db:"title"`
	BusinessType  sql.NullString `db:"business_type"`
	Method        sql.NullString `db:"method"`
	RequestMethod sql.NullString `db:"request_method"`
	DeviceType    sql.NullString `db:"device_type"`
	OperName      sql.NullString `db:"oper_name"`
	OperUrl       sql.NullString `db:"oper_url"`
	OperIp        sql.NullString `db:"oper_ip"`
	OperLocation  sql.NullString `db:"oper_location"`
	OperParam     sql.NullString `db:"oper_param"`
	JsonResult    sql.NullString `db:"json_result"`
	Status        sql.NullString `db:"status"`
	ErrorMsg      sql.NullString `db:"error_msg"`
	OperTime      sql.NullTime   `db:"oper_time"`
	CostTime      sql.NullInt64  `db:"cost_time"`
	UserAgent     sql.NullString `db:"user_agent"`
}

func getOperLogByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*operLogRow, error) {
	var row operLogRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `select id, title, business_type, method, request_method, device_type, oper_name, oper_url, oper_ip, oper_location, oper_param, json_result, status, error_msg, oper_time, cost_time, user_agent from public.s_oper_log where id = $1 limit 1`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("操作日志不存在")
		}
		return nil, err
	}
	return &row, nil
}

func operLogIn(ids []int64, start int) (string, []interface{}) {
	parts := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		parts = append(parts, fmt.Sprintf("$%d", start+i))
		args = append(args, id)
	}
	return strings.Join(parts, ", "), args
}

func toOperLogPB(row operLogRow) *pb.OperLog {
	return &pb.OperLog{
		Id:            row.Id,
		Title:         nullString(row.Title),
		BusinessType:  nullString(row.BusinessType),
		Method:        nullString(row.Method),
		RequestMethod: nullString(row.RequestMethod),
		DeviceType:    nullString(row.DeviceType),
		OperName:      nullString(row.OperName),
		OperUrl:       nullString(row.OperUrl),
		OperIp:        nullString(row.OperIp),
		OperLocation:  nullString(row.OperLocation),
		OperParam:     nullString(row.OperParam),
		JsonResult:    nullString(row.JsonResult),
		Status:        nullString(row.Status),
		ErrorMsg:      nullString(row.ErrorMsg),
		OperTime:      nullTime(row.OperTime),
		CostTime:      nullInt64(row.CostTime),
		UserAgent:     nullString(row.UserAgent),
	}
}

func toOperLogList(rows []operLogRow) []*pb.OperLog {
	list := make([]*pb.OperLog, 0, len(rows))
	for _, row := range rows {
		list = append(list, toOperLogPB(row))
	}
	return list
}
