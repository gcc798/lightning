package biz

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type DictUsecase struct{ res *data.Resources }

func NewDictUsecase(res *data.Resources) *DictUsecase { return &DictUsecase{res: res} }

func (uc *DictUsecase) Create(ctx context.Context, req *v1.DictItem) (*v1.MessageReply, error) {
	if err := uc.res.CreateDict(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *DictUsecase) Page(ctx context.Context, req *v1.PageDictRequest) (*v1.PageDictReply, error) {
	return uc.res.PageDicts(ctx, req)
}

func (uc *DictUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteDicts(ctx, ids...); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *DictUsecase) ByType(ctx context.Context, dictType string, parentID *int64) (*v1.DictListReply, error) {
	items, err := uc.res.FindDictsByType(ctx, dictType, parentID)
	if err != nil {
		return nil, err
	}
	return &v1.DictListReply{Items: items}, nil
}

func (uc *DictUsecase) Label(ctx context.Context, dictType, dictValue string) (*v1.DictLabelReply, error) {
	items, err := uc.res.FindDictsByType(ctx, dictType, nil)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.GetDictValue() == dictValue {
			return &v1.DictLabelReply{Label: item.GetDictLabel()}, nil
		}
	}
	return &v1.DictLabelReply{}, nil
}

func (uc *DictUsecase) Update(ctx context.Context, req *v1.UpdateDictRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateDict(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *DictUsecase) Get(ctx context.Context, id int64) (*v1.DictItem, error) {
	item, err := uc.res.GetDict(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "字典不存在")
}

func (uc *DictUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteDicts(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}
