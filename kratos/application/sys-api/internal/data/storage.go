package data

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
)

type StorageEnvRepo struct {
	client v1.StorageEnvServiceClient
}

func NewStorageEnvRepo(clients *RPCClientSet) *StorageEnvRepo {
	return &StorageEnvRepo{client: clients.StorageEnv}
}
func (r *StorageEnvRepo) Create(ctx context.Context, item *v1.StorageEnvItem) (*v1.StorageEnvItem, error) {
	return r.client.CreateStorageEnv(ctx, item)
}
func (r *StorageEnvRepo) Page(ctx context.Context, req *v1.PageStorageEnvRequest) (*v1.PageStorageEnvReply, error) {
	return r.client.PageStorageEnv(ctx, req)
}
func (r *StorageEnvRepo) Default(ctx context.Context) (*v1.StorageEnvItem, error) {
	return r.client.GetDefaultStorageEnv(ctx, &v1.StorageEmpty{})
}
func (r *StorageEnvRepo) SetDefault(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.SetDefaultStorageEnv(ctx, &v1.SetDefaultStorageEnvRequest{Id: id})
}
func (r *StorageEnvRepo) Update(ctx context.Context, req *v1.UpdateStorageEnvRequest) (*v1.MessageReply, error) {
	return r.client.UpdateStorageEnv(ctx, req)
}
func (r *StorageEnvRepo) Get(ctx context.Context, id int64) (*v1.StorageEnvItem, error) {
	return r.client.GetStorageEnv(ctx, &v1.StorageIdRequest{Id: id})
}
func (r *StorageEnvRepo) Test(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.TestStorageEnvConnection(ctx, &v1.StorageIdRequest{Id: id})
}
func (r *StorageEnvRepo) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	return r.client.DeleteStorageEnv(ctx, &v1.StorageIdRequest{Id: id})
}
