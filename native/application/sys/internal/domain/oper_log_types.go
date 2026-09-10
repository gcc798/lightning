package sys

import "github.com/gcc798/microservice-kit/internal/utils/pagination"

// CreateOperLogRequest 创建操作日志请求
type CreateOperLogRequest struct {
	Title         string `json:"title" binding:"required" msg:"模块标题不能为空"` // 模块标题
	BusinessType  string `json:"businessType"`                            // 业务类型
	Method        string `json:"method"`                                  // 调用方法
	RequestMethod string `json:"requestMethod"`                           // 请求方式：GET/POST
	DeviceType    string `json:"deviceType"`                              // 终端类型：web/ios/android/wechat
	OperName      string `json:"operName"`                                // 操作者
	OperUrl       string `json:"operUrl"`                                 // 请求URL
	OperIp        string `json:"operIp"`                                  // 操作IP
	OperLocation  string `json:"operLocation"`                            // 操作地点
	OperParam     string `json:"operParam"`                               // 请求参数
	JsonResult    string `json:"jsonResult"`                              // 返回结果
	Status        string `json:"status"`                                  // 操作状态：0成功 1失败
	ErrorMsg      string `json:"errorMsg"`                                // 错误信息
	CostTime      int64  `json:"costTime"`                                // 耗时（毫秒）
	UserAgent     string `json:"userAgent"`                               // UA
}

// UpdateOperLogRequest 更新操作日志请求
type UpdateOperLogRequest struct {
	ID            int64  `json:"id" binding:"required" msg:"日志ID不能为空"` // 日志ID
	Title         string `json:"title"`                                // 模块标题
	BusinessType  string `json:"businessType"`                         // 业务类型
	Method        string `json:"method"`                               // 调用方法
	RequestMethod string `json:"requestMethod"`                        // 请求方式：GET/POST
	DeviceType    string `json:"deviceType"`                           // 终端类型：web/ios/android/wechat
	OperName      string `json:"operName"`                             // 操作者
	OperUrl       string `json:"operUrl"`                              // 请求URL
	OperIp        string `json:"operIp"`                               // 操作IP
	OperLocation  string `json:"operLocation"`                         // 操作地点
	OperParam     string `json:"operParam"`                            // 请求参数
	JsonResult    string `json:"jsonResult"`                           // 返回结果
	Status        string `json:"status"`                               // 操作状态：0成功 1失败
	ErrorMsg      string `json:"errorMsg"`                             // 错误信息
	CostTime      int64  `json:"costTime"`                             // 耗时（毫秒）
	UserAgent     string `json:"userAgent"`                            // UA
}

// BatchDeleteOperLogRequest 批量删除操作日志请求
type BatchDeleteOperLogRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1" msg:"请至少选择一个日志"` // 日志ID列表
}

// PageOperLogRequest 操作日志列表查询请求
type PageOperLogRequest struct {
	pagination.PageQuery         // 嵌入分页参数
	Title                string  `json:"title"`        // 模块标题（可选,模糊查询）
	OperName             string  `json:"operName"`     // 操作者（可选,模糊查询）
	BusinessType         string  `json:"businessType"` // 业务类型（可选）
	Status               *string `json:"status"`       // 操作状态（可选,nil表示全部,"0"成功 "1"失败）
	StartTime            string  `json:"startTime"`    // 开始时间（可选）
	EndTime              string  `json:"endTime"`      // 结束时间（可选）
}

// CleanOperLogRequest 清理操作日志请求
type CleanOperLogRequest struct {
	Days int `json:"days" binding:"required,min=1" msg:"清理天数必须大于等于1天"` // 清理多少天之前的日志
}
