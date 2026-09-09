package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserImportLogic {
	return &UserImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserImportLogic) UserImport(in *pb.UserImportReq) (*pb.Ack, error) {
	successCount := 0
	failCount := 0
	for _, user := range in.Users {
		if _, err := NewUserCreateLogic(l.ctx, l.svcCtx).UserCreate(user); err != nil {
			failCount++
			l.Errorf("导入用户失败: %s, 错误: %v", user.UserName, err)
		} else {
			successCount++
		}
	}
	return &pb.Ack{Msg: fmt.Sprintf("导入完成：成功%d条，失败%d条", successCount, failCount)}, nil
}
