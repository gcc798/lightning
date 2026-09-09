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

// OperLogController 定义业务数据结构。
type OperLogController interface {
	CreateOperLog(ctx *echo.Context)      // 创建操作日志
	UpdateOperLog(ctx *echo.Context)      // 更新操作日志
	DeleteOperLog(ctx *echo.Context)      // 删除操作日志
	BatchDeleteOperLog(ctx *echo.Context) // 批量删除操作日志
	GetOperLogById(ctx *echo.Context)     // 根据ID查询操作日志
	PageOperLog(ctx *echo.Context)        // 分页查询操作日志列表
	CleanOperLog(ctx *echo.Context)       // 清理操作日志
}

type operLogController struct {
	ctr            container.Container
	operLogService sys.OperLogService
}

// NewOperLogController 创建组件实例。
func NewOperLogController(c container.Container) OperLogController {
	return &operLogController{
		ctr:            c,
		operLogService: sys.NewOperLogService(c.GetDB(), c.GetLogger()),
	}
}

// CreateOperLog 创建操作日志
//
//	@Summary		创建操作日志
//	@Description	创建新的操作日志记录
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateOperLogRequest	true	"创建操作日志请求"
//	@Success		200		{object}	response.Response				"创建成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/operLog [post]
//	@Security		Bearer
func (c *operLogController) CreateOperLog(ctx *echo.Context) {
	var req request.CreateOperLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.operLogService.Create(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "创建操作日志成功", nil)
}

// UpdateOperLog 更新操作日志
//
//	@Summary		更新操作日志
//	@Description	更新操作日志记录
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.UpdateOperLogRequest	true	"更新操作日志请求"
//	@Success		200		{object}	response.Response				"更新成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/operLog [put]
//	@Security		Bearer
func (c *operLogController) UpdateOperLog(ctx *echo.Context) {
	var req request.UpdateOperLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.operLogService.Update(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新操作日志成功", nil)
}

// DeleteOperLog 删除操作日志
//
//	@Summary		删除操作日志
//	@Description	删除单个操作日志记录
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"日志ID"
//	@Success		200	{object}	response.Response	"删除成功"
//	@Failure		400	{object}	response.Response	"请求参数错误"
//	@Failure		500	{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/operLog/{id} [delete]
//	@Security		Bearer
func (c *operLogController) DeleteOperLog(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的日志ID")
		return
	}

	if err := c.operLogService.Delete(ctx.Request().Context(), id); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除操作日志成功", nil)
}

// BatchDeleteOperLog 批量删除操作日志
//
//	@Summary		批量删除操作日志
//	@Description	批量删除操作日志记录
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.BatchDeleteOperLogRequest	true	"批量删除请求"
//	@Success		200		{object}	response.Response					"删除成功"
//	@Failure		400		{object}	response.Response					"请求参数错误"
//	@Failure		500		{object}	response.Response					"服务器内部错误"
//	@Router			/api/v1/operLog/batch [delete]
//	@Security		Bearer
func (c *operLogController) BatchDeleteOperLog(ctx *echo.Context) {
	var req request.BatchDeleteOperLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.operLogService.BatchDelete(ctx.Request().Context(), req.IDs); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "批量删除操作日志成功", nil)
}

// GetOperLogById 根据ID查询操作日志
//
//	@Summary		根据ID查询操作日志
//	@Description	根据日志ID查询操作日志详情
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int													true	"日志ID"
//	@Success		200	{object}	response.Response{data=response.OperLogResponse}	"查询成功"
//	@Failure		400	{object}	response.Response									"请求参数错误"
//	@Failure		500	{object}	response.Response									"服务器内部错误"
//	@Router			/api/v1/operLog/{id} [get]
//	@Security		Bearer
func (c *operLogController) GetOperLogById(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的日志ID")
		return
	}

	log, err := c.operLogService.GetById(ctx.Request().Context(), id)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, response.ToOperLogResponse(log))
}

// PageOperLog 分页查询操作日志列表
//
//	@Summary		分页查询操作日志列表
//	@Description	分页查询操作日志列表，支持按标题、操作者、业务类型、状态、时间范围筛选
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			body			body		request.PageOperLogRequest	true	"查询参数"
//	@Success		200				{object}	response.Response	"查询成功"
//	@Failure		400				{object}	response.Response	"请求参数错误"
//	@Failure		500				{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/operLog/page [post]
//	@Security		Bearer
func (c *operLogController) PageOperLog(ctx *echo.Context) {
	var req request.PageOperLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	page, err := c.operLogService.Page(ctx.Request().Context(), &req)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, page)
}

// CleanOperLog 清理操作日志
//
//	@Summary		清理操作日志
//	@Description	清理指定天数之前的操作日志
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CleanOperLogRequest	true	"清理日志请求"
//	@Success		200		{object}	response.Response			"清理成功"
//	@Failure		400		{object}	response.Response			"请求参数错误"
//	@Failure		500		{object}	response.Response			"服务器内部错误"
//	@Router			/api/v1/operLog/clean [post]
//	@Security		Bearer
func (c *operLogController) CleanOperLog(ctx *echo.Context) {
	var req request.CleanOperLogRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	count, err := c.operLogService.CleanOldLogs(ctx.Request().Context(), req.Days)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "清理操作日志成功", map[string]interface{}{
		"count": count,
		"days":  req.Days,
	})
}
