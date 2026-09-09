package response

import (
	"github.com/gcc798/lightning/application/sys/internal/domain/model"
	"github.com/gcc798/lightning/internal/utils"
)

// LoginLogResponse 登录日志响应
type LoginLogResponse struct {
	ID            int64           `json:"id"`            // 日志ID
	UserName      string          `json:"userName"`      // 用户名
	Ipaddr        string          `json:"ipaddr"`        // 登录IP
	LoginLocation string          `json:"loginLocation"` // 登录地点
	Browser       string          `json:"browser"`       // 浏览器类型
	Os            string          `json:"os"`            // 操作系统
	Status        int32           `json:"status"`        // 登录状态：0成功 1失败
	Msg           string          `json:"msg"`           // 提示消息
	LoginTime     utils.LocalTime `json:"loginTime"`     // 登录时间
	ClientId      string          `json:"clientId"`      // 客户端ID
}

// ToLoginLogResponse 转换为登录日志响应
func ToLoginLogResponse(log *model.LoginLog) LoginLogResponse {
	return LoginLogResponse{
		ID:            log.ID,
		UserName:      log.UserName,
		Ipaddr:        log.Ipaddr,
		LoginLocation: log.LoginLocation,
		Browser:       log.Browser,
		Os:            log.Os,
		Status:        log.Status,
		Msg:           log.Msg,
		LoginTime:     log.LoginTime,
		ClientId:      log.ClientId,
	}
}
