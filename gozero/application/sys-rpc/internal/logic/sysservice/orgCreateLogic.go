package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgCreateLogic {
	return &OrgCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgCreateLogic) OrgCreate(in *pb.OrgCreateReq) (*pb.Ack, error) {
	if in.OrgName == "" {
		return nil, errors.New("组织名称不能为空")
	}
	if in.OrgCode != "" {
		var count int64
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_org where org_code = $1`, in.OrgCode); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("组织编码已存在")
		}
	}
	ancestors, err := buildOrgAncestors(l.ctx, l.svcCtx, in.ParentId)
	if err != nil {
		return nil, err
	}
	orgType := in.OrgType
	if orgType == "" {
		orgType = "company"
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_org (parent_id, ancestors, org_name, org_code, org_type, leader, phone, email, status, sort, remark, create_by, update_by, created_time, updated_time)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, null, null, now(), now())`,
		in.ParentId, ancestors, in.OrgName, sql.NullString{String: in.OrgCode, Valid: in.OrgCode != ""}, orgType,
		sql.NullString{String: in.Leader, Valid: in.Leader != ""},
		sql.NullString{String: in.Phone, Valid: in.Phone != ""},
		sql.NullString{String: in.Email, Valid: in.Email != ""},
		in.Status, in.Sort, sql.NullString{String: in.Remark, Valid: in.Remark != ""}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
