package server

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func registerAttachmentHTTPRoutes(srv *khttp.Server, deps *GatewayDeps, svc v1.AttachmentServiceHTTPServer) {
	if srv == nil || deps == nil || deps.Attachment == nil || svc == nil {
		return
	}
	r := srv.Route("/")
	r.POST("/api/v1/attachment/upload-file", wrapOperation(v1.OperationAttachmentServiceUploadFile, func(ctx context.Context, httpCtx khttp.Context) error {
		return uploadAttachmentFile(ctx, httpCtx, deps.Attachment)
	}))
	r.POST("/api/v1/attachment/{attachmentId}/bind", wrapOperation(v1.OperationAttachmentServiceBindAttachmentToBusiness, func(ctx context.Context, httpCtx khttp.Context) error {
		var in v1.BindAttachmentRequest
		if err := httpCtx.Bind(&in); err != nil {
			return err
		}
		if err := httpCtx.BindQuery(&in); err != nil {
			return err
		}
		if err := httpCtx.BindVars(&in); err != nil {
			return err
		}
		reply, err := svc.BindAttachmentToBusiness(ctx, &in)
		if err != nil {
			return err
		}
		return httpCtx.Result(http.StatusOK, reply)
	}))
	r.GET("/api/v1/attachment/{attachmentId}", wrapOperation(v1.OperationAttachmentServiceGetAttachment, func(ctx context.Context, httpCtx khttp.Context) error {
		var in v1.AttachmentIdRequest
		if err := httpCtx.BindQuery(&in); err != nil {
			return err
		}
		if err := httpCtx.BindVars(&in); err != nil {
			return err
		}
		reply, err := svc.GetAttachment(ctx, &in)
		if err != nil {
			return err
		}
		return httpCtx.Result(http.StatusOK, reply)
	}))
	r.GET("/api/v1/attachment/business", wrapOperation(v1.OperationAttachmentServiceListAttachmentsByBusiness, func(ctx context.Context, httpCtx khttp.Context) error {
		var in v1.ListAttachmentsByBusinessRequest
		if err := httpCtx.BindQuery(&in); err != nil {
			return err
		}
		reply, err := svc.ListAttachmentsByBusiness(ctx, &in)
		if err != nil {
			return err
		}
		return httpCtx.Result(http.StatusOK, reply)
	}))
	r.POST("/api/v1/attachment/page", wrapOperation(v1.OperationAttachmentServicePageAttachments, func(ctx context.Context, httpCtx khttp.Context) error {
		var in v1.PageAttachmentsRequest
		if err := httpCtx.Bind(&in); err != nil {
			return err
		}
		if err := httpCtx.BindQuery(&in); err != nil {
			return err
		}
		reply, err := svc.PageAttachments(ctx, &in)
		if err != nil {
			return err
		}
		return httpCtx.Result(http.StatusOK, reply)
	}))
	r.GET("/api/v1/attachment/{attachmentId}/download", wrapOperation(v1.OperationAttachmentServiceDownloadAttachment, func(ctx context.Context, httpCtx khttp.Context) error {
		return downloadAttachmentFile(ctx, httpCtx, deps.Attachment)
	}))
	r.GET("/api/v1/attachment/{attachmentId}/url", wrapOperation(v1.OperationAttachmentServiceGetAttachmentURL, func(ctx context.Context, httpCtx khttp.Context) error {
		var in v1.AttachmentURLRequest
		if err := httpCtx.BindQuery(&in); err != nil {
			return err
		}
		if err := httpCtx.BindVars(&in); err != nil {
			return err
		}
		reply, err := svc.GetAttachmentURL(ctx, &in)
		if err != nil {
			return err
		}
		return httpCtx.Result(http.StatusOK, reply)
	}))
	r.DELETE("/api/v1/attachment/{attachmentId}", wrapOperation(v1.OperationAttachmentServiceDeleteAttachment, func(ctx context.Context, httpCtx khttp.Context) error {
		var in v1.AttachmentIdRequest
		if err := httpCtx.BindQuery(&in); err != nil {
			return err
		}
		if err := httpCtx.BindVars(&in); err != nil {
			return err
		}
		reply, err := svc.DeleteAttachment(ctx, &in)
		if err != nil {
			return err
		}
		return httpCtx.Result(http.StatusOK, reply)
	}))
}

func wrapOperation(operation string, fn func(ctx context.Context, httpCtx khttp.Context) error) func(ctx khttp.Context) error {
	return func(httpCtx khttp.Context) error {
		khttp.SetOperation(httpCtx, operation)
		h := httpCtx.Middleware(func(ctx context.Context, req any) (any, error) {
			return nil, fn(ctx, httpCtx)
		})
		_, err := h(httpCtx, nil)
		return err
	}
}

func uploadAttachmentFile(ctx context.Context, httpCtx khttp.Context, gateway AttachmentGateway) error {
	req := httpCtx.Request()
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		return kerrors.BadRequest("BAD_REQUEST", "文件上传参数错误")
	}
	fileHeader, err := extractMultipartFile(req, "file")
	if err != nil {
		return err
	}
	content, err := readUploadedContent(fileHeader)
	if err != nil {
		return err
	}
	item, err := gateway.Upload(ctx, &v1.UploadFileRequest{
		FileName: fileHeader.Filename,
		FileType: strings.TrimPrefix(filepath.Ext(fileHeader.Filename), "."),
		Content:  content,
		EnvCode:  strings.TrimSpace(req.FormValue("envCode")),
	})
	if err != nil {
		return err
	}
	return httpCtx.Result(http.StatusOK, item)
}

func downloadAttachmentFile(ctx context.Context, httpCtx khttp.Context, gateway AttachmentGateway) error {
	attachmentID, err := parsePathInt64(httpCtx, "attachmentId")
	if err != nil {
		return err
	}
	item, err := gateway.Download(ctx, attachmentID)
	if err != nil {
		return err
	}
	if item == nil || len(item.GetContent()) == 0 {
		return kerrors.NotFound("NOT_FOUND", "附件不存在")
	}
	fileName := strings.TrimSpace(item.GetFileName())
	if fileName == "" {
		fileName = fmt.Sprintf("attachment-%d", attachmentID)
	}
	httpCtx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	httpCtx.Response().Header().Set("Content-Type", "application/octet-stream")
	_, err = httpCtx.Response().Write(item.GetContent())
	return err
}

func extractMultipartFile(req *http.Request, field string) (*multipart.FileHeader, error) {
	file, header, err := req.FormFile(field)
	if err != nil {
		return nil, kerrors.BadRequest("BAD_REQUEST", "请选择要上传的文件")
	}
	_ = file.Close()
	return header, nil
}

func readUploadedContent(header *multipart.FileHeader) ([]byte, error) {
	file, err := header.Open()
	if err != nil {
		return nil, kerrors.BadRequest("BAD_REQUEST", "读取上传文件失败")
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, kerrors.BadRequest("BAD_REQUEST", "读取上传文件失败")
	}
	if len(content) == 0 {
		return nil, kerrors.BadRequest("BAD_REQUEST", "上传文件不能为空")
	}
	return content, nil
}

func parsePathInt64(ctx khttp.Context, name string) (int64, error) {
	value := strings.TrimSpace(ctx.Vars().Get(name))
	if value == "" {
		return 0, kerrors.BadRequest("BAD_REQUEST", "缺少路径参数")
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, kerrors.BadRequest("BAD_REQUEST", "路径参数格式错误")
	}
	return parsed, nil
}
