package sys

import "github.com/gcc798/microservice-kit/internal/utils/pagination"

// CreateLoginLogRequest 创建登录日志请求
type CreateLoginLogRequest struct {
	UserName      string `json:"userName" binding:"required" msg:"用户名不能为空"` // 用户名
	Ipaddr        string `json:"ipaddr"`                                    // 登录IP
	LoginLocation string `json:"loginLocation"`                             // 登录地点
	Browser       string `json:"browser"`                                   // 浏览器类型
	Os            string `json:"os"`                                        // 操作系统
	Status        int32  `json:"status"`                                    // 登录状态：0成功 1失败
	Msg           string `json:"msg"`                                       // 提示消息
	ClientId      string `json:"clientId"`                                  // 客户端ID
}

// UpdateLoginLogRequest 更新登录日志请求
type UpdateLoginLogRequest struct {
	ID            int64  `json:"id" binding:"required" msg:"日志ID不能为空"` // 日志ID
	UserName      string `json:"userName"`                             // 用户名
	Ipaddr        string `json:"ipaddr"`                               // 登录IP
	LoginLocation string `json:"loginLocation"`                        // 登录地点
	Browser       string `json:"browser"`                              // 浏览器类型
	Os            string `json:"os"`                                   // 操作系统
	Status        int32  `json:"status"`                               // 登录状态：0成功 1失败
	Msg           string `json:"msg"`                                  // 提示消息
	ClientId      string `json:"clientId"`                             // 客户端ID
}

// BatchDeleteLoginLogRequest 批量删除登录日志请求
type BatchDeleteLoginLogRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1" msg:"请至少选择一个日志"` // 日志ID列表
}

// PageLoginLogRequest 登录日志列表查询请求
type PageLoginLogRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	UserName             string `json:"userName"`  // 用户名（可选,模糊查询）
	Ipaddr               string `json:"ipaddr"`    // 登录IP（可选,模糊查询）
	Status               *int32 `json:"status"`    // 登录状态（可选,nil表示全部,0成功 1失败）
	StartTime            string `json:"startTime"` // 开始时间（可选）
	EndTime              string `json:"endTime"`   // 结束时间（可选）
}

// CleanLoginLogRequest 清理登录日志请求
type CleanLoginLogRequest struct {
	Days int `json:"days" binding:"required,min=1" msg:"清理天数必须大于等于1天"` // 清理多少天之前的日志
}
