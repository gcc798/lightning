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

type loginLogRow struct {
	Id            int64          `db:"id"`
	UserName      sql.NullString `db:"user_name"`
	Ipaddr        sql.NullString `db:"ipaddr"`
	LoginLocation sql.NullString `db:"login_location"`
	Browser       sql.NullString `db:"browser"`
	Os            sql.NullString `db:"os"`
	Status        int64          `db:"status"`
	Msg           sql.NullString `db:"msg"`
	LoginTime     sql.NullTime   `db:"login_time"`
	ClientId      sql.NullString `db:"client_id"`
}

func getLoginLogByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*loginLogRow, error) {
	var row loginLogRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `select id, user_name, ipaddr, login_location, browser, os, status, msg, login_time, client_id from public.s_login_log where id = $1 limit 1`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("登录日志不存在")
		}
		return nil, err
	}
	return &row, nil
}

func loginLogIn(ids []int64, start int) (string, []interface{}) {
	parts := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		parts = append(parts, fmt.Sprintf("$%d", start+i))
		args = append(args, id)
	}
	return strings.Join(parts, ", "), args
}

func toLoginLogPB(row loginLogRow) *pb.LoginLog {
	return &pb.LoginLog{
		Id:            row.Id,
		UserName:      nullString(row.UserName),
		Ipaddr:        nullString(row.Ipaddr),
		LoginLocation: nullString(row.LoginLocation),
		Browser:       nullString(row.Browser),
		Os:            nullString(row.Os),
		Status:        int32(row.Status),
		Msg:           nullString(row.Msg),
		LoginTime:     nullTime(row.LoginTime),
		ClientId:      nullString(row.ClientId),
	}
}

func toLoginLogList(rows []loginLogRow) []*pb.LoginLog {
	list := make([]*pb.LoginLog, 0, len(rows))
	for _, row := range rows {
		list = append(list, toLoginLogPB(row))
	}
	return list
}
