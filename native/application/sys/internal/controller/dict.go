package controller

import (
	"strconv"

	sys "github.com/gcc798/microservice-kit/application/sys/internal/domain"
	"github.com/gcc798/microservice-kit/application/sys/internal/request"
	"github.com/gcc798/microservice-kit/application/sys/internal/response"
	"github.com/gcc798/microservice-kit/internal/container"
	_ "github.com/gcc798/microservice-kit/internal/utils/pagination"
	"github.com/labstack/echo/v5"
)

// DictController 字典控制器接口
type DictController interface {
	CreateDict(ctx *echo.Context)      // 创建字典
	UpdateDict(ctx *echo.Context)      // 更新字典
	DeleteDict(ctx *echo.Context)      // 删除字典
	BatchDeleteDict(ctx *echo.Context) // 批量删除字典
	GetDictById(ctx *echo.Context)     // 根据ID查询字典
	PageDict(ctx *echo.Context)        // 分页查询字典列表
	GetDictByType(ctx *echo.Context)   // 根据类型获取字典列表
	GetDictLabel(ctx *echo.Context)    // 根据类型和键值获取标签
}

type dictController struct {
	ctr         container.Container
	dictService sys.DictService
}

// NewDictController 创建组件实例。
func NewDictController(c container.Container) DictController {
	return &dictController{
		ctr:         c,
		dictService: sys.NewDictService(c.GetDB(), c.GetLogger()),
	}
}

// CreateDict 创建字典
//
//	@Summary		创建字典
//	@Description	创建新的字典数据，支持树形结构（通过parentId）
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateDictRequest	true	"创建字典请求"
//	@Success		200		{object}	response.Response			"创建成功"
//	@Failure		400		{object}	response.Response			"请求参数错误"
//	@Failure		500		{object}	response.Response			"服务器内部错误"
//	@Router			/api/v1/dict [post]
//	@Security		Bearer
func (c *dictController) CreateDict(ctx *echo.Context) {
	var req request.CreateDictRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.dictService.Create(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "创建字典成功", nil)
}

// UpdateDict 更新字典
//
//	@Summary		更新字典
//	@Description	更新字典数据
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.UpdateDictRequest	true	"更新字典请求"
//	@Success		200		{object}	response.Response			"更新成功"
//	@Failure		400		{object}	response.Response			"请求参数错误"
//	@Failure		500		{object}	response.Response			"服务器内部错误"
//	@Router			/api/v1/dict [put]
//	@Security		Bearer
func (c *dictController) UpdateDict(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典ID")
		return
	}

	var req request.UpdateDictRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}
	req.ID = id

	if err := c.dictService.Update(ctx.Request().Context(), &req); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新字典成功", nil)
}

// DeleteDict 删除字典
//
//	@Summary		删除字典
//	@Description	删除单个字典数据（如果有子字典则无法删除）
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"字典ID"
//	@Success		200	{object}	response.Response	"删除成功"
//	@Failure		400	{object}	response.Response	"请求参数错误"
//	@Failure		500	{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/dict/{id} [delete]
//	@Security		Bearer
func (c *dictController) DeleteDict(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典ID")
		return
	}

	if err := c.dictService.Delete(ctx.Request().Context(), id); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除字典成功", nil)
}

// BatchDeleteDict 批量删除字典
//
//	@Summary		批量删除字典
//	@Description	批量删除字典数据（如果有子字典则无法删除）
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.BatchDeleteDictRequest	true	"批量删除请求"
//	@Success		200		{object}	response.Response				"删除成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/dict/batch [delete]
//	@Security		Bearer
func (c *dictController) BatchDeleteDict(ctx *echo.Context) {
	var req request.BatchDeleteDictRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.dictService.BatchDelete(ctx.Request().Context(), req.IDs); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "批量删除字典成功", nil)
}

// GetDictById 根据ID查询字典
//
//	@Summary		根据ID查询字典
//	@Description	根据字典ID查询字典详情
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int													true	"字典ID"
//	@Success		200	{object}	response.Response{data=response.DictDataResponse}	"查询成功"
//	@Failure		400	{object}	response.Response									"请求参数错误"
//	@Failure		500	{object}	response.Response									"服务器内部错误"
//	@Router			/api/v1/dict/{id} [get]
//	@Security		Bearer
func (c *dictController) GetDictById(ctx *echo.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典ID")
		return
	}

	dict, err := c.dictService.GetById(ctx.Request().Context(), id)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, response.ToDictDataResponse(dict))
}

// PageDict 分页查询字典列表
//
//	@Summary		分页查询字典列表
//	@Description	分页查询字典列表，支持按类型、标签、状态筛选
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			pageNum			query		int					true	"页码"	minimum(1)
//	@Param			pageSize		query		int					true	"每页数量"	minimum(1)
//	@Param			orderByColumn	query		string				false	"排序列"	example("sort")
//	@Param			isAsc			query		string				false	"排序方向"	Enums(asc,desc)	example("asc")
//	@Param			dictType		query		string				false	"字典类型"
//	@Param			dictLabel		query		string				false	"字典标签（模糊查询）"
//	@Param			status			query		int					false	"状态：0正常 1停用 -1全部"	Enums(0,1,-1)
//	@Success		200				{object}	response.Response	"查询成功"
//	@Failure		400				{object}	response.Response	"请求参数错误"
//	@Failure		500				{object}	response.Response	"服务器内部错误"
//	@Router			/api/v1/system/dict/page [get]
//	@Security		Bearer
func (c *dictController) PageDict(ctx *echo.Context) {
	var req request.PageDictRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	page, err := c.dictService.Page(ctx.Request().Context(), &req)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, page)
}

// GetDictByType 根据类型获取字典列表
//
//	@Summary		根据类型获取字典列表
//	@Description	根据字典类型获取字典列表，用于前端下拉框等场景。支持获取子字典列表（通过parentId参数）
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			dictType	query		string												true	"字典类型"				example("sys_user_sex")
//	@Param			parentId	query		int													false	"父字典ID（可选，用于获取子字典）"	example(0)
//	@Success		200			{object}	response.Response{data=[]response.DictDataResponse}	"查询成功"
//	@Failure		400			{object}	response.Response									"请求参数错误"
//	@Failure		500			{object}	response.Response									"服务器内部错误"
//	@Router			/api/v1/dict/type [get]
func (c *dictController) GetDictByType(ctx *echo.Context) {
	var req request.GetDictByTypeRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	var dicts []response.DictDataResponse

	if req.ParentId != nil {
		dictList, err := c.dictService.GetByTypeAndParent(ctx.Request().Context(), req.DictType, *req.ParentId)
		if err != nil {
			response.Fail(ctx, err.Error())
			return
		}
		for _, dict := range dictList {
			dicts = append(dicts, response.ToDictDataResponse(&dict))
		}
	} else {
		dictList, err := c.dictService.GetByType(ctx.Request().Context(), req.DictType)
		if err != nil {
			response.Fail(ctx, err.Error())
			return
		}
		for _, dict := range dictList {
			dicts = append(dicts, response.ToDictDataResponse(&dict))
		}
	}

	response.Success(ctx, dicts)
}

// GetDictLabel 根据类型和键值获取标签
//
//	@Summary		根据类型和键值获取标签
//	@Description	根据字典类型和键值获取对应的标签（用于数据展示）
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			dictType	query		string							true	"字典类型"	example("sys_user_sex")
//	@Param			dictValue	query		string							true	"字典键值"	example("0")
//	@Success		200			{object}	response.Response{data=string}	"查询成功"
//	@Failure		400			{object}	response.Response				"请求参数错误"
//	@Failure		500			{object}	response.Response				"服务器内部错误"
//	@Router			/api/v1/dict/label [get]
func (c *dictController) GetDictLabel(ctx *echo.Context) {
	dictType := ctx.QueryParam("dictType")
	dictValue := ctx.QueryParam("dictValue")

	if dictType == "" || dictValue == "" {
		response.BadRequest(ctx, "dictType和dictValue不能为空")
		return
	}

	label, err := c.dictService.GetDictLabel(ctx.Request().Context(), dictType, dictValue)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Success(ctx, label)
}
