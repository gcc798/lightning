package response

import (
	httpresponse "github.com/gcc798/microservice-kit/internal/httpresponse"
	"github.com/labstack/echo/v5"
)

const (
	CodeOK              = httpresponse.CodeOK
	CodeBadRequest      = httpresponse.CodeBadRequest
	CodeUnauthorized    = httpresponse.CodeUnauthorized
	CodeForbidden       = httpresponse.CodeForbidden
	CodeNotFound        = httpresponse.CodeNotFound
	CodeTimeout         = httpresponse.CodeTimeout
	CodeTooManyRequests = httpresponse.CodeTooManyRequests
	CodeServerError     = httpresponse.CodeServerError
	CodeInvalidParam    = httpresponse.CodeInvalidParam
)

type Response = httpresponse.Response

func Success(c *echo.Context, d interface{})                  { httpresponse.Success(c, d) }
func SuccessWithMsg(c *echo.Context, m string, d interface{}) { httpresponse.SuccessWithMsg(c, m, d) }
func Fail(c *echo.Context, m string)                          { httpresponse.Fail(c, m) }
func FailWithMsg(c *echo.Context, m string)                   { httpresponse.FailWithMsg(c, m) }
func BadRequest(c *echo.Context, m string)                    { httpresponse.BadRequest(c, m) }
func Unauthorized(c *echo.Context, m string)                  { httpresponse.Unauthorized(c, m) }
func Forbidden(c *echo.Context, m string)                     { httpresponse.Forbidden(c, m) }
func NotFound(c *echo.Context, m string)                      { httpresponse.NotFound(c, m) }
func InternalServerError(c *echo.Context, m string)           { httpresponse.InternalServerError(c, m) }
func Error(c *echo.Context, e error)                          { httpresponse.Error(c, e) }
func SuccessCode(c *echo.Context, code int, d interface{})    { httpresponse.SuccessCode(c, code, d) }
func FailCode(c *echo.Context, code int, m string)            { httpresponse.FailCode(c, code, m) }
