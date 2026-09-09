package response

import (
	"net/http"
	"runtime/debug"

	"github.com/gcc798/lightning/internal/utils/errors"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

const (
	// CodeOK 定义业务常量。
	CodeOK = 200
	// CodeBadRequest 定义业务常量。
	CodeBadRequest = 400
	// CodeUnauthorized 定义业务常量。
	CodeUnauthorized = 401
	// CodeForbidden 定义业务常量。
	CodeForbidden = 403
	// CodeNotFound 定义业务常量。
	CodeNotFound = 404
	// CodeTimeout 定义业务常量。
	CodeTimeout = 408
	// CodeTooManyRequests 定义业务常量。
	CodeTooManyRequests = 429
	// CodeServerError 定义业务常量。
	CodeServerError = 500
	// CodeInvalidParam 定义业务常量。
	CodeInvalidParam = 400

	defaultSuccessMsg = "操作成功"
	defaultFailMsg    = "操作失败"
)

// Response 统一响应结构
//
//	@Description	API 统一响应格式
type Response struct {
	Code int         `json:"code" example:"200"`        // 业务状态码
	Msg  string      `json:"msg" example:"操作成功"`        // 响应消息
	Data interface{} `json:"data" swaggertype:"object"` // 响应数据
}

// Success 执行业务逻辑。
func Success(c *echo.Context, data interface{}) {
	c.JSON(200, Response{Code: CodeOK, Msg: defaultSuccessMsg, Data: data})
}

// SuccessWithMsg 执行业务逻辑。
func SuccessWithMsg(c *echo.Context, msg string, data interface{}) {
	c.JSON(200, Response{Code: CodeOK, Msg: msg, Data: data})
}

// Fail 执行业务逻辑。
func Fail(c *echo.Context, msg string) {
	if msg == "" {
		msg = defaultFailMsg
	}
	c.JSON(http.StatusInternalServerError, Response{Code: CodeServerError, Msg: msg})
}

// FailWithMsg 执行业务逻辑。
func FailWithMsg(c *echo.Context, msg string) { Fail(c, msg) }

// BadRequest 执行业务逻辑。
func BadRequest(c *echo.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{Code: CodeBadRequest, Msg: msg})
}

// Unauthorized 执行业务逻辑。
func Unauthorized(c *echo.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{Code: CodeUnauthorized, Msg: msg})
}

// Forbidden 执行业务逻辑。
func Forbidden(c *echo.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{Code: CodeForbidden, Msg: msg})
}

// NotFound 执行业务逻辑。
func NotFound(c *echo.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{Code: CodeNotFound, Msg: msg})
}

// InternalServerError 执行业务逻辑。
func InternalServerError(c *echo.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{Code: CodeServerError, Msg: msg})
}

// Error 处理结构化错误（增强版：支持类型区分和日志记录）
func Error(c *echo.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		// 获取 logger（从 context 中获取，如果没有则使用默认行为）
		loggerValue := c.Get("logger")
		var logger *zap.Logger
		if loggerValue != nil {
			if l, ok := loggerValue.(*zap.Logger); ok {
				logger = l
			}
		}

		// 根据错误类型进行不同处理
		switch {
		case appErr.IsBusiness():
			// 业务错误：直接返回自定义消息给用户
			if logger != nil {
				logger.Warn("业务逻辑错误",
					zap.Int("code", int(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.String("path", c.Request().URL.Path))
			}
			c.JSON(appErr.HTTPStatus(), Response{
				Code: int(appErr.Code),
				Msg:  appErr.GetUserMessage(), // 直接返回业务消息
				Data: appErr.Details,
			})

		case appErr.IsInfrastructure():
			// 基础设施错误：记录详细日志，返回统一文案
			if logger != nil {
				logger.Error("基础设施错误",
					zap.Int("code", int(appErr.Code)),
					zap.String("internal_message", appErr.Message),
					zap.Error(appErr.RawErr),
					zap.String("path", c.Request().URL.Path),
					zap.Any("params", c.Request().URL.Query()),
					zap.String("method", c.Request().Method))
			}
			c.JSON(http.StatusInternalServerError, Response{
				Code: int(appErr.Code),
				Msg:  appErr.GetUserMessage(), // 返回统一的用户友好文案
				Data: nil,
			})

		case appErr.IsSystem():
			// 系统错误：记录详细堆栈，返回通用异常提示
			if logger != nil {
				logger.Error("系统异常",
					zap.Int("code", int(appErr.Code)),
					zap.String("internal_message", appErr.Message),
					zap.Error(appErr.RawErr),
					zap.String("path", c.Request().URL.Path),
					zap.Any("params", c.Request().URL.Query()),
					zap.String("method", c.Request().Method),
					zap.String("stack", string(debug.Stack()))) // 记录堆栈
			}
			c.JSON(http.StatusInternalServerError, Response{
				Code: int(appErr.Code),
				Msg:  appErr.GetUserMessage(), // 返回"系统异常"统一文案
				Data: nil,
			})

		default:
			// 未知类型，按系统错误处理
			if logger != nil {
				logger.Error("未分类错误",
					zap.Int("code", int(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.Error(appErr.RawErr))
			}
			c.JSON(http.StatusInternalServerError, Response{
				Code: CodeServerError,
				Msg:  "系统繁忙，请稍后重试",
				Data: nil,
			})
		}
		return
	}

	// 非 AppError 类型（普通 error），按系统错误处理
	loggerValue := c.Get("logger")
	if loggerValue != nil {
		if logger, ok := loggerValue.(*zap.Logger); ok {
			logger.Error("未捕获的系统错误",
				zap.Error(err),
				zap.String("path", c.Request().URL.Path),
				zap.String("stack", string(debug.Stack())))
		}
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code: CodeServerError,
		Msg:  "系统异常，请联系管理员",
		Data: nil,
	})
}

// SuccessCode 执行业务逻辑。
func SuccessCode(c *echo.Context, code int, data interface{}) {
	c.JSON(200, Response{Code: code, Msg: defaultSuccessMsg, Data: data})
}

// FailCode 执行业务逻辑。
func FailCode(c *echo.Context, code int, msg string) {
	status := http.StatusOK
	if code >= 400 && code <= 599 {
		status = code
	}
	c.JSON(status, Response{Code: code, Msg: msg})
}

// ValidationFieldError 单个字段验证错误
type ValidationFieldError struct {
	Field   string `json:"field"`   // 字段名（英文）
	Message string `json:"message"` // 错误信息（中文）
}

// ValidationErrorResponse 验证错误响应
type ValidationErrorResponse struct {
	Code   int                    `json:"code"`
	Msg    string                 `json:"msg"`
	Errors []ValidationFieldError `json:"errors"` // 详细的字段错误列表
}

// FailValidation 返回验证错误（包含字段详情）
func FailValidation(c *echo.Context, code int, msg string, errors []ValidationFieldError) {
	c.JSON(http.StatusBadRequest, ValidationErrorResponse{Code: code, Msg: msg, Errors: errors})
}
