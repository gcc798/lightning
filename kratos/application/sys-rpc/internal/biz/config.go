package biz

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type ConfigUsecase struct{ res *data.Resources }

func NewConfigUsecase(res *data.Resources) *ConfigUsecase { return &ConfigUsecase{res: res} }

func (uc *ConfigUsecase) Create(ctx context.Context, req *v1.ConfigItem) (*v1.MessageReply, error) {
	if err := uc.res.CreateConfig(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *ConfigUsecase) Page(ctx context.Context, req *v1.PageConfigRequest) (*v1.PageConfigReply, error) {
	return uc.res.PageConfigs(ctx, req)
}

func (uc *ConfigUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteConfigs(ctx, ids...); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *ConfigUsecase) ByCode(ctx context.Context, code string) (*v1.ConfigListReply, error) {
	items, err := uc.res.FindConfigsByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &v1.ConfigListReply{Items: items}, nil
}

func (uc *ConfigUsecase) Data(ctx context.Context, code string) (*v1.ConfigDataReply, error) {
	items, err := uc.res.FindConfigsByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return &v1.ConfigDataReply{}, nil
	}
	return &v1.ConfigDataReply{Data: items[0].GetData()}, nil
}

func (uc *ConfigUsecase) Update(ctx context.Context, req *v1.UpdateConfigRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateConfig(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *ConfigUsecase) Get(ctx context.Context, id int64) (*v1.ConfigItem, error) {
	item, err := uc.res.GetConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "参数配置不存在")
}

func (uc *ConfigUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteConfigs(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}
