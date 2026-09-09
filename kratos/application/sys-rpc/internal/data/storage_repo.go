package data

import (
	"context"
	"errors"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	entpkg "github.com/gcc798/lightning/kratos/application/sys-rpc/ent"
)

func storageEntityToItem(item *entpkg.StorageEnv) *v1.StorageEnvItem {
	if item == nil {
		return nil
	}
	return &v1.StorageEnvItem{
		Id:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		StorageType: item.StorageType,
		IsDefault:   item.IsDefault,
		Status:      item.Status,
		Config:      mapToProtoStruct(item.Config),
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func setDefaultStorageEnvTx(ctx context.Context, tx *entpkg.Tx, id int64) error {
	items, err := tx.StorageEnv.Query().All(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	for _, item := range items {
		if item.DeletedAt != nil {
			continue
		}
		if _, err := tx.StorageEnv.UpdateOneID(item.ID).SetIsDefault(item.ID == id).SetUpdatedTime(now).Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resources) CreateStorageEnv(ctx context.Context, req *v1.StorageEnvItem) (*v1.StorageEnvItem, error) {
	now := time.Now()
	createBy := operatorID(ctx, req.GetCreateBy())
	updateBy := operatorID(ctx, req.GetUpdateBy())
	var created *v1.StorageEnvItem
	err := r.withTx(ctx, func(tx *entpkg.Tx) error {
		items, err := tx.StorageEnv.Query().All(ctx)
		if err != nil {
			return err
		}
		hasDefault := false
		for _, item := range items {
			if item.DeletedAt != nil {
				continue
			}
			if strings.EqualFold(item.Code, req.GetCode()) {
				return errors.New("storage env code already exists")
			}
			if item.IsDefault {
				hasDefault = true
			}
		}
		isDefault := req.GetIsDefault()
		if !hasDefault {
			isDefault = true
		}
		item, err := tx.StorageEnv.Create().
			SetID(nextID()).
			SetName(req.GetName()).
			SetCode(req.GetCode()).
			SetStorageType(req.GetStorageType()).
			SetIsDefault(isDefault).
			SetStatus(req.GetStatus()).
			SetConfig(protoStructToMap(req.GetConfig())).
			SetRemark(req.GetRemark()).
			SetCreateBy(createBy).
			SetUpdateBy(updateBy).
			SetCreatedTime(now).
			SetUpdatedTime(now).
			Save(ctx)
		if err != nil {
			return err
		}
		created = storageEntityToItem(item)
		if isDefault {
			return setDefaultStorageEnvTx(ctx, tx, item.ID)
		}
		return nil
	})
	if err == nil && r.Storage != nil {
		r.Storage.InvalidateAll()
	}
	return created, err
}

func (r *Resources) PageStorageEnvs(ctx context.Context, req *v1.PageStorageEnvRequest) (*v1.PageStorageEnvReply, error) {
	items, err := r.Ent.StorageEnv.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.StorageEnvItem, 0)
	for _, item := range items {
		if item.DeletedAt != nil {
			continue
		}
		if req.GetName() != "" && !strings.Contains(item.Name, req.GetName()) {
			continue
		}
		if req.GetStorageType() != "" && item.StorageType != req.GetStorageType() {
			continue
		}
		filtered = append(filtered, storageEntityToItem(item))
	}
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageStorageEnvReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) GetStorageEnv(ctx context.Context, id int64) (*v1.StorageEnvItem, error) {
	item, err := r.Ent.StorageEnv.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return storageEntityToItem(item), nil
}

func (r *Resources) GetDefaultStorageEnv(ctx context.Context) (*v1.StorageEnvItem, error) {
	if r.Storage != nil {
		env, err := r.Storage.ResolveActiveEnv(ctx, "")
		if err == nil && env != nil {
			return storageEntityToItem(env), nil
		}
	}
	return nil, nil
}

func (r *Resources) SetDefaultStorageEnv(ctx context.Context, id int64) error {
	err := r.withTx(ctx, func(tx *entpkg.Tx) error {
		item, err := tx.StorageEnv.Get(ctx, id)
		if err != nil {
			return err
		}
		if item.DeletedAt != nil {
			return errors.New("storage env not found")
		}
		return setDefaultStorageEnvTx(ctx, tx, id)
	})
	if err == nil && r.Storage != nil {
		r.Storage.InvalidateAll()
	}
	return err
}

func (r *Resources) UpdateStorageEnv(ctx context.Context, req *v1.UpdateStorageEnvRequest) error {
	err := r.withTx(ctx, func(tx *entpkg.Tx) error {
		current, err := tx.StorageEnv.Get(ctx, req.GetId())
		if err != nil {
			return err
		}
		if current.DeletedAt != nil {
			return errors.New("storage env not found")
		}
		items, err := tx.StorageEnv.Query().All(ctx)
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.DeletedAt != nil || item.ID == req.GetId() {
				continue
			}
			if strings.EqualFold(item.Code, req.GetCode()) {
				return errors.New("storage env code already exists")
			}
		}
		if _, err := tx.StorageEnv.UpdateOneID(req.GetId()).
			SetName(req.GetName()).
			SetCode(req.GetCode()).
			SetStorageType(req.GetStorageType()).
			SetIsDefault(req.GetIsDefault()).
			SetStatus(req.GetStatus()).
			SetConfig(protoStructToMap(req.GetConfig())).
			SetRemark(req.GetRemark()).
			SetUpdateBy(operatorID(ctx, req.GetUpdateBy())).
			SetUpdatedTime(time.Now()).
			Save(ctx); err != nil {
			return err
		}
		if req.GetIsDefault() {
			return setDefaultStorageEnvTx(ctx, tx, req.GetId())
		}
		return nil
	})
	if err == nil && r.Storage != nil {
		r.Storage.InvalidateAll()
	}
	return err
}

func (r *Resources) TestStorageEnvConnection(ctx context.Context, id int64) error {
	item, err := r.GetStorageEnv(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("storage env not found")
	}
	if item.GetStatus() != 0 {
		return errors.New("storage env disabled")
	}
	if r.Storage == nil {
		return errors.New("storage manager is not initialized")
	}
	return r.Storage.TestConnection(ctx, id)
}

func (r *Resources) DeleteStorageEnvs(ctx context.Context, ids ...int64) error {
	now := time.Now()
	operator := currentOperatorID(ctx)
	for _, id := range ids {
		item, err := r.Ent.StorageEnv.Get(ctx, id)
		if err != nil {
			return err
		}
		if item.DeletedAt != nil {
			return errors.New("storage env not found")
		}
		if item.IsDefault {
			return errors.New("default storage env can not be deleted")
		}
		attachments, err := r.Ent.Attachment.Query().All(ctx)
		if err != nil {
			return err
		}
		for _, attachment := range attachments {
			if attachment.DeletedAt == nil && attachment.EnvID == id {
				return errors.New("storage env has attachments")
			}
		}
		if _, err := r.Ent.StorageEnv.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
		if r.Storage != nil {
			r.Storage.InvalidateEnv(id)
		}
	}
	return nil
}
