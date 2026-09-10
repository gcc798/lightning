package data

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
)

func dictEntityToItem(item *entpkg.DictData) *v1.DictItem {
	if item == nil {
		return nil
	}
	return &v1.DictItem{
		Id:          item.ID,
		ParentId:    item.ParentID,
		DictType:    item.DictType,
		DictLabel:   item.DictLabel,
		DictValue:   item.DictValue,
		IsDefault:   item.IsDefault,
		Status:      item.Status,
		Sort:        item.Sort,
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func (r *Resources) activeDicts(ctx context.Context) ([]*entpkg.DictData, error) {
	items, err := r.Ent.DictData.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.DictData, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *Resources) CreateDict(ctx context.Context, req *v1.DictItem) error {
	items, err := r.activeDicts(ctx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.DictType == req.GetDictType() && item.DictValue == req.GetDictValue() {
			return errors.New("dict value already exists")
		}
		if req.GetParentId() > 0 && item.ID == req.GetParentId() && item.DictType != req.GetDictType() {
			return errors.New("parent dict type mismatch")
		}
	}
	if req.GetParentId() > 0 {
		parentFound := false
		for _, item := range items {
			if item.ID == req.GetParentId() {
				parentFound = true
				break
			}
		}
		if !parentFound {
			return errors.New("parent dict not found")
		}
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	_, err = r.Ent.DictData.Create().SetID(nextID()).SetParentID(req.GetParentId()).SetDictType(req.GetDictType()).SetDictLabel(req.GetDictLabel()).SetDictValue(req.GetDictValue()).SetIsDefault(req.GetIsDefault()).SetStatus(req.GetStatus()).SetSort(req.GetSort()).SetRemark(req.GetRemark()).SetCreateBy(operator).SetUpdateBy(operator).SetCreatedTime(now).SetUpdatedTime(now).Save(ctx)
	return err
}

func (r *Resources) PageDicts(ctx context.Context, req *v1.PageDictRequest) (*v1.PageDictReply, error) {
	items, err := r.activeDicts(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.DictItem, 0, len(items))
	for _, item := range items {
		if req.GetDictType() != "" && item.DictType != req.GetDictType() {
			continue
		}
		if req.GetDictLabel() != "" && !strings.Contains(item.DictLabel, req.GetDictLabel()) {
			continue
		}
		if req.Status != nil && req.GetStatus() >= 0 && item.Status != req.GetStatus() {
			continue
		}
		filtered = append(filtered, dictEntityToItem(item))
	}
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageDictReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) GetDict(ctx context.Context, id int64) (*v1.DictItem, error) {
	item, err := r.Ent.DictData.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return dictEntityToItem(item), nil
}

func (r *Resources) FindDictsByType(ctx context.Context, dictType string, parentID *int64) ([]*v1.DictItem, error) {
	items, err := r.activeDicts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.DictItem, 0)
	for _, item := range items {
		if item.DictType == dictType {
			if parentID != nil && item.ParentID != *parentID {
				continue
			}
			out = append(out, dictEntityToItem(item))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetSort() < out[j].GetSort() })
	return out, nil
}

func (r *Resources) UpdateDict(ctx context.Context, req *v1.UpdateDictRequest) error {
	if req.GetParentId() == req.GetId() {
		return errors.New("dict can not be its own parent")
	}
	current, err := r.Ent.DictData.Get(ctx, req.GetId())
	if err != nil {
		return err
	}
	if current.DeletedAt != nil {
		return errors.New("dict not found")
	}
	items, err := r.activeDicts(ctx)
	if err != nil {
		return err
	}
	parentFound := req.GetParentId() == 0
	for _, item := range items {
		if item.ID != req.GetId() && item.DictType == req.GetDictType() && item.DictValue == req.GetDictValue() {
			return errors.New("dict value already exists")
		}
		if req.GetParentId() > 0 && item.ID == req.GetParentId() {
			parentFound = true
			if item.DictType != req.GetDictType() {
				return errors.New("parent dict type mismatch")
			}
		}
	}
	if !parentFound {
		return errors.New("parent dict not found")
	}
	_, err = r.Ent.DictData.UpdateOneID(req.GetId()).SetParentID(req.GetParentId()).SetDictType(req.GetDictType()).SetDictLabel(req.GetDictLabel()).SetDictValue(req.GetDictValue()).SetIsDefault(req.GetIsDefault()).SetStatus(req.GetStatus()).SetSort(req.GetSort()).SetRemark(req.GetRemark()).SetUpdateBy(currentOperatorID(ctx)).SetUpdatedTime(time.Now()).Save(ctx)
	return err
}

func (r *Resources) DeleteDicts(ctx context.Context, ids ...int64) error {
	items, err := r.activeDicts(ctx)
	if err != nil {
		return err
	}
	children := make(map[int64][]int64)
	for _, item := range items {
		children[item.ParentID] = append(children[item.ParentID], item.ID)
	}
	deleteIDs := make(map[int64]struct{})
	var walk func(int64)
	walk = func(id int64) {
		deleteIDs[id] = struct{}{}
		for _, childID := range children[id] {
			walk(childID)
		}
	}
	for _, id := range ids {
		found := false
		for _, item := range items {
			if item.ID == id {
				found = true
				break
			}
		}
		if !found {
			return errors.New("dict not found")
		}
		walk(id)
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	for id := range deleteIDs {
		if _, err := r.Ent.DictData.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
	}
	return nil
}
