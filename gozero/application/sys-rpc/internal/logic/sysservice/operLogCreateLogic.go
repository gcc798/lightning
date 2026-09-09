package sysservicelogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCreateLogic {
	return &OperLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogCreateLogic) OperLogCreate(in *pb.OperLogReq) (*pb.Ack, error) {
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_oper_log (title, business_type, method, request_method, device_type, oper_name, oper_url, oper_ip, oper_location, oper_param, json_result, status, error_msg, oper_time, cost_time, user_agent) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		sql.NullString{String: in.Title, Valid: in.Title != ""},
		sql.NullString{String: in.BusinessType, Valid: in.BusinessType != ""},
		sql.NullString{String: in.Method, Valid: in.Method != ""},
		sql.NullString{String: in.RequestMethod, Valid: in.RequestMethod != ""},
		sql.NullString{String: in.DeviceType, Valid: in.DeviceType != ""},
		sql.NullString{String: in.OperName, Valid: in.OperName != ""},
		sql.NullString{String: in.OperUrl, Valid: in.OperUrl != ""},
		sql.NullString{String: in.OperIp, Valid: in.OperIp != ""},
		sql.NullString{String: in.OperLocation, Valid: in.OperLocation != ""},
		sql.NullString{String: in.OperParam, Valid: in.OperParam != ""},
		sql.NullString{String: in.JsonResult, Valid: in.JsonResult != ""},
		sql.NullString{String: in.Status, Valid: in.Status != ""},
		sql.NullString{String: in.ErrorMsg, Valid: in.ErrorMsg != ""},
		time.Now(),
		sql.NullInt64{Int64: in.CostTime, Valid: in.CostTime > 0},
		sql.NullString{String: in.UserAgent, Valid: in.UserAgent != ""},
	); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
