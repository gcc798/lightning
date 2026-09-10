package data

import (
	"context"
	"errors"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
)

func configEntityToItem(item *entpkg.SystemConfig) *v1.ConfigItem {
	if item == nil {
		return nil
	}
	return &v1.ConfigItem{
		Id:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		Data:        rawJSONToProtoValue(item.Data),
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		Remark:      item.Remark,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func (r *Resources) activeConfigs(ctx context.Context) ([]*entpkg.SystemConfig, error) {
	items, err := r.Ent.SystemConfig.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.SystemConfig, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *Resources) CreateConfig(ctx context.Context, req *v1.ConfigItem) error {
	items, err := r.activeConfigs(ctx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.Name == req.GetName() {
			return errors.New("config name already exists")
		}
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	_, err = r.Ent.SystemConfig.Create().SetID(nextID()).SetName(req.GetName()).SetCode(req.GetCode()).SetData(protoValueToRawJSON(req.GetData())).SetRemark(req.GetRemark()).SetCreateBy(operator).SetUpdateBy(operator).SetCreatedTime(now).SetUpdatedTime(now).Save(ctx)
	return err
}

func (r *Resources) PageConfigs(ctx context.Context, req *v1.PageConfigRequest) (*v1.PageConfigReply, error) {
	items, err := r.activeConfigs(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.ConfigItem, 0, len(items))
	for _, item := range items {
		if req.GetName() != "" && !strings.Contains(item.Name, req.GetName()) {
			continue
		}
		if req.GetCode() != "" && item.Code != req.GetCode() {
			continue
		}
		filtered = append(filtered, configEntityToItem(item))
	}
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageConfigReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) GetConfig(ctx context.Context, id int64) (*v1.ConfigItem, error) {
	item, err := r.Ent.SystemConfig.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return configEntityToItem(item), nil
}

func (r *Resources) FindConfigsByCode(ctx context.Context, code string) ([]*v1.ConfigItem, error) {
	items, err := r.activeConfigs(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.ConfigItem, 0)
	for _, item := range items {
		if item.Code == code {
			out = append(out, configEntityToItem(item))
		}
	}
	return out, nil
}

func (r *Resources) UpdateConfig(ctx context.Context, req *v1.UpdateConfigRequest) error {
	current, err := r.Ent.SystemConfig.Get(ctx, req.GetId())
	if err != nil {
		return err
	}
	if current.DeletedAt != nil {
		return errors.New("config not found")
	}
	items, err := r.activeConfigs(ctx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ID != req.GetId() && item.Name == req.GetName() {
			return errors.New("config name already exists")
		}
	}
	_, err = r.Ent.SystemConfig.UpdateOneID(req.GetId()).SetName(req.GetName()).SetCode(req.GetCode()).SetData(protoValueToRawJSON(req.GetData())).SetRemark(req.GetRemark()).SetUpdateBy(currentOperatorID(ctx)).SetUpdatedTime(time.Now()).Save(ctx)
	return err
}

func (r *Resources) DeleteConfigs(ctx context.Context, ids ...int64) error {
	now := time.Now()
	operator := currentOperatorID(ctx)
	for _, id := range ids {
		item, err := r.Ent.SystemConfig.Get(ctx, id)
		if err != nil {
			return err
		}
		if item.DeletedAt != nil {
			return errors.New("config not found")
		}
		if _, err := r.Ent.SystemConfig.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
	}
	return nil
}
