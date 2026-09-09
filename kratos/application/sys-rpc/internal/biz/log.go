package biz

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type LogUsecase struct {
	res  *data.Resources
	kind string
}

type LoginLogUsecase struct{ *LogUsecase }
type OperLogUsecase struct{ *LogUsecase }

func NewLogUsecase(res *data.Resources, kind string) *LogUsecase {
	return &LogUsecase{res: res, kind: kind}
}

func NewLoginLogUsecase(res *data.Resources) *LoginLogUsecase {
	return &LoginLogUsecase{LogUsecase: NewLogUsecase(res, "login")}
}

func NewOperLogUsecase(res *data.Resources) *OperLogUsecase {
	return &OperLogUsecase{LogUsecase: NewLogUsecase(res, "oper")}
}

func (uc *LogUsecase) CreateLogin(ctx context.Context, req *v1.CreateLoginLogRequest) (*v1.MessageReply, error) {
	if err := uc.res.CreateLoginLog(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *LogUsecase) CreateOper(ctx context.Context, req *v1.CreateOperLogRequest) (*v1.MessageReply, error) {
	if err := uc.res.CreateOperLog(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *LogUsecase) PageLogin(ctx context.Context, req *v1.PageLoginLogRequest) (*v1.PageLogReply, error) {
	return uc.res.PageLoginLogs(ctx, req)
}

func (uc *LogUsecase) PageOper(ctx context.Context, req *v1.PageOperLogRequest) (*v1.PageLogReply, error) {
	return uc.res.PageOperLogs(ctx, req)
}

func (uc *LogUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteLogs(ctx, uc.kind, ids...); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *LogUsecase) Clean(ctx context.Context, days int32) (*v1.LogCleanReply, error) {
	count, err := uc.res.CleanLogs(ctx, uc.kind, int(days))
	if err != nil {
		return nil, err
	}
	message := "清理登录日志成功"
	if uc.kind == "oper" {
		message = "清理操作日志成功"
	}
	return &v1.LogCleanReply{
		Message: message,
		Count:   count,
		Days:    days,
	}, nil
}

func (uc *LogUsecase) UpdateLogin(ctx context.Context, req *v1.UpdateLoginLogRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateLoginLog(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *LogUsecase) UpdateOper(ctx context.Context, req *v1.UpdateOperLogRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateOperLog(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *LogUsecase) Get(ctx context.Context, id int64) (*v1.LogItem, error) {
	item, err := uc.res.GetLog(ctx, uc.kind, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "日志不存在")
}

func (uc *LogUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteLogs(ctx, uc.kind, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}
