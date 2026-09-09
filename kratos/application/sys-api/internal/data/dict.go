package data

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
)

type DictRepo struct {
	client v1.DictServiceClient
}

func NewDictRepo(clients *RPCClientSet) *DictRepo {
	return &DictRepo{client: clients.Dict}
}
func (r *DictRepo) Create(ctx context.Context, item *v1.DictItem) (*v1.MessageReply, error) {
	return r.client.CreateDict(ctx, item)
}
func (r *DictRepo) Page(ctx context.Context, req *v1.PageDictRequest) (*v1.PageDictReply, error) {
	return r.client.PageDict(ctx, req)
}
func (r *DictRepo) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	return r.client.BatchDeleteDict(ctx, &v1.DictBatchIdsRequest{Ids: ids})
}
func (r *DictRepo) ByType(ctx context.Context, t string, parentID *int64) (*v1.DictListReply, error) {
	return r.client.GetDictByType(ctx, &v1.GetDictByTypeRequest{DictType: t, ParentId: parentID})
}
func (r *DictRepo) Label(ctx context.Context, t, v string) (*v1.DictLabelReply, error) {
	return r.client.GetDictLabel(ctx, &v1.GetDictLabelRequest{DictType: t, DictValue: v})
}
func (r *DictRepo) Update(ctx context.Context, req *v1.UpdateDictRequest) (*v1.MessageReply, error) {
	return r.client.UpdateDict(ctx, req)
}
func (r *DictRepo) Get(ctx context.Context, id int64) (*v1.DictItem, error) {
	return r.client.GetDictById(ctx, &v1.DictIdRequest{Id: id})
}
func (r *DictRepo) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.DeleteDict(ctx, &v1.DictIdRequest{Id: id})
}
