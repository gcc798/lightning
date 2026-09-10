package data

import (
	"context"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type LoginLogRepo struct {
	client v1.LoginLogServiceClient
}

type OperLogRepo struct {
	client v1.OperLogServiceClient
}

func NewLoginLogRepo(clients *RPCClientSet) *LoginLogRepo {
	return &LoginLogRepo{client: clients.LoginLog}
}

func NewOperLogRepo(clients *RPCClientSet) *OperLogRepo {
	return &OperLogRepo{client: clients.OperLog}
}

func (r *LoginLogRepo) Create(ctx context.Context, item *v1.CreateLoginLogRequest) (*v1.MessageReply, error) {
	return r.client.CreateLoginLog(ctx, item)
}

func (r *LoginLogRepo) Page(ctx context.Context, req *v1.PageLoginLogRequest) (*v1.PageLogReply, error) {
	return r.client.PageLoginLog(ctx, req)
}

func (r *LoginLogRepo) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return r.client.BatchDeleteLoginLog(ctx, &v1.LogBatchIdsRequest{Ids: ids})
}

func (r *LoginLogRepo) Clean(ctx context.Context, days int32) (*v1.LogCleanReply, error) {
	return r.client.CleanLoginLog(ctx, &v1.CleanLogRequest{Days: days})
}

func (r *LoginLogRepo) Update(ctx context.Context, req *v1.UpdateLoginLogRequest) (*v1.MessageReply, error) {
	return r.client.UpdateLoginLog(ctx, req)
}

func (r *LoginLogRepo) Get(ctx context.Context, id int64) (*v1.LogItem, error) {
	return r.client.GetLoginLogById(ctx, &v1.LogIdRequest{Id: id})
}

func (r *LoginLogRepo) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.DeleteLoginLog(ctx, &v1.LogIdRequest{Id: id})
}

func (r *OperLogRepo) Create(ctx context.Context, item *v1.CreateOperLogRequest) (*v1.MessageReply, error) {
	return r.client.CreateOperLog(ctx, item)
}

func (r *OperLogRepo) Page(ctx context.Context, req *v1.PageOperLogRequest) (*v1.PageLogReply, error) {
	return r.client.PageOperLog(ctx, req)
}

func (r *OperLogRepo) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return r.client.BatchDeleteOperLog(ctx, &v1.LogBatchIdsRequest{Ids: ids})
}

func (r *OperLogRepo) Clean(ctx context.Context, days int32) (*v1.LogCleanReply, error) {
	return r.client.CleanOperLog(ctx, &v1.CleanLogRequest{Days: days})
}

func (r *OperLogRepo) Update(ctx context.Context, req *v1.UpdateOperLogRequest) (*v1.MessageReply, error) {
	return r.client.UpdateOperLog(ctx, req)
}

func (r *OperLogRepo) Get(ctx context.Context, id int64) (*v1.LogItem, error) {
	return r.client.GetOperLogById(ctx, &v1.LogIdRequest{Id: id})
}

func (r *OperLogRepo) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.DeleteOperLog(ctx, &v1.LogIdRequest{Id: id})
}
