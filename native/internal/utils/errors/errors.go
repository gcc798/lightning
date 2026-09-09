package errors

import (
	"fmt"
	"net/http"
)

// ErrorType 错误类型
type ErrorType string

const (
	// ErrorTypeBusiness 业务逻辑错误（可向用户展示）
	ErrorTypeBusiness ErrorType = "BUSINESS"
	// ErrorTypeInfrastructure 基础设施错误（需要记录日志，不向用户暴露细节）
	ErrorTypeInfrastructure ErrorType = "INFRASTRUCTURE"
	// ErrorTypeSystem 系统级错误（如panic、数组越界等，需要记录堆栈）
	ErrorTypeSystem ErrorType = "SYSTEM"
)

// ErrorCode 业务错误码
type ErrorCode int

const (
	// ========== 基础错误码 (200-599) ==========
	CodeSuccess       ErrorCode = 200
	CodeBadRequest    ErrorCode = 400
	CodeUnauthorized  ErrorCode = 401
	CodeForbidden     ErrorCode = 403
	CodeNotFound      ErrorCode = 404
	CodeInternalError ErrorCode = 500

	// ========== 业务模块错误码 (10000-19999) ==========
	// 用户模块 (10xxx)
	CodeUserNotFound    ErrorCode = 10001
	CodeInvalidPassword ErrorCode = 10002
	CodeUserDisabled    ErrorCode = 10003
	CodeUserNameExists  ErrorCode = 10004
	CodePhoneExists     ErrorCode = 10005
	CodeEmailExists     ErrorCode = 10006
	CodeNoPermission    ErrorCode = 10007

	// 角色模块 (11xxx)
	CodeRoleNotFound   ErrorCode = 11001
	CodeRoleNameExists ErrorCode = 11002

	// 菜单模块 (12xxx)
	CodeMenuNotFound ErrorCode = 12001

	// ========== 基础设施错误码 (20000-29999) ==========
	// 数据库错误 (200xx)
	CodeDatabaseError          ErrorCode = 20001
	CodeDatabaseTimeout        ErrorCode = 20002
	CodeDatabaseConnectionLost ErrorCode = 20003

	// 缓存错误 (201xx)
	CodeRedisError   ErrorCode = 20101
	CodeRedisTimeout ErrorCode = 20102

	// 存储错误 (202xx)
	CodeS3Error        ErrorCode = 20201
	CodeS3UploadFailed ErrorCode = 20202

	// ========== 系统级错误码 (30000-39999) ==========
	CodePanicError       ErrorCode = 30001
	CodeIndexOutOfBounds ErrorCode = 30002
	CodeNullPointerError ErrorCode = 30003
)

// AppError 结构化业务错误
type AppError struct {
	Type    ErrorType   `json:"-"`                 // 错误类型（不暴露给前端）
	Code    ErrorCode   `json:"code"`              // 业务错误码
	Message string      `json:"message"`           // 错误描述（面向用户）
	RawErr  error       `json:"-"`                 // 原始错误（不暴露给前端）
	Details interface{} `json:"details,omitempty"` // 错误详情
}

func (e *AppError) Error() string {
	if e.RawErr != nil {
		return fmt.Sprintf("[%s][%d] %s: %v", e.Type, e.Code, e.Message, e.RawErr)
	}
	return fmt.Sprintf("[%s][%d] %s", e.Type, e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.RawErr
}

// NewBusiness 创建业务逻辑错误（可直接向用户展示）
// 示例: return errors.NewBusiness(errors.CodeUserNameExists, "用户名已存在")
func NewBusiness(code ErrorCode, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeBusiness,
		Code:    code,
		Message: message,
	}
}

// NewInfrastructure 创建基础设施错误（需记录日志，向用户返回统一文案）
// 示例: return errors.NewInfrastructure(errors.CodeDatabaseError, "数据库查询失败", err)
func NewInfrastructure(code ErrorCode, internalMsg string, rawErr error) *AppError {
	return &AppError{
		Type:    ErrorTypeInfrastructure,
		Code:    code,
		Message: internalMsg, // 内部消息，用于日志记录
		RawErr:  rawErr,
	}
}

// NewSystem 创建系统级错误（如panic、数组越界，需记录堆栈）
// 示例: return errors.NewSystem(errors.CodePanicError, "系统异常", err)
func NewSystem(code ErrorCode, internalMsg string, rawErr error) *AppError {
	return &AppError{
		Type:    ErrorTypeSystem,
		Code:    code,
		Message: internalMsg,
		RawErr:  rawErr,
	}
}

// New 创建一个新的 AppError（默认为业务错误，兼容旧代码）
// 推荐使用 NewBusiness 替代
func New(code ErrorCode, message string) *AppError {
	return NewBusiness(code, message)
}

// NewWithErr 创建一个带有原始错误的 AppError（默认为基础设施错误）
// 推荐使用 NewInfrastructure 或 NewSystem 替代
func NewWithErr(code ErrorCode, message string, err error) *AppError {
	return NewInfrastructure(code, message, err)
}

// Wrap 包装一个现有错误（自动判断类型）
func Wrap(err error, code ErrorCode, message string) *AppError {
	// 如果是系统错误（如 panic），标记为 System 类型
	if code >= 30000 {
		return NewSystem(code, message, err)
	}
	// 如果是基础设施错误
	if code >= 20000 {
		return NewInfrastructure(code, message, err)
	}
	// 默认为业务错误
	return NewBusiness(code, message)
}

// HTTPStatus 根据错误码返回对应的 HTTP 状态码
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInternalError:
		return http.StatusInternalServerError
	default:
		// 业务错误通常也返回 200 或特定的 4xx/5xx
		if e.Code >= 10000 && e.Code < 20000 {
			return http.StatusOK // 业务逻辑错误返回 200，由 code 区分
		}
		if e.Code >= 20000 {
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}
}

// IsBusiness 判断是否为业务错误
func (e *AppError) IsBusiness() bool {
	return e.Type == ErrorTypeBusiness
}

// IsInfrastructure 判断是否为基础设施错误
func (e *AppError) IsInfrastructure() bool {
	return e.Type == ErrorTypeInfrastructure
}

// IsSystem 判断是否为系统错误
func (e *AppError) IsSystem() bool {
	return e.Type == ErrorTypeSystem
}

// GetUserMessage 获取面向用户的错误信息
func (e *AppError) GetUserMessage() string {
	// 业务错误：直接返回自定义消息
	if e.IsBusiness() {
		return e.Message
	}

	// 基础设施错误：返回通用文案，隐藏内部细节
	if e.IsInfrastructure() {
		return "服务暂时不可用，请稍后重试"
	}

	// 系统错误：返回通用异常提示
	if e.IsSystem() {
		return "系统异常，请联系管理员"
	}

	return "系统繁忙，请稍后重试"
}
