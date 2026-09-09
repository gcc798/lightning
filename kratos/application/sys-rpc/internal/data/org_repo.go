package data

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	entpkg "github.com/gcc798/lightning/kratos/application/sys-rpc/ent"
)

func orgEntityToItem(item *entpkg.Org) *v1.OrgItem {
	if item == nil {
		return nil
	}
	return &v1.OrgItem{
		OrgId:       item.ID,
		ParentId:    item.ParentID,
		OrgName:     item.OrgName,
		OrgCode:     item.OrgCode,
		OrgType:     item.OrgType,
		Leader:      item.Leader,
		Phone:       item.Phone,
		Email:       item.Email,
		Status:      item.Status,
		Sort:        item.Sort,
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func (r *Resources) activeOrgs(ctx context.Context) ([]*entpkg.Org, error) {
	items, err := r.Ent.Org.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.Org, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func buildOrgAncestors(items []*entpkg.Org, parentID int64) (string, error) {
	if parentID == 0 {
		return "0", nil
	}
	for _, item := range items {
		if item.ID == parentID {
			if item.Ancestors == "" {
				return "0", nil
			}
			return item.Ancestors + "," + strconv.FormatInt(parentID, 10), nil
		}
	}
	return "", errors.New("parent org not found")
}

func (r *Resources) CreateOrg(ctx context.Context, req *v1.CreateOrgRequest) (int64, error) {
	items, err := r.activeOrgs(ctx)
	if err != nil {
		return 0, err
	}
	for _, item := range items {
		if req.GetOrgCode() != "" && item.OrgCode == req.GetOrgCode() {
			return 0, errors.New("org code already exists")
		}
	}
	ancestors, err := buildOrgAncestors(items, req.GetParentId())
	if err != nil {
		return 0, err
	}
	now := time.Now()
	orgType := req.GetOrgType()
	if orgType == "" {
		orgType = "company"
	}
	operator := currentOperatorID(ctx)
	item, err := r.Ent.Org.Create().
		SetID(nextID()).
		SetParentID(req.GetParentId()).
		SetAncestors(ancestors).
		SetOrgName(req.GetOrgName()).
		SetOrgCode(req.GetOrgCode()).
		SetOrgType(orgType).
		SetLeader(req.GetLeader()).
		SetPhone(req.GetPhone()).
		SetEmail(req.GetEmail()).
		SetStatus(req.GetStatus()).
		SetSort(req.GetSort()).
		SetRemark(req.GetRemark()).
		SetCreateBy(operator).
		SetUpdateBy(operator).
		SetCreatedTime(now).
		SetUpdatedTime(now).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *Resources) UpdateOrg(ctx context.Context, req *v1.UpdateOrgRequest) error {
	current, err := r.Ent.Org.Get(ctx, req.GetOrgId())
	if err != nil {
		return err
	}
	if current.DeletedAt != nil {
		return errors.New("org not found")
	}
	if req.GetParentId() == req.GetOrgId() {
		return errors.New("org can not be its own parent")
	}
	items, err := r.activeOrgs(ctx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ID != req.GetOrgId() && req.GetOrgCode() != "" && item.OrgCode == req.GetOrgCode() {
			return errors.New("org code already exists")
		}
	}
	ancestors, err := buildOrgAncestors(items, req.GetParentId())
	if err != nil {
		return err
	}
	if strings.Contains(","+ancestors+",", ","+strconv.FormatInt(req.GetOrgId(), 10)+",") {
		return errors.New("org can not be moved under its child")
	}
	orgType := req.GetOrgType()
	if orgType == "" {
		orgType = current.OrgType
	}
	operator := currentOperatorID(ctx)
	_, err = r.Ent.Org.UpdateOneID(req.GetOrgId()).
		SetParentID(req.GetParentId()).
		SetAncestors(ancestors).
		SetOrgName(req.GetOrgName()).
		SetOrgCode(req.GetOrgCode()).
		SetOrgType(orgType).
		SetLeader(req.GetLeader()).
		SetPhone(req.GetPhone()).
		SetEmail(req.GetEmail()).
		SetStatus(req.GetStatus()).
		SetSort(req.GetSort()).
		SetRemark(req.GetRemark()).
		SetUpdateBy(operator).
		SetUpdatedTime(time.Now()).
		Save(ctx)
	return err
}

func (r *Resources) GetOrg(ctx context.Context, id int64) (*v1.OrgItem, error) {
	item, err := r.Ent.Org.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return orgEntityToItem(item), nil
}

func (r *Resources) ListOrgs(ctx context.Context) ([]*v1.OrgItem, error) {
	items, err := r.activeOrgs(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.OrgItem, 0, len(items))
	for _, item := range items {
		out = append(out, orgEntityToItem(item))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].GetSort() == out[j].GetSort() {
			return out[i].GetOrgId() < out[j].GetOrgId()
		}
		return out[i].GetSort() < out[j].GetSort()
	})
	return out, nil
}

func (r *Resources) DeleteOrgs(ctx context.Context, ids ...int64) error {
	orgs, err := r.activeOrgs(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	for _, id := range ids {
		exists := false
		for _, item := range orgs {
			if item.ID == id {
				exists = true
				break
			}
		}
		if !exists {
			return errors.New("org not found")
		}
		for _, item := range orgs {
			if item.ParentID == id {
				return errors.New("org has children")
			}
		}
		if _, err := r.Ent.Org.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
	}
	return nil
}
