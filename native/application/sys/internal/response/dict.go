package response

import (
	"github.com/gcc798/microservice-kit/application/sys/internal/domain/model"
	"github.com/gcc798/microservice-kit/internal/utils"
)

// DictDataResponse 字典数据响应
type DictDataResponse struct {
	ID          int64           `json:"id"`          // 字典ID
	ParentId    int64           `json:"parentId"`    // 父字典ID
	DictType    string          `json:"dictType"`    // 字典类型
	DictLabel   string          `json:"dictLabel"`   // 字典标签
	DictValue   string          `json:"dictValue"`   // 字典键值
	Sort        int64           `json:"sort"`        // 排序
	IsDefault   bool            `json:"isDefault"`   // 是否默认
	Status      int32           `json:"status"`      // 状态
	Remark      string          `json:"remark"`      // 备注
	CreateBy    int64           `json:"createBy"`    // 创建者
	CreatedTime utils.LocalTime `json:"createdTime"` // 创建时间
	UpdateBy    int64           `json:"updateBy"`    // 更新者
	UpdatedTime utils.LocalTime `json:"updatedTime"` // 更新时间
}

// ToDictDataResponse 转换为字典响应
func ToDictDataResponse(dict *model.DictData) DictDataResponse {
	return DictDataResponse{
		ID:          dict.ID,
		ParentId:    dict.ParentId,
		DictType:    dict.DictType,
		DictLabel:   dict.DictLabel,
		DictValue:   dict.DictValue,
		Sort:        dict.Sort,
		IsDefault:   dict.IsDefault,
		Status:      dict.Status,
		Remark:      dict.Remark,
		CreateBy:    dict.CreateBy,
		CreatedTime: dict.CreatedTime,
		UpdateBy:    dict.UpdateBy,
		UpdatedTime: dict.UpdatedTime,
	}
}
