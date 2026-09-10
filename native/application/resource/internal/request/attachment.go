package request

import "github.com/gcc798/microservice-kit/internal/utils"

// BindAttachmentToBusinessRequest 绑定附件到业务请求（步骤2：绑定业务信息）
type BindAttachmentToBusinessRequest struct {
	BusinessType  string           `json:"businessType" binding:"required" msg:"业务类型不能为空"` // 业务类型
	BusinessId    string           `json:"businessId" binding:"required" msg:"业务ID不能为空"`   // 业务ID
	BusinessField string           `json:"businessField"`                                  // 业务字段
	IsPublic      bool             `json:"isPublic"`                                       // 是否公开
	Metadata      map[string]any   `json:"metadata"`                                       // 元数据（JSON对象）
	ExpireTime    *utils.LocalTime `json:"expireTime"`                                     // 过期时间
}

// GetAttachmentURLRequest 获取附件 URL 请求（Query 参数）
type GetAttachmentURLRequest struct {
	Expires int `form:"expires" query:"expires" binding:"omitempty,min=0" msg:"过期时间不能为负数"` // 过期时间（秒），0 表示永久，默认 3600
}

// ListAttachmentsByBusinessRequest 根据业务查询附件列表请求（Query 参数）
type ListAttachmentsByBusinessRequest struct {
	BusinessType string `form:"businessType" query:"businessType" binding:"required" msg:"业务类型不能为空"`
	BusinessId   string `form:"businessId" query:"businessId" binding:"required" msg:"业务ID不能为空"`
}

// PageAttachmentsRequest 分页查询附件列表请求
type PageAttachmentsRequest struct {
	PageNum      int    `json:"pageNum" binding:"required,min=1" msg:"页码必须大于等于1"`
	PageSize     int    `json:"pageSize" binding:"required,min=1,max=100" msg:"每页数量必须是1-100之间"`
	FileName     string `json:"fileName"`
	FileType     string `json:"fileType"`
	BusinessType string `json:"businessType"`
}
