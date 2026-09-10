package controller

import (
	"strconv"

	iam "github.com/gcc798/microservice-kit/application/iam/internal/domain"
	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	"github.com/gcc798/microservice-kit/application/iam/internal/response"
	"github.com/labstack/echo/v5"
)

// MenuController 菜单控制器接口
type MenuController interface {
	GetUserMenuTree(ctx *echo.Context) // 获取当前用户的菜单树
	GetMenuTree(ctx *echo.Context)     // 获取所有菜单树
	GetMenuList(ctx *echo.Context)     // 获取菜单列表
	GetMenuById(ctx *echo.Context)     // 获取菜单详情
	CreateMenu(ctx *echo.Context)      // 创建菜单
	UpdateMenu(ctx *echo.Context)      // 更新菜单
	DeleteMenu(ctx *echo.Context)      // 删除菜单
}

type menuController struct {
	menuService *iam.MenuService
}

// NewMenuController 创建组件实例。
func NewMenuController(menuService *iam.MenuService) MenuController {
	return &menuController{
		menuService: menuService,
	}
}

// GetUserMenuTree 获取当前用户的菜单树
//
//	@Summary		获取用户菜单树
//	@Description	获取当前登录用户的菜单树，用于前端生成动态路由
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	response.Response{data=[]iam.MenuTree}	"成功"
//	@Failure		401	{object}	response.Response							"未授权"
//	@Failure		500	{object}	response.Response							"服务器错误"
//	@Router			/api/v1/menu/user/tree [get]
func (c *menuController) GetUserMenuTree(ctx *echo.Context) {
	userId := ctx.Get("userId")
	if userId == nil {
		response.FailCode(ctx, response.CodeUnauthorized, "未授权")
		return
	}

	userIdInt64, ok := userId.(int64)
	if !ok {
		response.FailCode(ctx, response.CodeServerError, "用户ID类型错误")
		return
	}

	tree, err := c.menuService.GetUserMenuTree(userIdInt64)
	if err != nil {
		response.FailCode(ctx, response.CodeServerError, err.Error())
		return
	}

	response.Success(ctx, tree)
}

// GetMenuTree 获取所有菜单树
//
//	@Summary		获取菜单树
//	@Description	获取所有菜单的树形结构，用于菜单管理页面
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	response.Response{data=[]iam.MenuTree}	"成功"
//	@Failure		500	{object}	response.Response							"服务器错误"
//	@Router			/api/v1/menu/tree [get]
func (c *menuController) GetMenuTree(ctx *echo.Context) {
	tree, err := c.menuService.GetAllMenuTree()
	if err != nil {
		response.FailCode(ctx, response.CodeServerError, err.Error())
		return
	}

	response.Success(ctx, tree)
}

// GetMenuList 获取菜单列表
//
//	@Summary		获取菜单列表
//	@Description	获取所有菜单的列表形式
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	response.Response{data=[]model.Menu}	"成功"
//	@Failure		500	{object}	response.Response						"服务器错误"
//	@Router			/api/v1/menu [get]
func (c *menuController) GetMenuList(ctx *echo.Context) {
	menus, err := c.menuService.GetMenuList()
	if err != nil {
		response.FailCode(ctx, response.CodeServerError, err.Error())
		return
	}

	response.Success(ctx, menus)
}

// GetMenuById 获取菜单详情
//
//	@Summary		获取菜单详情
//	@Description	根据菜单ID获取菜单详细信息
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		int									true	"菜单ID"
//	@Success		200	{object}	response.Response{data=model.Menu}	"成功"
//	@Failure		400	{object}	response.Response					"参数错误"
//	@Failure		404	{object}	response.Response					"菜单不存在"
//	@Failure		500	{object}	response.Response					"服务器错误"
//	@Router			/api/v1/menu/{id} [get]
func (c *menuController) GetMenuById(ctx *echo.Context) {
	menuId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.FailCode(ctx, response.CodeInvalidParam, "无效的菜单ID")
		return
	}

	menu, err := c.menuService.GetMenuById(menuId)
	if err != nil {
		response.FailCode(ctx, response.CodeNotFound, "菜单不存在")
		return
	}

	response.Success(ctx, menu)
}

// CreateMenu 创建菜单
//
//	@Summary		创建菜单
//	@Description	创建新的菜单项
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			menu	body		model.Menu			true	"菜单信息"
//	@Success		200		{object}	response.Response	"成功"
//	@Failure		400		{object}	response.Response	"参数错误"
//	@Failure		500		{object}	response.Response	"服务器错误"
//	@Router			/api/v1/menu [post]
func (c *menuController) CreateMenu(ctx *echo.Context) {
	var menu model.Menu
	if err := ctx.Bind(&menu); err != nil {
		response.FailCode(ctx, response.CodeInvalidParam, "参数错误: "+err.Error())
		return
	}

	userId := ctx.Get("userId")
	menu.CreateBy = userId.(int64)
	menu.UpdateBy = userId.(int64)

	if err := c.menuService.CreateMenu(&menu); err != nil {
		response.FailCode(ctx, response.CodeServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// UpdateMenu 更新菜单
//
//	@Summary		更新菜单
//	@Description	更新菜单信息
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		int					true	"菜单ID"
//	@Param			menu	body		model.Menu			true	"菜单信息"
//	@Success		200		{object}	response.Response	"成功"
//	@Failure		400		{object}	response.Response	"参数错误"
//	@Failure		500		{object}	response.Response	"服务器错误"
//	@Router			/api/v1/menu/{id} [put]
func (c *menuController) UpdateMenu(ctx *echo.Context) {
	menuId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.FailCode(ctx, response.CodeInvalidParam, "无效的菜单ID")
		return
	}

	var menu model.Menu
	if err := ctx.Bind(&menu); err != nil {
		response.FailCode(ctx, response.CodeInvalidParam, "参数错误: "+err.Error())
		return
	}

	menu.ID = menuId

	userId := ctx.Get("userId")
	menu.UpdateBy = userId.(int64)

	if err := c.menuService.UpdateMenu(&menu); err != nil {
		response.FailCode(ctx, response.CodeServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// DeleteMenu 删除菜单
//
//	@Summary		删除菜单
//	@Description	删除指定的菜单
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		int					true	"菜单ID"
//	@Success		200	{object}	response.Response	"成功"
//	@Failure		400	{object}	response.Response	"参数错误"
//	@Failure		500	{object}	response.Response	"服务器错误"
//	@Router			/api/v1/menu/{id} [delete]
func (c *menuController) DeleteMenu(ctx *echo.Context) {
	menuId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.FailCode(ctx, response.CodeInvalidParam, "无效的菜单ID")
		return
	}

	if err := c.menuService.DeleteMenu(menuId); err != nil {
		response.FailCode(ctx, response.CodeServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}
