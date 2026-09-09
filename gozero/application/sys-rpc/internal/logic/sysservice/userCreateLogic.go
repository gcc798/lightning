package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCreateLogic) UserCreate(in *pb.UserCreateReq) (*pb.Ack, error) {
	if in.UserName == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_user where user_name = $1`, in.UserName); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("用户名已存在")
	}
	if in.Phonenumber != "" {
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_user where phonenumber = $1`, in.Phonenumber); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("手机号已存在")
		}
	}
	if in.Email != "" {
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_user where email = $1`, in.Email); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("邮箱已存在")
		}
	}
	password := ""
	if in.Password != "" {
		bs, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(bs)
	}
	_, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		insert into public.s_user (user_name, nick_name, password, user_type, email, phonenumber, sex, avatar, status, remark, create_by, update_by, created_time, updated_time)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, nullif($11, 0), nullif($12, 0), now(), now())
	`,
		in.UserName,
		sql.NullString{String: in.NickName, Valid: in.NickName != ""},
		sql.NullString{String: password, Valid: password != ""},
		in.UserType,
		sql.NullString{String: in.Email, Valid: in.Email != ""},
		sql.NullString{String: in.Phonenumber, Valid: in.Phonenumber != ""},
		in.Sex,
		sql.NullString{String: in.Avatar, Valid: in.Avatar != ""},
		in.Status,
		sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		in.CreateBy,
		in.UpdateBy,
	)
	if err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
