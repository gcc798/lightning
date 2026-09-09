package menu

import "github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

type nativeMenuTree struct {
	Id          int64             `json:"id"`
	MenuName    string            `json:"menuName"`
	ParentId    int64             `json:"parentId"`
	Sort        int64             `json:"sort"`
	Path        string            `json:"path"`
	Component   string            `json:"component"`
	Query       string            `json:"query"`
	IsFrame     int64             `json:"isFrame"`
	IsCache     int64             `json:"isCache"`
	MenuType    int64             `json:"menuType"`
	Visible     int64             `json:"visible"`
	Status      int64             `json:"status"`
	Perms       string            `json:"perms"`
	Icon        string            `json:"icon"`
	Remark      string            `json:"remark"`
	CreateBy    int64             `json:"createBy"`
	UpdateBy    int64             `json:"updateBy"`
	CreatedTime *string           `json:"createdTime"`
	UpdatedTime *string           `json:"updatedTime"`
	Children    []*nativeMenuTree `json:"children,omitempty"`
}

func toNativeMenuTreeList(records []*sysservice.Menu) []*nativeMenuTree {
	out := make([]*nativeMenuTree, 0, len(records))
	for _, item := range records {
		out = append(out, toNativeMenuTree(item))
	}
	return out
}

func toNativeMenuTree(item *sysservice.Menu) *nativeMenuTree {
	if item == nil {
		return nil
	}
	node := &nativeMenuTree{
		Id:          item.Id,
		MenuName:    item.MenuName,
		ParentId:    item.ParentId,
		Sort:        item.Sort,
		Path:        item.Path,
		Component:   item.Component,
		Query:       item.Query,
		IsFrame:     item.IsFrame,
		IsCache:     item.IsCache,
		MenuType:    item.MenuType,
		Visible:     item.Visible,
		Status:      item.Status,
		Perms:       item.Perms,
		Icon:        item.Icon,
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		CreatedTime: optionalTime(item.CreatedTime),
		UpdatedTime: optionalTime(item.UpdatedTime),
	}
	if len(item.Children) > 0 {
		node.Children = toNativeMenuTreeList(item.Children)
	}
	return node
}

func optionalTime(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
