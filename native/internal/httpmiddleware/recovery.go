package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gcc798/microservice-kit/internal/httpresponse"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	apperrors "github.com/gcc798/microservice-kit/internal/utils/errors"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// Recovery 全局 Panic 恢复中间件
// 捕获系统级错误（如数组越界、空指针等），记录堆栈，返回统一文案
// 参数：logger 需要传入 *zap.Logger，调用方使用 logger.Get() 获取
func Recovery(logger logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) (handlerErr error) {
			defer func() {
				if err := recover(); err != nil {
					// 记录详细的 panic 信息和堆栈
					logging.WithContext(c.Request().Context(), logger).Error("系统 Panic 捕获",
						zap.Any("error", err),
						zap.String("path", c.Request().URL.Path),
						zap.String("method", c.Request().Method),
						zap.Any("params", c.Request().URL.Query()),
						zap.String("stack", string(debug.Stack())),
					)

					// 构造系统错误
					var panicErr error
					switch e := err.(type) {
					case error:
						panicErr = e
					case string:
						panicErr = fmt.Errorf("%s", e)
					default:
						panicErr = fmt.Errorf("%v", e)
					}

					// 创建系统级错误
					systemErr := apperrors.NewSystem(
						apperrors.CodePanicError,
						"系统发生严重错误",
						panicErr,
					)

					// 返回统一的系统异常响应
					handlerErr = c.JSON(http.StatusInternalServerError, response.Response{
						Code: int(systemErr.Code),
						Msg:  systemErr.GetUserMessage(), // "系统异常，请联系管理员"
						Data: nil,
					})

				}
			}()

			return next(c)
		}
	}
}
