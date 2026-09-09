package biz

import (
	"context"
	"strings"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
)

type OrgUsecase struct {
	res *data.Resources
}

func NewOrgUsecase(res *data.Resources) *OrgUsecase { return &OrgUsecase{res: res} }

func (uc *OrgUsecase) Create(ctx context.Context, req *v1.CreateOrgRequest) (*v1.OrgIdReply, error) {
	id, err := uc.res.CreateOrg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.OrgIdReply{Id: id}, nil
}

func (uc *OrgUsecase) Update(ctx context.Context, req *v1.UpdateOrgRequest) (*v1.MessageReply, error) {
	if err := uc.res.UpdateOrg(ctx, req); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *OrgUsecase) Delete(ctx context.Context, id int64) (*v1.MessageReply, error) {
	if err := uc.res.DeleteOrgs(ctx, id); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *OrgUsecase) BatchDelete(ctx context.Context, ids []int64) (*v1.MessageReply, error) {
	for _, id := range ids {
		if err := uc.res.DeleteOrgs(ctx, id); err != nil {
			return nil, err
		}
	}
	return okReply(), nil
}

func (uc *OrgUsecase) GetByID(ctx context.Context, id int64) (*v1.OrgItem, error) {
	item, err := uc.res.GetOrg(ctx, id)
	if err != nil {
		return nil, err
	}
	return requireFound(item, "组织不存在")
}

func (uc *OrgUsecase) Tree(ctx context.Context) (*v1.OrgTreeReply, error) {
	items, err := uc.res.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.OrgTreeReply{Items: buildOrgTree(items)}, nil
}

func (uc *OrgUsecase) Page(ctx context.Context, req *v1.PageOrgsRequest) (*v1.PageOrgsReply, error) {
	items, err := uc.res.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.OrgItem, 0)
	for _, item := range items {
		if req.GetOrgName() != "" && !strings.Contains(item.GetOrgName(), req.GetOrgName()) {
			continue
		}
		if req.GetOrgCode() != "" && !strings.Contains(item.GetOrgCode(), req.GetOrgCode()) {
			continue
		}
		if req.Status != nil && item.GetStatus() != req.GetStatus() {
			continue
		}
		if req.ParentId != nil && item.GetParentId() != req.GetParentId() {
			continue
		}
		filtered = append(filtered, item)
	}
	pageNum, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	paged, total := paginateOrgs(filtered, pageNum, pageSize)
	return &v1.PageOrgsReply{List: paged, Total: total, PageNum: pageNum, PageSize: pageSize}, nil
}

func normalizePage(pageNum, pageSize int64) (int64, int64) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return pageNum, pageSize
}

func paginateOrgs(items []*v1.OrgItem, pageNum, pageSize int64) ([]*v1.OrgItem, int64) {
	total := int64(len(items))
	start := (pageNum - 1) * pageSize
	if start >= total {
		return []*v1.OrgItem{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}

func buildOrgTree(items []*v1.OrgItem) []*v1.OrgItem {
	byParent := make(map[int64][]*v1.OrgItem)
	for _, item := range items {
		item.Children = nil
		byParent[item.GetParentId()] = append(byParent[item.GetParentId()], item)
	}
	var build func(parentID int64) []*v1.OrgItem
	build = func(parentID int64) []*v1.OrgItem {
		nodes := byParent[parentID]
		for _, node := range nodes {
			node.Children = build(node.GetOrgId())
		}
		return nodes
	}
	return build(0)
}
