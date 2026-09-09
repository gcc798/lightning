package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgUpdateLogic {
	return &OrgUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgUpdateLogic) OrgUpdate(in *pb.OrgUpdateReq) (*pb.Ack, error) {
	row, err := getOrgByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.ParentId == in.Id && in.Id != 0 {
		return nil, errors.New("不能将组织设置为自己的子组织")
	}
	if in.OrgCode != "" && in.OrgCode != nullString(row.OrgCode) {
		var count int64
		if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.s_org where org_code = $1 and id <> $2`, in.OrgCode, in.Id); err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("组织编码已被占用")
		}
	}
	parentID := row.ParentId
	ancestors := nullString(row.Ancestors)
	if in.ParentId != row.ParentId {
		nextAncestors, err := buildOrgAncestors(l.ctx, l.svcCtx, in.ParentId)
		if err != nil {
			return nil, err
		}
		if strings.Contains(","+nextAncestors+",", ","+strconv.FormatInt(in.Id, 10)+",") {
			return nil, errors.New("不能将组织移动到其子组织下")
		}
		parentID = in.ParentId
		ancestors = nextAncestors
	}
	orgName := row.OrgName
	if in.OrgName != "" {
		orgName = in.OrgName
	}
	orgCode := row.OrgCode
	if in.OrgCode != "" {
		orgCode = sql.NullString{String: in.OrgCode, Valid: true}
	}
	orgType := row.OrgType
	if in.OrgType != "" {
		orgType = in.OrgType
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_org
		set parent_id = $2, ancestors = $3, org_name = $4, org_code = $5, org_type = $6, leader = $7, phone = $8, email = $9, status = $10, sort = $11, remark = $12, updated_time = now()
		where id = $1`,
		in.Id, parentID, ancestors, orgName, orgCode, orgType,
		sql.NullString{String: in.Leader, Valid: in.Leader != ""},
		sql.NullString{String: in.Phone, Valid: in.Phone != ""},
		sql.NullString{String: in.Email, Valid: in.Email != ""},
		in.Status, in.Sort, sql.NullString{String: in.Remark, Valid: in.Remark != ""}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
