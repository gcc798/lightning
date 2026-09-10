package controller

import (
	"net/http"
	"time"

	"github.com/gcc798/microservice-kit/application/resource/internal/domain"
	"github.com/gcc798/microservice-kit/application/resource/internal/request"
	"github.com/gcc798/microservice-kit/internal/container"
	"github.com/gcc798/microservice-kit/internal/httpresponse"
	"github.com/gcc798/microservice-kit/internal/httputils"
	"github.com/gcc798/microservice-kit/internal/logger"
	_ "github.com/gcc798/microservice-kit/internal/utils/pagination"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// AttachmentController 定义业务数据结构。
type AttachmentController interface {
	UploadFile(ctx *echo.Context)                // 上传文件
	BindAttachmentToBusiness(ctx *echo.Context)  // 绑定附件到业务
	DownloadAttachment(ctx *echo.Context)        // 下载附件
	DeleteAttachment(ctx *echo.Context)          // 删除附件
	GetAttachmentURL(ctx *echo.Context)          // 获取附件访问URL
	GetAttachment(ctx *echo.Context)             // 获取附件详情
	ListAttachmentsByBusiness(ctx *echo.Context) // 根据业务查询附件列表
	PageAttachments(ctx *echo.Context)           // 分页查询附件列表
}

type attachmentController struct {
	attachmentService resource.AttachmentService
	logger            logger.Logger
}

// NewAttachmentController 创建组件实例。
func NewAttachmentController(c container.Container) AttachmentController {
	return &attachmentController{
		attachmentService: resource.NewAttachmentService(c.GetDB(), c.GetStorage(), c.GetLogger()),
		logger:            c.GetLogger(),
	}
}

// UploadFile 上传文件（步骤1：只上传文件）
//
//	@Summary		上传文件
//	@Description	上传文件到指定存储环境（步骤1：只上传文件，返回附件ID）
//	@Tags			附件管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			file			formData	file	true	"文件"
//	@Success		200				{object}	response.Response{data=object}
//	@Router			/api/v1/attachment/upload-file [post]
func (c *attachmentController) UploadFile(ctx *echo.Context) {
	header, err := ctx.FormFile("file")
	if err != nil {
		response.BadRequest(ctx, "请选择要上传的文件")
		return
	}

	file, err := header.Open()
	if err != nil {
		response.BadRequest(ctx, "打开文件失败: "+err.Error())
		return
	}
	defer file.Close()

	attachment, err := c.attachmentService.UploadFile(ctx.Request().Context(),
		file, header.Filename, header.Header.Get("Content-Type"), header.Size)
	if err != nil {
		c.logger.Error("上传文件失败", zap.Error(err))
		response.InternalServerError(ctx, "上传文件失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "上传文件成功", attachment)
}

// BindAttachmentToBusiness 绑定附件到业务（步骤2：绑定业务信息）
//
//	@Summary		绑定附件到业务
//	@Description	将附件绑定到指定业务（步骤2：绑定业务信息）
//	@Tags			附件管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {token}"
//	@Param			attachmentId	path		int										true	"附件ID"
//	@Param			body			body		request.BindAttachmentToBusinessRequest	true	"业务信息"
//	@Success		200				{object}	response.Response
//	@Router			/api/v1/attachment/{attachmentId}/bind [post]
func (c *attachmentController) BindAttachmentToBusiness(ctx *echo.Context) {
	attachmentId, err := utils.ParseInt64Param(ctx, "attachmentId", "required")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	var req request.BindAttachmentToBusinessRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	in := resource.BindInput{
		BusinessType:  req.BusinessType,
		BusinessId:    req.BusinessId,
		BusinessField: req.BusinessField,
		IsPublic:      req.IsPublic,
		Metadata:      req.Metadata,
		ExpireTime:    req.ExpireTime,
	}
	if err := c.attachmentService.BindToBusiness(ctx.Request().Context(), attachmentId, in); err != nil {
		c.logger.Error("绑定附件到业务失败", zap.Error(err))
		response.InternalServerError(ctx, "绑定附件到业务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "绑定附件到业务成功", nil)
}

// DownloadAttachment 下载附件
//
//	@Summary		下载附件
//	@Description	下载指定附件
//	@Tags			附件管理
//	@Accept			json
//	@Produce		application/octet-stream
//	@Param			Authorization	header	string	true	"Bearer {token}"
//	@Param			attachmentId	path	int		true	"附件ID"
//	@Success		200				{file}	binary
//	@Router			/api/v1/attachment/{attachmentId}/download [get]
func (c *attachmentController) DownloadAttachment(ctx *echo.Context) {
	attachmentId, err := utils.ParseInt64Param(ctx, "attachmentId", "required")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	reader, filename, err := c.attachmentService.Download(ctx.Request().Context(), attachmentId)
	if err != nil {
		c.logger.Error("下载附件失败", zap.Error(err))
		response.InternalServerError(ctx, "下载附件失败: "+err.Error())
		return
	}
	defer reader.Close()

	ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	_ = ctx.Stream(http.StatusOK, "application/octet-stream", reader)
}

// DeleteAttachment 删除附件
//
//	@Summary		删除附件
//	@Description	删除指定附件
//	@Tags			附件管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			attachmentId	path		int		true	"附件ID"
//	@Success		200				{object}	response.Response
//	@Router			/api/v1/attachment/{attachmentId} [delete]
func (c *attachmentController) DeleteAttachment(ctx *echo.Context) {
	attachmentId, err := utils.ParseInt64Param(ctx, "attachmentId", "required")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := c.attachmentService.Delete(ctx.Request().Context(), attachmentId); err != nil {
		c.logger.Error("删除附件失败", zap.Error(err))
		response.InternalServerError(ctx, "删除附件失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除附件成功", nil)
}

// GetAttachmentURL 获取附件访问URL
//
//	@Summary		获取附件访问URL
//	@Description	获取附件的访问URL（支持临时URL）
//	@Tags			附件管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			attachmentId	path		int		true	"附件ID"
//	@Param			expires			query		int		false	"过期时间（秒）"	default(3600)
//	@Success		200				{object}	response.Response{data=object}
//	@Router			/api/v1/attachment/{attachmentId}/url [get]
func (c *attachmentController) GetAttachmentURL(ctx *echo.Context) {
	attachmentId, err := utils.ParseInt64Param(ctx, "attachmentId", "required")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	var req request.GetAttachmentURLRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	expires := time.Duration(req.Expires) * time.Second
	url, err := c.attachmentService.GetURL(ctx.Request().Context(), attachmentId, expires)
	if err != nil {
		c.logger.Error("获取附件URL失败", zap.Error(err))
		response.InternalServerError(ctx, "获取附件URL失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "获取附件URL成功", map[string]any{
		"url":     url,
		"expires": req.Expires,
	})
}

// GetAttachment 获取附件详情
//
//	@Summary		获取附件详情
//	@Description	根据附件ID获取附件详情
//	@Tags			附件管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			attachmentId	path		int		true	"附件ID"
//	@Success		200				{object}	response.Response{data=object}
//	@Router			/api/v1/attachment/{attachmentId} [get]
func (c *attachmentController) GetAttachment(ctx *echo.Context) {
	attachmentId, err := utils.ParseInt64Param(ctx, "attachmentId", "required")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	attachment, err := c.attachmentService.GetById(ctx.Request().Context(), attachmentId)
	if err != nil {
		c.logger.Error("获取附件详情失败", zap.Error(err))
		response.InternalServerError(ctx, "获取附件详情失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "获取附件详情成功", attachment)
}

// ListAttachmentsByBusiness 根据业务查询附件列表
//
//	@Summary		根据业务查询附件列表
//	@Description	根据业务类型和业务ID查询附件列表
//	@Tags			附件管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {token}"
//	@Param			businessType	query		string	true	"业务类型"
//	@Param			businessId		query		string	true	"业务ID"
//	@Success		200				{object}	response.Response{data=[]object}
//	@Router			/api/v1/attachment/business [get]
func (c *attachmentController) ListAttachmentsByBusiness(ctx *echo.Context) {
	var req request.ListAttachmentsByBusinessRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	attachments, err := c.attachmentService.ListByBusiness(ctx.Request().Context(), req.BusinessType, req.BusinessId)
	if err != nil {
		c.logger.Error("查询业务附件列表失败", zap.Error(err))
		response.InternalServerError(ctx, "查询业务附件列表失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "查询业务附件列表成功", attachments)
}

// PageAttachments 分页查询附件列表
//
//	@Summary		分页查询附件列表
//	@Description	分页查询附件列表，支持按文件名、文件类型、业务类型筛选
//	@Tags			附件管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {token}"
//	@Param			body			body		request.PageAttachmentsRequest	true	"查询参数"
//	@Success		200				{object}	response.Response{data=object}
//	@Router			/api/v1/attachment/page [post]
func (c *attachmentController) PageAttachments(ctx *echo.Context) {
	var req request.PageAttachmentsRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	page, err := c.attachmentService.Page(ctx.Request().Context(), resource.PageQuery{
		PageNum:      req.PageNum,
		PageSize:     req.PageSize,
		FileName:     req.FileName,
		FileType:     req.FileType,
		BusinessType: req.BusinessType,
	})
	if err != nil {
		c.logger.Error("分页查询附件列表失败", zap.Error(err))
		response.InternalServerError(ctx, "分页查询附件列表失败: "+err.Error())
		return
	}

	response.Success(ctx, page)
}
