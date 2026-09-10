package sysservicelogic

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictBatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictBatchDeleteLogic {
	return &DictBatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictBatchDeleteLogic) DictBatchDelete(in *pb.BatchIdsReq) (*pb.Ack, error) {
	if len(in.Ids) == 0 {
		return nil, errors.New("ids 不能为空")
	}
	placeholders, args := buildInt64In(in.Ids, 1)
	query := fmt.Sprintf(`
		with recursive dict_tree as (
			select id from public.s_dict_data where id in (%s)
			union all
			select d.id from public.s_dict_data d inner join dict_tree dt on d.parent_id = dt.id
		)
		delete from public.s_dict_data where id in (select id from dict_tree)
	`, placeholders)
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, query, args...); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
