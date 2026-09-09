package controller

import (
	iam "github.com/gcc798/lightning/application/iam/internal/domain"
	"github.com/gcc798/lightning/application/iam/internal/request"
	"github.com/gcc798/lightning/application/iam/internal/response"
	"github.com/gcc798/lightning/internal/container"
	"github.com/gcc798/lightning/internal/httputils"
	_ "github.com/gcc798/lightning/internal/utils/pagination"
	"github.com/gcc798/lightning/internal/validator"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// OrgController 组织控制器接口
type OrgController interface {
	Create(c *echo.Context)      // 创建组织
	Update(c *echo.Context)      // 更新组织
	Delete(c *echo.Context)      // 删除组织
	BatchDelete(c *echo.Context) // 批量删除组织
	GetById(c *echo.Context)     // 根据ID查询组织
	GetTree(c *echo.Context)     // 获取组织树
	PageOrg(c *echo.Context)     // 分页查询组织列表
}

type orgController struct {
	ctr        container.Container
	base       *BaseController
	orgService iam.OrgService
}

// NewOrgController 创建组件实例。
func NewOrgController(c container.Container) OrgController {
	return &orgController{
		ctr:        c,
		base:       NewBaseController(c),
		orgService: iam.NewOrgService(c.GetDB(), c.GetLogger()),
	}
}

// Create 创建组织
//
//	@Summary		创建组织
//	@Description	创建新组织，需要管理员权限
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			request			body		request.CreateOrgRequest		true	"创建组织请求"
//	@Success		200				{object}	response.Response{data=int64}	"创建成功，返回组织ID"
//	@Failure		400				{object}	response.Response				"参数错误"
//	@Failure		401				{object}	response.Response				"未授权"
//	@Failure		500				{object}	response.Response				"服务器错误"
//	@Router			/api/v1/org [post]
//	@Security		Bearer
func (h *orgController) Create(c *echo.Context) {
	var req request.CreateOrgRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.TranslateValidationError(err))
		return
	}

	currentUserId, _ := h.base.GetUserId(c)
	req.CreateBy = currentUserId
	req.UpdateBy = currentUserId

	orgId, err := h.orgService.Create(c.Request().Context(), &req)
	if err != nil {
		h.ctr.GetLogger().Error("创建组织失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, orgId)
}

// Update 更新组织
//
//	@Summary		更新组织信息
//	@Description	更新指定组织的信息，需要管理员权限
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			id				path		int								true	"组织ID"
//	@Param			request			body		request.UpdateOrgRequest		true	"更新组织请求"
//	@Success		200				{object}	response.Response{data=string}	"更新成功"
//	@Failure		400				{object}	response.Response				"参数错误"
//	@Failure		401				{object}	response.Response				"未授权"
//	@Failure		404				{object}	response.Response				"组织不存在"
//	@Failure		500				{object}	response.Response				"服务器错误"
//	@Router			/api/v1/org/{id} [put]
//	@Security		Bearer
func (h *orgController) Update(c *echo.Context) {
	orgId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	var req request.UpdateOrgRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.TranslateValidationError(err))
		return
	}

	req.OrgId = orgId
	currentUserId, _ := h.base.GetUserId(c)
	req.UpdateBy = currentUserId

	if err := h.orgService.Update(c.Request().Context(), &req); err != nil {
		h.ctr.GetLogger().Error("更新组织失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// Delete 删除组织
//
//	@Summary		删除组织
//	@Description	删除指定组织，需要管理员权限
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			id				path		int								true	"组织ID"
//	@Success		200				{object}	response.Response{data=string}	"删除成功"
//	@Failure		400				{object}	response.Response				"参数错误"
//	@Failure		401				{object}	response.Response				"未授权"
//	@Failure		404				{object}	response.Response				"组织不存在"
//	@Failure		500				{object}	response.Response				"服务器错误"
//	@Router			/api/v1/org/{id} [delete]
//	@Security		Bearer
func (h *orgController) Delete(c *echo.Context) {
	orgId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	if err := h.orgService.Delete(c.Request().Context(), orgId); err != nil {
		h.ctr.GetLogger().Error("删除组织失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// BatchDelete 批量删除组织
//
//	@Summary		批量删除组织
//	@Description	批量删除多个组织，需要管理员权限
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			request			body		request.BatchDeleteOrgsRequest	true	"批量删除请求"
//	@Success		200				{object}	response.Response{data=string}	"删除成功"
//	@Failure		400				{object}	response.Response				"参数错误"
//	@Failure		401				{object}	response.Response				"未授权"
//	@Failure		500				{object}	response.Response				"服务器错误"
//	@Router			/api/v1/org/batch [delete]
//	@Security		Bearer
func (h *orgController) BatchDelete(c *echo.Context) {
	var req request.BatchDeleteOrgsRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.TranslateValidationError(err))
		return
	}

	if err := h.orgService.BatchDelete(c.Request().Context(), req.IDs); err != nil {
		h.ctr.GetLogger().Error("批量删除组织失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// GetById 根据ID查询组织
//
//	@Summary		查询组织详情
//	@Description	根据组织ID查询组织详细信息
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			id				path		int								true	"组织ID"
//	@Success		200				{object}	response.Response{data=object}	"查询成功"
//	@Failure		400				{object}	response.Response				"参数错误"
//	@Failure		401				{object}	response.Response				"未授权"
//	@Failure		404				{object}	response.Response				"组织不存在"
//	@Failure		500				{object}	response.Response				"服务器错误"
//	@Router			/api/v1/org/{id} [get]
//	@Security		Bearer
func (h *orgController) GetById(c *echo.Context) {
	orgId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	org, err := h.orgService.GetById(c.Request().Context(), orgId)
	if err != nil {
		h.ctr.GetLogger().Error("查询组织失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, org)
}

// GetTree 获取组织树
//
//	@Summary		获取组织树
//	@Description	获取完整的组织树结构
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer {token}"
//	@Success		200				{object}	response.Response{data=[]object}	"查询成功"
//	@Failure		401				{object}	response.Response					"未授权"
//	@Failure		500				{object}	response.Response					"服务器错误"
//	@Router			/api/v1/org/tree [get]
//	@Security		Bearer
func (h *orgController) GetTree(c *echo.Context) {
	orgs, err := h.orgService.GetTree(c.Request().Context())
	if err != nil {
		h.ctr.GetLogger().Error("查询组织树失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, orgs)
}

// PageOrg 分页查询组织列表
//
//	@Summary		分页查询组织列表
//	@Description	使用 Paginator 分页查询组织列表
//	@Tags			组织管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer {token}"
//	@Param			body			body		request.PageOrgsRequest	true	"查询参数"
//	@Success		200				{object}	response.Response{data=object}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/org/page [post]
//	@Security		Bearer
func (h *orgController) PageOrg(c *echo.Context) {
	var req request.PageOrgsRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.TranslateValidationError(err))
		return
	}

	if req.PageNum < 1 {
		req.PageNum = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	page, err := h.orgService.Page(c.Request().Context(), req.PageNum, req.PageSize, req.OrgName, req.OrgCode, req.Status, req.ParentId)
	if err != nil {
		h.ctr.GetLogger().Error("分页查询组织列表失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, page)
}
