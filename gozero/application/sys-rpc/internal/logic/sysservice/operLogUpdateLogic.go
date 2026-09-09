package sysservicelogic

import (
	"context"
	"database/sql"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogUpdateLogic {
	return &OperLogUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogUpdateLogic) OperLogUpdate(in *pb.OperLogUpdateReq) (*pb.Ack, error) {
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_oper_log set title = $2, business_type = $3, method = $4, request_method = $5, device_type = $6, oper_name = $7, oper_url = $8, oper_ip = $9, oper_location = $10, oper_param = $11, json_result = $12, status = $13, error_msg = $14, cost_time = $15, user_agent = $16 where id = $1`,
		in.Id,
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
		sql.NullInt64{Int64: in.CostTime, Valid: in.CostTime > 0},
		sql.NullString{String: in.UserAgent, Valid: in.UserAgent != ""},
	); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
