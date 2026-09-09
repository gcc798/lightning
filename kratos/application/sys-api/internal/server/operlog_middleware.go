package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/data"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func operLogMiddleware(deps *GatewayDeps) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if deps == nil || deps.OperLog == nil {
				return next(ctx, req)
			}
			request, _ := khttp.RequestFromServerContext(ctx)
			operation := transportOperation(ctx)
			if shouldSkipOperLog(operation, request) {
				return next(ctx, req)
			}
			bodySummary := captureRequestSummary(request)
			startedAt := time.Now()
			resp, err := next(ctx, req)
			entry := buildOperLogEntry(ctx, request, operation, bodySummary, resp, err, startedAt)
			if entry != nil {
				_, _ = deps.OperLog.Create(ctx, entry)
			}
			return resp, err
		}
	}
}

func shouldSkipOperLog(operation string, req *http.Request) bool {
	if _, ok := publicOperations[operation]; ok {
		return true
	}
	path := ""
	if req != nil && req.URL != nil {
		path = req.URL.Path
	}
	for _, item := range []string{"/health", "/health/ready", "/health/live", "/health/startup", "/metrics", "/swagger", "/ws", "/api/v1/loginLog", "/api/v1/operLog"} {
		if path == item || strings.HasPrefix(path, item+"/") {
			return true
		}
	}
	return false
}

func captureRequestSummary(req *http.Request) string {
	if req == nil {
		return ""
	}
	contentType := strings.ToLower(strings.TrimSpace(req.Header.Get("Content-Type")))
	if strings.Contains(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			return "multipart/form-data"
		}
		parts := make([]string, 0)
		if req.MultipartForm != nil {
			for key, values := range req.MultipartForm.Value {
				for _, value := range values {
					parts = append(parts, fmt.Sprintf("%s=%s", key, value))
				}
			}
			for key, files := range req.MultipartForm.File {
				for _, file := range files {
					parts = append(parts, fmt.Sprintf("%s=[file:%s,size:%d]", key, file.Filename, file.Size))
				}
			}
		}
		return strings.Join(parts, ",")
	}
	if req.Body == nil || req.Method == http.MethodGet || req.Method == http.MethodDelete {
		if req.URL == nil {
			return ""
		}
		return req.URL.RawQuery
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func buildOperLogEntry(ctx context.Context, req *http.Request, operation string, requestSummary string, resp any, err error, startedAt time.Time) *v1.CreateOperLogRequest {
	if req == nil {
		return nil
	}
	resource, action := inferOperationResource(operation, req)
	status := "0"
	message := ""
	if err != nil {
		status = "1"
		message = normalizeOperLogError(err)
	}
	return &v1.CreateOperLogRequest{
		Title:         resource,
		OperName:      buildOperLogUserName(ctx),
		OperIp:        data.CurrentClientIP(ctx),
		Status:        status,
		ErrorMsg:      message,
		BusinessType:  action,
		Method:        operation,
		RequestMethod: req.Method,
		DeviceType:    inferDeviceType(req.UserAgent()),
		OperUrl:       req.URL.Path,
		OperLocation:  "",
		OperParam:     truncateOperLogValue(requestSummary, 2000),
		JsonResult:    truncateOperLogValue(marshalOperLogResponse(resp), 2000),
		CostTime:      time.Since(startedAt).Milliseconds(),
		UserAgent:     strings.TrimSpace(req.UserAgent()),
	}
}

func inferOperationResource(operation string, req *http.Request) (string, string) {
	method := "unknown"
	if req != nil {
		method = strings.ToLower(req.Method)
	}
	if permission, ok := permissionForOperation(operation); ok {
		parts := strings.SplitN(permission, ".", 2)
		resource := parts[0]
		action := "operate"
		if len(parts) == 2 {
			action = parts[1]
		}
		return resource, action
	}
	path := ""
	if req != nil && req.URL != nil {
		path = req.URL.Path
	}
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) >= 3 && segments[0] == "api" && segments[1] == "v1" {
		return strings.ReplaceAll(segments[2], "-", "_"), method
	}
	if len(segments) > 0 && segments[0] != "" {
		return strings.ReplaceAll(segments[0], "-", "_"), method
	}
	return "system", method
}

func buildOperLogUserName(ctx context.Context) string {
	userID := data.CurrentUserID(ctx)
	userName := data.CurrentUserName(ctx)
	if userID > 0 && strings.TrimSpace(userName) != "" {
		return strconv.FormatInt(userID, 10) + "-" + strings.TrimSpace(userName)
	}
	if userID > 0 {
		return strconv.FormatInt(userID, 10)
	}
	if strings.TrimSpace(userName) != "" {
		return strings.TrimSpace(userName)
	}
	return "anonymous"
}

func normalizeOperLogError(err error) string {
	if err == nil {
		return ""
	}
	se := kerrors.FromError(err)
	message := strings.TrimSpace(se.Message)
	if message == "" {
		message = strings.TrimSpace(err.Error())
	}
	return message
}

func marshalOperLogResponse(resp any) string {
	if resp == nil {
		return ""
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		return ""
	}
	return string(payload)
}

func truncateOperLogValue(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func inferDeviceType(userAgent string) string {
	lower := strings.ToLower(strings.TrimSpace(userAgent))
	switch {
	case strings.Contains(lower, "micromessenger"):
		return "wechat"
	case strings.Contains(lower, "iphone"), strings.Contains(lower, "ios"):
		return "ios"
	case strings.Contains(lower, "android"):
		return "android"
	default:
		return "web"
	}
}
