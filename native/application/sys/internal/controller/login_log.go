package controller

import (
	"strconv"

	sys "github.com/gcc798/lightning/application/sys/internal/domain"
	"github.com/gcc798/lightning/application/sys/internal/request"
	"github.com/gcc798/lightning/application/sys/internal/response"
	"github.com/gcc798/lightning/internal/container"
	_ "github.com/gcc798/lightning/internal/utils/pagination"
	"github.com/labstack/echo/v5"
)

// LoginLogController 定义业务数据结构。
type LoginLogController interface {
	CreateLoginLog(ctx *echo.Context)      // 创建登录日志
	UpdateLoginLog(ctx *echo.Context)      // 更新登录日志
	DeleteLoginLog(ctx *echo.Context)      // 删除登录日志
	BatchDeleteLoginLog(ctx *echo.Context) // 批量删除登录日志
	GetLoginLogById(ctx *echo.Context)     // 根据ID查询登录日志
	PageLoginLog(ctx *echo.Context)        // 分页查询登录日志列表
	CleanLoginLog(ctx *echo.Context)       // 清理登录日志
}

type loginLogController struct {
	ctr             container.Container
	loginLogService sys.LoginLogService
}

// NewLoginLogController 创建组件实例。
func NewLoginLogController(c container.Container) LoginLogController {
	return &loginLogController{
		ctr:             c,
		loginLogService: sys.NewLoginLogService(c.GetDB(), c.GetLogger()),
	}
}

// CreateLoginLog 创建登录日志
//
//	@Summary		创建登录日志
//	@Description	创建新的登录日志记录
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateLoginLogRequest	true	"创建登录日志请求"
//	@Success		200		{object}	response.Response				"创建成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/loginLog [post]
//	@Security		Bearer
func (c *loginLogController) CreateLoginLog(ctx *echo.Context) {
	var req request.CreateLoginLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.loginLogService.Create(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "创建登录日志成功", nil)
}

// UpdateLoginLog 更新登录日志
//
//	@Summary		更新登录日志
//	@Description	更新登录日志记录
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.UpdateLoginLogRequest	true	"更新登录日志请求"
//	@Success		200		{object}	response.Response				"更新成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/loginLog [put]
//	@Security		Bearer
func (c *loginLogController) UpdateLoginLog(ctx *echo.Context) {
	var req request.UpdateLoginLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.loginLogService.Update(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新登录日志成功", nil)
}

// DeleteLoginLog 删除登录日志
//
//	@Summary		删除登录日志
//	@Description	删除单个登录日志记录
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"日志ID"
//	@Success		200	{object}	response.Response	"删除成功"
//	@Failure		400	{object}	response.Response	"请求参数错误"
//	@Failure		500	{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/loginLog/{id} [delete]
//	@Security		Bearer
func (c *loginLogController) DeleteLoginLog(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的日志ID")
		return
	}

	if err := c.loginLogService.Delete(ctx.Request().Context(), id); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除登录日志成功", nil)
}

// BatchDeleteLoginLog 批量删除登录日志
//
//	@Summary		批量删除登录日志
//	@Description	批量删除登录日志记录
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.BatchDeleteLoginLogRequest	true	"批量删除请求"
//	@Success		200		{object}	response.Response					"删除成功"
//	@Failure		400		{object}	response.Response					"请求参数错误"
//	@Failure		500		{object}	response.Response					"服务器内部错误"
//	@Router			/api/v1/loginLog/batch [delete]
//	@Security		Bearer
func (c *loginLogController) BatchDeleteLoginLog(ctx *echo.Context) {
	var req request.BatchDeleteLoginLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.loginLogService.BatchDelete(ctx.Request().Context(), req.IDs); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "批量删除登录日志成功", nil)
}

// GetLoginLogById 根据ID查询登录日志
//
//	@Summary		根据ID查询登录日志
//	@Description	根据日志ID查询登录日志详情
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int													true	"日志ID"
//	@Success		200	{object}	response.Response{data=response.LoginLogResponse}	"查询成功"
//	@Failure		400	{object}	response.Response									"请求参数错误"
//	@Failure		500	{object}	response.Response									"服务器内部错误"
//	@Router			/api/v1/loginLog/{id} [get]
//	@Security		Bearer
func (c *loginLogController) GetLoginLogById(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的日志ID")
		return
	}

	log, err := c.loginLogService.GetById(ctx.Request().Context(), id)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, response.ToLoginLogResponse(log))
}

// PageLoginLog 分页查询登录日志列表
//
//	@Summary		分页查询登录日志列表
//	@Description	分页查询登录日志列表，支持按用户名、IP、状态、时间范围筛选
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			pageNum			query		int					true	"页码"	minimum(1)
//	@Param			pageSize		query		int					true	"每页数量"	minimum(1)
//	@Param			orderByColumn	query		string				false	"排序列"	example("loginTime")
//	@Param			isAsc			query		string				false	"排序方向"	Enums(asc,desc)	example("desc")
//	@Param			userName		query		string				false	"用户名（模糊查询）"
//	@Param			ipaddr			query		string				false	"登录IP（模糊查询）"
//	@Param			status			query		int					false	"登录状态：0成功 1失败，不传或null表示全部"	Enums(0,1)
//	@Param			startTime		query		string				false	"开始时间"				example("2024-01-01 00:00:00")
//	@Param			endTime			query		string				false	"结束时间"				example("2024-12-31 23:59:59")
//	@Success		200				{object}	response.Response	"查询成功"
//	@Failure		400				{object}	response.Response	"请求参数错误"
//	@Failure		500				{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/system/loginLog/page [post]
//	@Security		Bearer
func (c *loginLogController) PageLoginLog(ctx *echo.Context) {
	var req request.PageLoginLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	page, err := c.loginLogService.Page(ctx.Request().Context(), &req)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, page)
}

// CleanLoginLog 清理登录日志
//
//	@Summary		清理登录日志
//	@Description	清理指定天数之前的登录日志
//	@Tags			登录日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CleanLoginLogRequest	true	"清理日志请求"
//	@Success		200		{object}	response.Response				"清理成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/loginLog/clean [post]
//	@Security		Bearer
func (c *loginLogController) CleanLoginLog(ctx *echo.Context) {
	var req request.CleanLoginLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	count, err := c.loginLogService.CleanOldLogs(ctx.Request().Context(), req.Days)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "清理登录日志成功", map[string]interface{}{
		"count": count,
		"days":  req.Days,
	})
}
