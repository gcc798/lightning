package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvPageLogic {
	return &StorageEnvPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvPageLogic) StorageEnvPage(in *pb.StorageEnvPageReq) (*pb.StorageEnvPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.Name != "" {
		args = append(args, "%"+in.Name+"%")
		where = append(where, fmt.Sprintf("name like $%d", len(args)))
	}
	if in.StorageType != "" {
		args = append(args, in.StorageType)
		where = append(where, fmt.Sprintf("storage_type = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_storage_env where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []storageEnvRow
	query := `select id, name, code, storage_type, is_default, status, config, remark, create_by, created_time, update_by, updated_time from public.s_storage_env where ` + whereSQL + ` order by is_default desc, id desc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.StorageEnvPageResp{Records: toStorageEnvList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
