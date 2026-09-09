package router

import (
	"github.com/gcc798/lightning/application/resource/internal/controller"
	"github.com/gcc798/lightning/internal/constants"
	"github.com/gcc798/lightning/internal/httpmiddleware"
	"github.com/gcc798/lightning/internal/httpx"
)

// registerAttachmentRoutes 注册附件管理路由
func registerAttachmentRoutes(r *httpx.Router, ctx *RouterContext) {
	// 初始化 controller
	attachmentController := controller.NewAttachmentController(ctx.Container)

	// 附件管理路由组（需要认证和权限）
	attachments := r.Group("/api/v1/attachment")
	attachments.Use(ctx.AuthMiddleware)
	{
		// 分两步上传（符合参数传递规范）
		// 步骤1：上传文件 - 需要 attachment.upload 权限
		attachments.POST("/upload-file", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentUpload, "write"), attachmentController.UploadFile)
		// 步骤2：绑定业务信息 - 需要 attachment.bind 权限
		attachments.POST("/:attachmentId/bind", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentBind, "write"), attachmentController.BindAttachmentToBusiness)

		// 附件查询 - 需要 attachment.read 权限
		// 字面量路由必须先于 /:attachmentId 注册，否则会被参数路由吞掉
		attachments.GET("/business", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentRead, "read"), attachmentController.ListAttachmentsByBusiness)
		attachments.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentRead, "read"), attachmentController.PageAttachments)
		attachments.GET("/:attachmentId", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentRead, "read"), attachmentController.GetAttachment)

		// 附件下载 - 需要 attachment.download 权限
		attachments.GET("/:attachmentId/download", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentDownload, "write"), attachmentController.DownloadAttachment)

		// 获取附件URL - 需要 attachment.read 权限
		attachments.GET("/:attachmentId/url", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentRead, "read"), attachmentController.GetAttachmentURL)

		// 附件删除 - 需要 attachment.delete 权限
		attachments.DELETE("/:attachmentId", middleware.Permission(ctx.PermissionService, constants.ResourceAttachmentDelete, "write"), attachmentController.DeleteAttachment)
	}
}
