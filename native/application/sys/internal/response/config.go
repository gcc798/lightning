package response

import (
	"encoding/json"

	"github.com/gcc798/microservice-kit/application/sys/internal/domain/model"
	"github.com/gcc798/microservice-kit/internal/utils"
)

// ConfigResponse 配置响应
type ConfigResponse struct {
	ID          int64           `json:"id"`          // 配置ID
	Name        string          `json:"name"`        // 配置名称
	Code        string          `json:"code"`        // 配置编码
	Data        json.RawMessage `json:"data"`        // 配置数据（JSON格式）
	Remark      string          `json:"remark"`      // 备注
	CreateBy    int64           `json:"createBy"`    // 创建者
	CreatedTime utils.LocalTime `json:"createdTime"` // 创建时间
	UpdateBy    int64           `json:"updateBy"`    // 更新者
	UpdatedTime utils.LocalTime `json:"updatedTime"` // 更新时间
}

// ConfigDataResponse 配置数据响应（仅返回data字段）
type ConfigDataResponse struct {
	Code string          `json:"code"` // 配置编码
	Data json.RawMessage `json:"data"` // 配置数据
}

// ToConfigResponse 转换为配置响应
func ToConfigResponse(config *model.Config) ConfigResponse {
	return ConfigResponse{
		ID:          config.ID,
		Name:        config.Name,
		Code:        config.Code,
		Data:        config.Data,
		Remark:      config.Remark,
		CreateBy:    config.CreateBy,
		CreatedTime: config.CreatedTime,
		UpdateBy:    config.UpdateBy,
		UpdatedTime: config.UpdatedTime,
	}
}
