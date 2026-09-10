package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangePasswordLogic {
	return &UserChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserChangePasswordLogic) UserChangePassword(in *pb.UserChangePasswordReq) (*pb.Ack, error) {
	if in.UserId <= 0 || in.OldPassword == "" || in.NewPassword == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	var row struct {
		Password string `db:"password"`
	}
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &row, `select password from public.s_user where id = $1 limit 1`, in.UserId); err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(in.OldPassword)); err != nil {
		return nil, fmt.Errorf("旧密码错误")
	}
	password, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_user set password = $2, updated_time = now() where id = $1`, in.UserId, string(password)); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
