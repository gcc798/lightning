package sys

import "github.com/gcc798/lightning/internal/utils/pagination"

// CreateDictRequest 创建字典请求
type CreateDictRequest struct {
	ParentId  int64  `json:"parentId"`                                    // 父字典ID（0表示根节点）
	DictType  string `json:"dictType" binding:"required" msg:"字典类型不能为空"`  // 字典类型
	DictLabel string `json:"dictLabel" binding:"required" msg:"字典标签不能为空"` // 字典标签
	DictValue string `json:"dictValue" binding:"required" msg:"字典键值不能为空"` // 字典键值
	Sort      int64  `json:"sort"`                                        // 排序
	IsDefault bool   `json:"isDefault"`                                   // 是否默认
	Status    int32  `json:"status"`                                      // 状态：0正常 1停用
	Remark    string `json:"remark"`                                      // 备注
	CreateBy  int64  `json:"createBy"`                                    // 创建者
	UpdateBy  int64  `json:"updateBy"`                                    // 更新者
}

// UpdateDictRequest 更新字典请求
type UpdateDictRequest struct {
	ID        int64  `json:"id"`                                          // 字典ID（由路径参数注入）
	ParentId  int64  `json:"parentId"`                                    // 父字典ID
	DictType  string `json:"dictType" binding:"required" msg:"字典类型不能为空"`  // 字典类型
	DictLabel string `json:"dictLabel" binding:"required" msg:"字典标签不能为空"` // 字典标签
	DictValue string `json:"dictValue" binding:"required" msg:"字典键值不能为空"` // 字典键值
	Sort      int64  `json:"sort"`                                        // 排序
	IsDefault bool   `json:"isDefault"`                                   // 是否默认
	Status    int32  `json:"status"`                                      // 状态
	Remark    string `json:"remark"`                                      // 备注
	UpdateBy  int64  `json:"updateBy"`                                    // 更新者
}

// PageDictRequest 查询字典列表请求
type PageDictRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	DictType             string `json:"dictType"`  // 字典类型
	DictLabel            string `json:"dictLabel"` // 字典标签（模糊查询）
	Status               int32  `json:"status"`    // 状态：0正常 1停用 -1全部
}

// GetDictByTypeRequest 根据类型获取字典请求
type GetDictByTypeRequest struct {
	DictType string `form:"dictType" query:"dictType" binding:"required" msg:"字典类型不能为空"` // 字典类型
	ParentId *int64 `form:"parentId" query:"parentId"`                                   // 父字典ID（可选，用于获取子字典）
}

// BatchDeleteDictRequest 批量删除字典请求
type BatchDeleteDictRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1" msg:"请至少选择一个字典"` // 字典ID列表
}
