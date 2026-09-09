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

// ConfigController 配置控制器接口
type ConfigController interface {
	CreateConfig(ctx *echo.Context)        // 创建配置
	UpdateConfig(ctx *echo.Context)        // 更新配置
	DeleteConfig(ctx *echo.Context)        // 删除配置
	BatchDeleteConfig(ctx *echo.Context)   // 批量删除配置
	GetConfigById(ctx *echo.Context)       // 根据ID查询配置
	PageConfig(ctx *echo.Context)          // 分页查询配置列表
	GetConfigByCode(ctx *echo.Context)     // 根据编码获取配置列表
	GetConfigDataByCode(ctx *echo.Context) // 根据编码获取配置数据
}

type configController struct {
	ctr           container.Container
	configService sys.ConfigService
}

// NewConfigController 创建组件实例。
func NewConfigController(c container.Container) ConfigController {
	return &configController{
		ctr:           c,
		configService: sys.NewConfigService(c.GetDB(), c.GetLogger(), c.GetRuntimeConfig(), c),
	}
}

// CreateConfig 创建配置
//
//	@Summary		创建配置
//	@Description	创建新的配置数据，支持存储JSON格式的配置信息
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateConfigRequest	true	"创建配置请求"
//	@Success		200		{object}	response.Response			"创建成功"
//	@Failure		400		{object}	response.Response			"请求参数错误"
//	@Failure		500		{object}	response.Response			"服务器内部错误"
//	@Router			/api/v1/config [post]
//	@Security		Bearer
func (c *configController) CreateConfig(ctx *echo.Context) {
	var req request.CreateConfigRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.configService.Create(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "创建配置成功", nil)
}

// UpdateConfig 更新配置
//
//	@Summary		更新配置
//	@Description	更新配置数据
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.UpdateConfigRequest	true	"更新配置请求"
//	@Success		200		{object}	response.Response			"更新成功"
//	@Failure		400		{object}	response.Response			"请求参数错误"
//	@Failure		500		{object}	response.Response			"服务器内部错误"
//	@Router			/api/v1/config [put]
//	@Security		Bearer
func (c *configController) UpdateConfig(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的配置ID")
		return
	}

	var req request.UpdateConfigRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}
	req.ID = id

	if err := c.configService.Update(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新配置成功", nil)
}

// DeleteConfig 删除配置
//
//	@Summary		删除配置
//	@Description	删除单个配置数据
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"配置ID"
//	@Success		200	{object}	response.Response	"删除成功"
//	@Failure		400	{object}	response.Response	"请求参数错误"
//	@Failure		500	{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/config/{id} [delete]
//	@Security		Bearer
func (c *configController) DeleteConfig(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的配置ID")
		return
	}

	if err := c.configService.Delete(ctx.Request().Context(), id); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除配置成功", nil)
}

// BatchDeleteConfig 批量删除配置
//
//	@Summary		批量删除配置
//	@Description	批量删除配置数据
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.BatchDeleteConfigRequest	true	"批量删除请求"
//	@Success		200		{object}	response.Response					"删除成功"
//	@Failure		400		{object}	response.Response					"请求参数错误"
//	@Failure		500		{object}	response.Response					"服务器内部错误"
//	@Router			/api/v1/config/batch [delete]
//	@Security		Bearer
func (c *configController) BatchDeleteConfig(ctx *echo.Context) {
	var req request.BatchDeleteConfigRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.configService.BatchDelete(ctx.Request().Context(), req.IDs); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "批量删除配置成功", nil)
}

// GetConfigById 根据ID查询配置
//
//	@Summary		根据ID查询配置
//	@Description	根据配置ID查询配置详情
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int												true	"配置ID"
//	@Success		200	{object}	response.Response{data=response.ConfigResponse}	"查询成功"
//	@Failure		400	{object}	response.Response								"请求参数错误"
//	@Failure		500	{object}	response.Response								"服务器内部错误"
//	@Router			/api/v1/config/{id} [get]
//	@Security		Bearer
func (c *configController) GetConfigById(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的配置ID")
		return
	}

	config, err := c.configService.GetById(ctx.Request().Context(), id)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, response.ToConfigResponse(config))
}

// PageConfig 分页查询配置列表
//
//	@Summary		分页查询配置列表
//	@Description	使用 Paginator 分页查询配置列表，支持按类型、名称筛选
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {token}"
//	@Param			body			body		request.PageConfigRequest	true	"查询参数"
//	@Success		200				{object}	response.Response{data=object}	"查询成功"
//	@Failure		400				{object}	response.Response			"请求参数错误"
//	@Router			/api/v1/config/page [post]
//	@Security		Bearer
func (c *configController) PageConfig(ctx *echo.Context) {
	var req request.PageConfigRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	page, err := c.configService.Page(
		ctx.Request().Context(),
		req.PageNum,
		req.PageSize,
		req.Code,
		req.Name,
	)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, page)
}

// GetConfigByCode 根据编码获取配置列表
//
//	@Summary		根据编码获取配置列表
//	@Description	根据配置编码获取配置列表
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string												true	"配置编码"	example("system_settings")
//	@Success		200		{object}	response.Response{data=response.ConfigResponse}	"查询成功"
//	@Failure		400		{object}	response.Response									"请求参数错误"
//	@Failure		500		{object}	response.Response									"服务器内部错误"
//	@Router			/api/v1/config/code [get]
func (c *configController) GetConfigByCode(ctx *echo.Context) {
	var req request.GetConfigByCodeRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	config, err := c.configService.GetByCode(ctx.Request().Context(), req.Code)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, response.ToConfigResponse(config))
}

// GetConfigDataByCode 根据编码获取配置数据
//
//	@Summary		根据编码获取配置数据
//	@Description	根据配置编码获取配置的data字段（JSON格式），返回第一个匹配的配置
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string												true	"配置编码"	example("system_settings")
//	@Success		200		{object}	response.Response{data=response.ConfigDataResponse}	"查询成功"
//	@Failure		400		{object}	response.Response									"请求参数错误"
//	@Failure		500		{object}	response.Response									"服务器内部错误"
//	@Router			/api/v1/config/data [get]
func (c *configController) GetConfigDataByCode(ctx *echo.Context) {
	var req request.GetConfigByCodeRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	data, err := c.configService.GetDataByCode(ctx.Request().Context(), req.Code)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, response.ConfigDataResponse{
		Code: req.Code,
		Data: data,
	})
}
