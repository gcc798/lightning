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

type UserPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPageLogic {
	return &UserPageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *UserPageLogic) UserPage(in *pb.UserPageReq) (*pb.UserPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.Username != "" {
		args = append(args, "%"+in.Username+"%")
		where = append(where, fmt.Sprintf("user_name like $%d", len(args)))
	}
	if in.Phonenumber != "" {
		args = append(args, "%"+in.Phonenumber+"%")
		where = append(where, fmt.Sprintf("phonenumber like $%d", len(args)))
	}
	if in.Status == 0 || in.Status == 1 {
		args = append(args, in.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_user where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []struct {
		Id          int64          `db:"id"`
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
		OrgId       sql.NullInt64  `db:"org_id"`
	}
	query := `select id, user_name, nick_name, user_type, email, phonenumber, sex, avatar, status, sort, login_ip, login_date, open_id, union_id, remark, create_by, update_by, created_time, updated_time, org_id
		from public.s_user where ` + whereSQL + ` order by sort asc, created_time desc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	records := make([]*pb.User, 0, len(rows))
	for _, row := range rows {
		records = append(records, &pb.User{
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
		})
	}
	return &pb.UserPageResp{Records: records, Page: toPageInfo(total, pageNum, pageSize)}, nil
}
