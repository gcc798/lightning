package controller

import (
	"strconv"

	iam "github.com/gcc798/microservice-kit/application/iam/internal/domain"
	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	"github.com/gcc798/microservice-kit/application/iam/internal/request"
	"github.com/gcc798/microservice-kit/application/iam/internal/response"
	"github.com/gcc798/microservice-kit/internal/container"
	"github.com/gcc798/microservice-kit/internal/httputils"
	_ "github.com/gcc798/microservice-kit/internal/utils/pagination"
	"github.com/gcc798/microservice-kit/internal/validator"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// UserController 用户控制器接口
type UserController interface {
	Create(c *echo.Context)         // 创建用户
	Update(c *echo.Context)         // 更新用户
	Delete(c *echo.Context)         // 删除用户
	BatchDelete(c *echo.Context)    // 批量删除用户
	GetById(c *echo.Context)        // 根据ID查询用户
	BatchImport(c *echo.Context)    // 批量导入用户
	ResetPassword(c *echo.Context)  // 重置用户密码
	PageUser(c *echo.Context)       // 分页查询用户列表
	ChangePassword(c *echo.Context) // 用户修改密码
	XcxGetInfo(c *echo.Context)     // 获取小程序用户信息
}

type userController struct {
	ctr         container.Container
	base        *BaseController
	userService iam.UserService
}

// NewUserController 创建组件实例。
func NewUserController(c container.Container) UserController {
	return &userController{
		ctr:         c,
		base:        NewBaseController(c),
		userService: iam.NewUserService(c.GetDB(), c.GetLogger()),
	}
}

// Create 创建用户
//
//	@Summary		创建用户
//	@Description	创建新用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {token}"
//	@Param			body			body		request.CreateUserRequest	true	"用户信息"
//	@Success		200				{object}	response.Response{data=object}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user [post]
func (h *userController) Create(c *echo.Context) {
	var req request.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.TranslateWithMsg(err, &req))
		return
	}

	currentUserId, _ := h.base.GetUserId(c)
	req.CreateBy = currentUserId
	req.UpdateBy = currentUserId

	if req.Status == 0 {
		req.Status = 0
	}
	if req.UserType == 0 {
		req.UserType = 0
	}

	if err := h.userService.Create(c.Request().Context(), &req); err != nil {
		h.ctr.GetLogger().Error("创建用户失败", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, map[string]any{"userId": "ok"})
}

// Update 更新用户
//
//	@Summary		更新用户
//	@Description	更新用户信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {token}"
//	@Param			id				path		int							true	"用户ID"
//	@Param			body			body		request.UpdateUserRequest	true	"用户信息"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/{id} [put]
func (h *userController) Update(c *echo.Context) {
	userId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	var req request.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.Translate(err))
		return
	}

	req.UserId = userId
	currentUserId, _ := h.base.GetUserId(c)
	req.UpdateBy = currentUserId

	if err := h.userService.Update(c.Request().Context(), &req); err != nil {
		h.ctr.GetLogger().Error("更新用户失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// Delete 删除用户
//
//	@Summary		删除用户
//	@Description	删除指定用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			id				path		int		true	"用户ID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/{id} [delete]
func (h *userController) Delete(c *echo.Context) {
	userId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	if err := h.userService.Delete(c.Request().Context(), userId); err != nil {
		h.ctr.GetLogger().Error("删除用户失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// BatchDelete 批量删除用户
//
//	@Summary		批量删除用户
//	@Description	批量删除多个用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			body			body		request.BatchDeleteUsersRequest	true	"用户ID列表"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/batch [delete]
func (h *userController) BatchDelete(c *echo.Context) {
	var req request.BatchDeleteUsersRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.Translate(err))
		return
	}

	if err := h.userService.BatchDelete(c.Request().Context(), req.IDs); err != nil {
		h.ctr.GetLogger().Error("批量删除用户失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// GetById 根据ID查询用户
//
//	@Summary		获取用户详情
//	@Description	根据用户ID获取用户详细信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			id				path		int		true	"用户ID"
//	@Success		200				{object}	response.Response{data=object}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/{id} [get]
func (h *userController) GetById(c *echo.Context) {
	userId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	user, err := h.userService.GetById(c.Request().Context(), userId)
	if err != nil {
		h.ctr.GetLogger().Error("查询用户失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, user)
}

// BatchImport 批量导入用户
//
//	@Summary		批量导入用户
//	@Description	批量导入多个用户，返回成功和失败的统计信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			body			body		request.BatchImportUsersRequest	true	"用户列表"
//	@Success		200				{object}	response.Response{data=object}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/import [post]
func (h *userController) BatchImport(c *echo.Context) {
	var req request.BatchImportUsersRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.Translate(err))
		return
	}

	currentUserId, _ := h.base.GetUserId(c)

	for i := range req.Users {
		req.Users[i].CreateBy = currentUserId
		req.Users[i].UpdateBy = currentUserId

		if req.Users[i].Status == 0 {
			req.Users[i].Status = 0
		}
		if req.Users[i].UserType == 0 {
			req.Users[i].UserType = 0
		}
	}

	successCount, failCount, errors, err := h.userService.BatchImport(c.Request().Context(), &req)
	if err != nil {
		h.ctr.GetLogger().Error("批量导入用户失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, map[string]any{
		"successCount": successCount,
		"failCount":    failCount,
		"errors":       errors,
	})
}

// ResetPassword 重置用户密码
//
//	@Summary		重置用户密码
//	@Description	管理员重置指定用户的密码
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			id				path		int								true	"用户ID"
//	@Param			body			body		request.ResetPasswordRequest	true	"新密码"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/{id}/password [put]
func (h *userController) ResetPassword(c *echo.Context) {
	userId, err := utils.ParseInt64Param(c, "id", "required")
	if err != nil {
		response.FailCode(c, response.CodeInvalidParam, err.Error())
		return
	}

	var req request.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.Translate(err))
		return
	}

	if err := h.userService.ResetPassword(c.Request().Context(), userId, req.NewPassword); err != nil {
		h.ctr.GetLogger().Error("重置密码失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// PageUser 分页查询用户列表
//
//	@Summary		分页查询用户列表
//	@Description	使用 Paginator 分页查询用户列表
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {token}"
//	@Param			body			body		request.PageUsersRequest	true	"查询参数"
//	@Success		200				{object}	response.Response{data=object}
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/page [post]
func (h *userController) PageUser(c *echo.Context) {
	var req request.PageUsersRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.Translate(err))
		return
	}

	page, err := h.userService.Page(c.Request().Context(), req.PageNum, req.PageSize, req.UserName, req.Phonenumber, req.Status)
	if err != nil {
		h.ctr.GetLogger().Error("分页查询用户列表失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, page)
}

// ChangePassword 用户修改密码
//
//	@Summary		用户修改密码
//	@Description	用户修改自己的密码（需要验证旧密码）
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			body			body		request.ChangePasswordRequest	true	"密码信息"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Router			/api/v1/user/password/change [post]
func (h *userController) ChangePassword(c *echo.Context) {
	var req request.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		response.FailCode(c, response.CodeInvalidParam, validator.Translate(err))
		return
	}

	currentUserId, _ := h.base.GetUserId(c)

	if err := h.userService.ChangePassword(c.Request().Context(), currentUserId, req.OldPassword, req.NewPassword); err != nil {
		h.ctr.GetLogger().Error("修改密码失败", zap.Error(err))
		response.FailWithMsg(c, err.Error())
		return
	}

	response.Success(c, "ok")
}

// XcxGetInfo 获取小程序用户信息
func (h *userController) XcxGetInfo(c *echo.Context) {
	currentUserId, err := h.base.GetUserId(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	user, err := h.userService.GetById(c.Request().Context(), currentUserId)
	if err != nil {
		h.ctr.GetLogger().Error("查询用户失败", zap.Error(err))
		response.FailWithMsg(c, "没有权限访问用户数据!")
		return
	}

	info := response.XcxUserInfo{
		UserID:       user.ID,
		OrgID:        user.OrgID,
		Phonenumber:  user.Phonenumber,
		OpenID:       user.OpenId,
		UnionID:      user.UnionId,
		UserName:     user.UserName,
		NickName:     user.NickName,
		Sex:          formatUserSex(user.Sex),
		HeadPortrait: user.Avatar,
	}

	roles, err := (&model.Role{}).FindByUserId(h.ctr.GetDB(), user.ID)
	if err != nil {
		h.ctr.GetLogger().Warn("查询用户角色失败", zap.Error(err))
	} else if len(roles) > 0 {
		info.RoleKey = roles[0].RoleKey
		info.RoleName = roles[0].RoleName
	}

	response.Success(c, info)
}

func formatUserSex(sex int32) string {
	return strconv.FormatInt(int64(sex), 10)
}
