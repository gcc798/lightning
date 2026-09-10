package response

import (
	"github.com/gcc798/microservice-kit/application/sys/internal/domain/model"
	"github.com/gcc798/microservice-kit/internal/utils"
)

// OperLogResponse 操作日志响应
type OperLogResponse struct {
	ID            int64           `json:"id"`            // 日志ID
	Title         string          `json:"title"`         // 模块标题
	BusinessType  string          `json:"businessType"`  // 业务类型
	Method        string          `json:"method"`        // 调用方法
	RequestMethod string          `json:"requestMethod"` // 请求方式：GET/POST
	DeviceType    string          `json:"deviceType"`    // 终端类型：web/ios/android/wechat
	OperName      string          `json:"operName"`      // 操作者
	OperUrl       string          `json:"operUrl"`       // 请求URL
	OperIp        string          `json:"operIp"`        // 操作IP
	OperLocation  string          `json:"operLocation"`  // 操作地点
	OperParam     string          `json:"operParam"`     // 请求参数
	JsonResult    string          `json:"jsonResult"`    // 返回结果
	Status        string          `json:"status"`        // 操作状态：0成功 1失败
	ErrorMsg      string          `json:"errorMsg"`      // 错误信息
	OperTime      utils.LocalTime `json:"operTime"`      // 操作时间
	CostTime      int64           `json:"costTime"`      // 耗时（毫秒）
	UserAgent     string          `json:"userAgent"`     // UA
}

// ToOperLogResponse 转换为操作日志响应
func ToOperLogResponse(log *model.OperLog) OperLogResponse {
	return OperLogResponse{
		ID:            log.ID,
		Title:         log.Title,
		BusinessType:  log.BusinessType,
		Method:        log.Method,
		RequestMethod: log.RequestMethod,
		DeviceType:    log.DeviceType,
		OperName:      log.OperName,
		OperUrl:       log.OperUrl,
		OperIp:        log.OperIp,
		OperLocation:  log.OperLocation,
		OperParam:     log.OperParam,
		JsonResult:    log.JsonResult,
		Status:        log.Status,
		ErrorMsg:      log.ErrorMsg,
		OperTime:      log.OperTime,
		CostTime:      log.CostTime,
		UserAgent:     log.UserAgent,
	}
}
