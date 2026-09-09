package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
)

type userRow struct {
	Id          int64          `db:"id"`
	OrgId       sql.NullInt64  `db:"org_id"`
	UserName    string         `db:"user_name"`
	NickName    sql.NullString `db:"nick_name"`
	UserType    int64          `db:"user_type"`
	Email       sql.NullString `db:"email"`
	Phonenumber sql.NullString `db:"phonenumber"`
	Sex         int64          `db:"sex"`
	Avatar      sql.NullString `db:"avatar"`
	Status      int64          `db:"status"`
	Sort        int64          `db:"sort"`
	LoginIp     sql.NullString `db:"login_ip"`
	LoginDate   sql.NullInt64  `db:"login_date"`
	OpenId      sql.NullString `db:"open_id"`
	UnionId     sql.NullString `db:"union_id"`
	Remark      sql.NullString `db:"remark"`
	CreateBy    sql.NullInt64  `db:"create_by"`
	UpdateBy    sql.NullInt64  `db:"update_by"`
	CreatedTime sql.NullTime   `db:"created_time"`
	UpdatedTime sql.NullTime   `db:"updated_time"`
}

func listUsersByRole(ctx context.Context, svcCtx *svc.ServiceContext, roleID int64) ([]userRow, error) {
	var rows []userRow
	err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select u.id, u.user_name, u.nick_name, u.user_type, u.email, u.phonenumber, u.sex, u.avatar, u.status, u.sort,
		       u.login_ip, u.login_date, u.open_id, u.union_id, u.remark, u.create_by, u.update_by, u.created_time, u.updated_time, u.org_id
		from public.s_user u
		inner join public.m_user_role ur on u.id = ur.user_id
		where ur.role_id = $1
		order by u.created_time desc
	`, roleID)
	if err != nil {
		return nil, fmt.Errorf("查询角色用户失败: %w", err)
	}
	return rows, nil
}

func toUserPB(row userRow) *pb.User {
	return &pb.User{
		UserId:      row.Id,
		UserName:    row.UserName,
		NickName:    nullString(row.NickName),
		UserType:    int32(row.UserType),
		Email:       nullString(row.Email),
		Phonenumber: nullString(row.Phonenumber),
		Sex:         int32(row.Sex),
		Avatar:      nullString(row.Avatar),
		Status:      int32(row.Status),
		Sort:        row.Sort,
		LoginIp:     nullString(row.LoginIp),
		LoginDate:   nullInt64(row.LoginDate),
		OpenId:      nullString(row.OpenId),
		UnionId:     nullString(row.UnionId),
		Remark:      nullString(row.Remark),
		CreateBy:    nullInt64(row.CreateBy),
		UpdateBy:    nullInt64(row.UpdateBy),
		CreatedAt:   nullTime(row.CreatedTime),
		UpdatedAt:   nullTime(row.UpdatedTime),
		OrgId:       nullInt64(row.OrgId),
	}
}

func toUserListPB(rows []userRow) []*pb.User {
	list := make([]*pb.User, 0, len(rows))
	for _, row := range rows {
		list = append(list, toUserPB(row))
	}
	return list
}
