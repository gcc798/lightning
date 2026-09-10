package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictCreateLogic {
	return &DictCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictCreateLogic) DictCreate(in *pb.DictCreateReq) (*pb.Ack, error) {
	if in.DictType == "" || in.DictLabel == "" || in.DictValue == "" {
		return nil, errors.New("字典类型、标签和值不能为空")
	}
	exists, err := dictValueExists(l.ctx, l.svcCtx, in.DictType, in.DictValue, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("字典值已存在")
	}
	if in.ParentId > 0 {
		parent, err := getDictByID(l.ctx, l.svcCtx, in.ParentId)
		if err != nil {
			return nil, errors.New("父字典不存在")
		}
		if parent.DictType.Valid && parent.DictType.String != in.DictType {
			return nil, errors.New("父字典类型不匹配")
		}
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		insert into public.s_dict_data (parent_id, dict_type, dict_label, dict_value, sort, is_default, status, remark, create_by, update_by, created_time, updated_time)
		values ($1, $2, $3, $4, $5, $6, $7, $8, nullif($9, 0), nullif($10, 0), now(), now())
	`, in.ParentId, in.DictType, in.DictLabel, in.DictValue, in.Sort, in.IsDefault, in.Status, sql.NullString{String: in.Remark, Valid: in.Remark != ""}, in.CreateBy, in.UpdateBy); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
