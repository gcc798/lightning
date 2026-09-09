package controller

import (
	"strconv"

	iam "github.com/gcc798/lightning/application/iam/internal/domain"
	"github.com/gcc798/lightning/application/iam/internal/request"
	"github.com/gcc798/lightning/application/iam/internal/response"
	"github.com/gcc798/lightning/internal/container"
	"github.com/labstack/echo/v5"
)

// ApiPermissionController API 权限控制器。
type ApiPermissionController interface {
	Tree(ctx *echo.Context)
	List(ctx *echo.Context)
	Create(ctx *echo.Context)
	Update(ctx *echo.Context)
	Delete(ctx *echo.Context)
	GetRolePermissions(ctx *echo.Context)
	AssignRolePermissions(ctx *echo.Context)
	GetUserPermissions(ctx *echo.Context)
	AssignUserPermissions(ctx *echo.Context)
}

type apiPermissionController struct {
	service iam.ApiPermissionService
}

func NewApiPermissionController(c container.Container) ApiPermissionController {
	return &apiPermissionController{
		service: iam.NewApiPermissionService(c.GetDB()),
	}
}

func (c *apiPermissionController) Tree(ctx *echo.Context) {
	tree, err := c.service.Tree(ctx.Request().Context())
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, tree)
}

func (c *apiPermissionController) List(ctx *echo.Context) {
	list, err := c.service.List(ctx.Request().Context())
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, list)
}

func (c *apiPermissionController) Create(ctx *echo.Context) {
	var req request.ApiPermissionSaveRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	userId := currentUserId(ctx)
	permission, err := c.service.Create(ctx.Request().Context(), &req, userId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, permission)
}

func (c *apiPermissionController) Update(ctx *echo.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "权限ID格式错误")
		return
	}
	var req request.ApiPermissionSaveRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	if err := c.service.Update(ctx.Request().Context(), id, &req, currentUserId(ctx)); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *apiPermissionController) Delete(ctx *echo.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "权限ID格式错误")
		return
	}
	if err := c.service.Delete(ctx.Request().Context(), id); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *apiPermissionController) GetRolePermissions(ctx *echo.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}
	grants, err := c.service.GetRolePermissions(ctx.Request().Context(), roleId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, grants)
}

func (c *apiPermissionController) AssignRolePermissions(ctx *echo.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}
	var req request.ApiPermissionAssignRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	if err := c.service.AssignRolePermissions(ctx.Request().Context(), roleId, req.PermissionIds, currentUserId(ctx)); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *apiPermissionController) GetUserPermissions(ctx *echo.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户ID格式错误")
		return
	}
	ids, err := c.service.GetUserPermissionIds(ctx.Request().Context(), userId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, ids)
}

func (c *apiPermissionController) AssignUserPermissions(ctx *echo.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户ID格式错误")
		return
	}
	var req request.ApiPermissionAssignRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	if err := c.service.AssignUserPermissions(ctx.Request().Context(), userId, req.PermissionIds, currentUserId(ctx)); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func currentUserId(ctx *echo.Context) int64 {
	value := ctx.Get("userId")
	if value == nil {
		return 0
	}
	userId, ok := value.(int64)
	if !ok {
		return 0
	}
	return userId
}
