package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserResetPasswordLogic {
	return &UserResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserResetPasswordLogic) UserResetPassword(in *pb.UserPasswordReq) (*pb.Ack, error) {
	if in.Id <= 0 || in.NewPassword == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	if _, err := getUserByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_user set password = $2, updated_time = now() where id = $1`, in.Id, sql.NullString{String: string(password), Valid: true}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
