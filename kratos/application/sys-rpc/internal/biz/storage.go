package biz

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type StorageEnvUsecase struct{ res *data.Resources }

func NewStorageEnvUsecase(res *data.Resources) *StorageEnvUsecase {
	return &StorageEnvUsecase{res: res}
}

func (uc *StorageEnvUsecase) Create(ctx context.Context, req *v1.StorageEnvItem) (*v1.StorageEnvItem, error) {
	return uc.res.CreateStorageEnv(ctx, req)
}

func (uc *StorageEnvUsecase) Page(ctx context.Context, req *v1.PageStorageEnvRequest) (*v1.PageStorageEnvReply, error) {
	return uc.res.PageStorageEnvs(ctx, req)
}

func (uc *StorageEnvUsecase) Default(ctx context.Context) (*v1.StorageEnvItem, error) {
	item, err := uc.res.GetDefaultStorageEnv(ctx)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return &v1.StorageEnvItem{}, nil
	}
	return item, nil
}

func (uc *StorageEnvUsecase) SetDefault(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.SetDefaultStorageEnv(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *StorageEnvUsecase) Update(ctx context.Context, req *v1.UpdateStorageEnvRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateStorageEnv(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *StorageEnvUsecase) Get(ctx context.Context, id int64) (*v1.StorageEnvItem, error) {
	item, err := uc.res.GetStorageEnv(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "存储环境不存在")
}

func (uc *StorageEnvUsecase) Test(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.TestStorageEnvConnection(ctx, id); err != nil {
		if err.Error() == "storage env not found" {
			return nil, notFound("存储环境不存在")
		}
		return nil, badRequest(err.Error())
	}
	return okReply(), nil
}

func (uc *StorageEnvUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteStorageEnvs(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}
