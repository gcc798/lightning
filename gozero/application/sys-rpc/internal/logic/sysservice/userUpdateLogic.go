package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateLogic) UserUpdate(in *pb.UserUpdateReq) (*pb.Ack, error) {
	if in.Id <= 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if _, err := getUserByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	var count int64
	if in.UserName != "" {
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_user where id <> $1 and user_name = $2`, in.Id, in.UserName); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("用户名已被占用")
		}
	}
	if in.Phonenumber != "" {
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_user where id <> $1 and phonenumber = $2`, in.Id, in.Phonenumber); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("手机号已被占用")
		}
	}
	if in.Email != "" {
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_user where id <> $1 and email = $2`, in.Id, in.Email); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("邮箱已被占用")
		}
	}
	_, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		update public.s_user
		set user_name = $2, nick_name = $3, user_type = $4, email = $5, phonenumber = $6, sex = $7, avatar = $8, status = $9, remark = $10, update_by = nullif($11, 0), updated_time = now()
		where id = $1
	`,
		in.Id, in.UserName,
		sql.NullString{String: in.NickName, Valid: in.NickName != ""},
		in.UserType,
		sql.NullString{String: in.Email, Valid: in.Email != ""},
		sql.NullString{String: in.Phonenumber, Valid: in.Phonenumber != ""},
		in.Sex,
		sql.NullString{String: in.Avatar, Valid: in.Avatar != ""},
		in.Status,
		sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		in.UpdateBy,
	)
	if err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
