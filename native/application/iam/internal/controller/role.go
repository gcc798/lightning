package controller

import (
	"strconv"

	iam "github.com/gcc798/microservice-kit/application/iam/internal/domain"
	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	"github.com/gcc798/microservice-kit/application/iam/internal/request"
	"github.com/gcc798/microservice-kit/application/iam/internal/response"
	"github.com/gcc798/microservice-kit/internal/container"
	"github.com/gcc798/microservice-kit/internal/logger"
	_ "github.com/gcc798/microservice-kit/internal/utils/pagination"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// RoleController 角色控制器接口
type RoleController interface {
	CreateRole(ctx *echo.Context)          // 创建角色
	UpdateRole(ctx *echo.Context)          // 更新角色
	DeleteRole(ctx *echo.Context)          // 删除角色
	GetRole(ctx *echo.Context)             // 获取角色详情
	PageRole(ctx *echo.Context)            // 分页查询角色列表
	AssignRoleToUser(ctx *echo.Context)    // 为用户分配角色
	RemoveRoleFromUser(ctx *echo.Context)  // 移除用户的角色
	GetUserRoles(ctx *echo.Context)        // 获取用户的所有角色
	GetRoleUsers(ctx *echo.Context)        // 获取角色下的用户
	AssignUsersToRole(ctx *echo.Context)   // 批量为角色添加用户
	RemoveUsersFromRole(ctx *echo.Context) // 批量移除角色下的用户
	GetRoleMenus(ctx *echo.Context)        // 获取角色菜单
	AssignRoleMenus(ctx *echo.Context)     // 分配角色菜单
}

type roleController struct {
	roleService iam.RoleService
	logger      logger.Logger
}

// NewRoleController 创建组件实例。
func NewRoleController(c container.Container) RoleController {
	return &roleController{
		roleService: iam.NewRoleService(c.GetDB(), c.GetLogger()),
		logger:      c.GetLogger(),
	}
}

// CreateRole 创建角色
//
//	@Summary		创建角色
//	@Description	创建新角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {token}"
//	@Param			body			body		request.CreateRoleRequest	true	"角色信息"
//	@Success		200				{object}	response.Response{data=model.Role}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role [post]
func (c *roleController) CreateRole(ctx *echo.Context) {
	var req request.CreateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	userId := ctx.Get("userId")

	role := &model.Role{
		RoleKey:   req.RoleKey,
		RoleName:  req.RoleName,
		Sort:      req.Sort,
		Status:    req.Status,
		DataScope: req.DataScope,
		IsSystem:  false,
		Remark:    req.Remark,
	}
	role.CreateBy = userId.(int64)

	if err := c.roleService.Create(ctx.Request().Context(), role); err != nil {
		c.logger.Error("创建角色失败", zap.Error(err))
		response.InternalServerError(ctx, "创建角色失败: "+err.Error())
		return
	}

	response.Success(ctx, role)
}

// UpdateRole 更新角色
//
//	@Summary		更新角色
//	@Description	更新角色信息
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {token}"
//	@Param			roleId			path		int							true	"角色ID"
//	@Param			body			body		request.UpdateRoleRequest	true	"角色信息"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role/{roleId} [put]
func (c *roleController) UpdateRole(ctx *echo.Context) {
	roleIdStr := ctx.Param("roleId")
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	var req request.UpdateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	req.RoleId = roleId

	userId := ctx.Get("userId")

	role := &model.Role{
		ID:        req.RoleId,
		RoleName:  req.RoleName,
		Sort:      req.Sort,
		Status:    req.Status,
		DataScope: req.DataScope,
		Remark:    req.Remark,
	}
	role.UpdateBy = userId.(int64)

	if err := c.roleService.Update(ctx.Request().Context(), role); err != nil {
		c.logger.Error("更新角色失败", zap.Error(err))
		response.InternalServerError(ctx, "更新角色失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

// DeleteRole 删除角色
//
//	@Summary		删除角色
//	@Description	删除角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			roleId			path		int		true	"角色ID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role/{roleId} [delete]
func (c *roleController) DeleteRole(ctx *echo.Context) {
	roleIdStr := ctx.Param("roleId")
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	if err := c.roleService.Delete(ctx.Request().Context(), roleId); err != nil {
		c.logger.Error("删除角色失败", zap.Error(err))
		response.InternalServerError(ctx, "删除角色失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetRole 获取角色详情
//
//	@Summary		获取角色详情
//	@Description	根据角色ID获取角色详情
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			roleId			path		int		true	"角色ID"
//	@Success		200				{object}	response.Response{data=model.Role}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role/{roleId} [get]
func (c *roleController) GetRole(ctx *echo.Context) {
	roleIdStr := ctx.Param("roleId")
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	role, err := c.roleService.GetById(ctx.Request().Context(), roleId)
	if err != nil {
		c.logger.Error("获取角色详情失败", zap.Error(err))
		response.InternalServerError(ctx, "获取角色详情失败: "+err.Error())
		return
	}

	response.Success(ctx, role)
}

// PageRole 分页查询角色列表
//
//	@Summary		分页查询角色列表
//	@Description	使用 Paginator 分页查询角色列表
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer {token}"
//	@Param			body			body		request.PageRoleRequest	true	"查询参数"
//	@Success		200				{object}	response.Response{data=object}
//	@Router			/api/v1/role/page [post]
func (c *roleController) PageRole(ctx *echo.Context) {
	var req request.PageRoleRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	page, err := c.roleService.Page(ctx.Request().Context(), req.PageNum, req.PageSize, req.RoleName, req.Status)
	if err != nil {
		c.logger.Error("分页查询角色列表失败", zap.Error(err))
		response.InternalServerError(ctx, "分页查询角色列表失败: "+err.Error())
		return
	}

	response.Success(ctx, page)
}

// AssignRoleToUser 为用户分配角色
//
//	@Summary		为用户分配角色
//	@Description	为用户分配角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			body			body		request.AssignRoleToUserRequest	true	"分配信息"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role/assign [post]
func (c *roleController) AssignRoleToUser(ctx *echo.Context) {
	var req request.AssignRoleToUserRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if err := c.roleService.AssignRoleToUser(ctx.Request().Context(), req.UserId.Int64(), req.RoleId.Int64()); err != nil {
		c.logger.Error("为用户分配角色失败", zap.Error(err))
		response.InternalServerError(ctx, "为用户分配角色失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

// RemoveRoleFromUser 移除用户的角色
//
//	@Summary		移除用户的角色
//	@Description	移除指定用户的角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			userId			query		int		true	"用户ID"
//	@Param			roleId			query		int		true	"角色ID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role/remove [delete]
func (c *roleController) RemoveRoleFromUser(ctx *echo.Context) {
	userIdStr := ctx.QueryParam("userId")
	roleIdStr := ctx.QueryParam("roleId")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户ID格式错误")
		return
	}
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	if err := c.roleService.RemoveRoleFromUser(ctx.Request().Context(), userId, roleId); err != nil {
		c.logger.Error("移除用户角色失败", zap.Error(err))
		response.InternalServerError(ctx, "移除用户角色失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetUserRoles 获取用户的所有角色
//
//	@Summary		获取用户的所有角色
//	@Description	获取用户的所有角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			userId			query		int		true	"用户ID"
//	@Success		200				{object}	response.Response{data=[]model.Role}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/role/user [get]
func (c *roleController) GetUserRoles(ctx *echo.Context) {
	userIdStr := ctx.QueryParam("userId")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户ID格式错误")
		return
	}

	roles, err := c.roleService.GetUserRoles(ctx.Request().Context(), userId)
	if err != nil {
		c.logger.Error("获取用户角色失败", zap.Error(err))
		response.InternalServerError(ctx, "获取用户角色失败: "+err.Error())
		return
	}

	response.Success(ctx, roles)
}

// GetRoleUsers 获取角色下的用户
func (c *roleController) GetRoleUsers(ctx *echo.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	users, err := c.roleService.GetRoleUsers(ctx.Request().Context(), roleId)
	if err != nil {
		c.logger.Error("获取角色用户失败", zap.Error(err))
		response.InternalServerError(ctx, "获取角色用户失败: "+err.Error())
		return
	}

	response.Success(ctx, users)
}

// AssignUsersToRole 批量为角色添加用户
func (c *roleController) AssignUsersToRole(ctx *echo.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	var req request.BatchRoleUsersRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if err := c.roleService.AssignUsersToRole(ctx.Request().Context(), roleId, toInt64IDs(req.UserIds), currentUserId(ctx)); err != nil {
		c.logger.Error("批量为角色添加用户失败", zap.Error(err))
		response.InternalServerError(ctx, "批量为角色添加用户失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

// RemoveUsersFromRole 批量移除角色下的用户
func (c *roleController) RemoveUsersFromRole(ctx *echo.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}

	var req request.BatchRoleUsersRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if err := c.roleService.RemoveUsersFromRole(ctx.Request().Context(), roleId, toInt64IDs(req.UserIds)); err != nil {
		c.logger.Error("批量移除角色用户失败", zap.Error(err))
		response.InternalServerError(ctx, "批量移除角色用户失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

func toInt64IDs(ids []request.Int64ID) []int64 {
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		result = append(result, id.Int64())
	}
	return result
}

// GetRoleMenus 获取角色菜单 ID 列表
func (c *roleController) GetRoleMenus(ctx *echo.Context) {
	roleIdStr := ctx.Param("roleId")
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}
	menus, err := c.roleService.GetRoleMenus(ctx.Request().Context(), roleId)
	if err != nil {
		c.logger.Error("获取角色菜单失败", zap.Error(err))
		response.InternalServerError(ctx, "获取角色菜单失败: "+err.Error())
		return
	}
	menuIds := make([]int64, 0, len(menus))
	for _, menu := range menus {
		menuIds = append(menuIds, menu.ID)
	}
	response.Success(ctx, menuIds)
}

// AssignRoleMenus 分配角色菜单
func (c *roleController) AssignRoleMenus(ctx *echo.Context) {
	roleIdStr := ctx.Param("roleId")
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "角色ID格式错误")
		return
	}
	var req request.AssignRoleMenusRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	if err := c.roleService.AssignMenusToRole(ctx.Request().Context(), roleId, toInt64IDs(req.MenuIds)); err != nil {
		c.logger.Error("分配角色菜单失败", zap.Error(err))
		response.InternalServerError(ctx, "分配角色菜单失败: "+err.Error())
		return
	}
	response.Success(ctx, nil)
}
