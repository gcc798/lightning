package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDeleteLogic {
	return &DictDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDeleteLogic) DictDelete(in *pb.IdReq) (*pb.Ack, error) {
	if _, err := getDictByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	_, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		with recursive dict_tree as (
			select id from public.s_dict_data where id = $1
			union all
			select d.id from public.s_dict_data d inner join dict_tree dt on d.parent_id = dt.id
		)
		delete from public.s_dict_data where id in (select id from dict_tree)
	`, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
