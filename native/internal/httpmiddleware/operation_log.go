package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/utils"
	"github.com/gcc798/microservice-kit/internal/utils/idgen"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// 操作日志上下文键
const (
	// OperLogTitleKey 定义业务常量。
	OperLogTitleKey = "oper_log_title"
	// OperLogBusinessTypeKey 定义业务常量。
	OperLogBusinessTypeKey = "oper_log_business_type"
)

// 操作日志批量写入配置
const (
	operLogBatchSize     = 100              // 批量写入大小
	operLogFlushInterval = 10 * time.Second // 定时刷新间隔
	operLogChannelSize   = 1000             // 通道缓冲大小
	operLogBodyLimit     = 64 << 10         // 最多记录 64 KiB 请求体
	redactedValue        = "***"
)

// OperLogWriter 操作日志写入器
type OperLogWriter struct {
	api      sysv1.API
	logger   logging.Logger
	logChan  chan *sysv1.OperationLog
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	stopOnce sync.Once
}

// NewOperLogWriter creates and starts one application-scoped operation log writer.
func NewOperLogWriter(api sysv1.API, logger logging.Logger) *OperLogWriter {
	ctx, cancel := context.WithCancel(context.Background())
	w := &OperLogWriter{
		api:     api,
		logger:  logger,
		logChan: make(chan *sysv1.OperationLog, operLogChannelSize),
		ctx:     ctx,
		cancel:  cancel,
	}
	w.start()
	return w
}

// start 启动日志写入协程
func (w *OperLogWriter) start() {
	w.wg.Add(1)
	go w.batchWrite()
}

// batchWrite 批量写入日志
func (w *OperLogWriter) batchWrite() {
	defer w.wg.Done()

	buffer := make([]*sysv1.OperationLog, 0, operLogBatchSize)
	ticker := time.NewTicker(operLogFlushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(buffer) == 0 {
			return
		}

		// 批量插入
		if err := w.api.RecordOperations(w.ctx, &sysv1.RecordOperationsRequest{Logs: buffer}); err != nil {
			w.logger.Error("批量写入操作日志失败",
				zap.Error(err),
				zap.Int("count", len(buffer)))
		} else {
			w.logger.Debug("批量写入操作日志成功",
				zap.Int("count", len(buffer)))
		}

		// 清空缓冲区
		buffer = buffer[:0]
	}

	for {
		select {
		case <-w.ctx.Done():
			// 程序退出前刷新剩余日志
			flush()
			w.logger.Info("操作日志写入器已停止")
			return

		case log := <-w.logChan:
			buffer = append(buffer, log)
			// 达到批量大小，立即刷新
			if len(buffer) >= operLogBatchSize {
				flush()
			}

		case <-ticker.C:
			// 定时刷新
			flush()
		}
	}
}

// Write 写入日志到通道
func (w *OperLogWriter) Write(log *sysv1.OperationLog) {
	select {
	case w.logChan <- log:
		// 成功写入通道
	default:
		// 通道已满，记录警告
		w.logger.Warn("操作日志通道已满，丢弃日志",
			zap.String("title", log.Title),
			zap.String("operName", log.OperatorName))
	}
}

// Stop 停止日志写入器
func (w *OperLogWriter) Stop() {
	w.stopOnce.Do(func() {
		w.cancel()
		w.wg.Wait()
	})
}

// OperationLog 操作日志中间件
func OperationLog(writer *OperLogWriter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// 跳过操作日志相关的请求，避免递归记录
			if shouldSkipLogging(c.Request().URL.Path) {
				return next(c)
			}

			start := utils.Now()
			bodyBytes := readBody(c)

			handlerErr := next(c)
			if handlerErr != nil {
				c.Set("handler_error", handlerErr)
			}

			duration := time.Since(start.Time()).Milliseconds()
			operParam := extractParams(c, bodyBytes)
			status, errMsg := resolveStatus(c)

			// 获取用户信息：格式为 "用户ID-用户名称"
			userId := getStringValue(c, "userId")
			userName := getStringValue(c, "userName")
			operName := formatOperName(userId, userName)

			// 获取终端类型（从 token 中解析）
			deviceType := getStringValue(c, "deviceType")
			if deviceType == "" {
				deviceType = "unknown"
			}

			// 获取标题和业务类型（优先从上下文获取，否则自动推断）
			title := getOperLogTitle(c)
			businessType := getBusinessType(c)

			logEntry := &sysv1.OperationLog{
				Id: idgen.MustNextID(), Title: title, BusinessType: businessType,
				Method: c.RouteInfo().Name, RequestMethod: c.Request().Method, DeviceType: deviceType,
				OperatorName: operName, Url: c.Request().URL.Path, IpAddress: utils.GetClientIP(c),
				Parameters: truncate(operParam, 2000), Status: status, ErrorMessage: truncate(errMsg, 1000),
				OperationTimeUnixMilli: start.Time().UnixMilli(), CostMillis: duration, UserAgent: c.Request().UserAgent(),
			}

			// 写入日志到通道（批量写入）
			writer.Write(logEntry)
			return handlerErr
		}
	}
}

func readBody(c *echo.Context) []byte {
	if c.Request() == nil || c.Request().Body == nil {
		return nil
	}
	// 仅对可重复读场景处理
	if c.Request().Method == http.MethodGet || c.Request().Method == http.MethodDelete {
		return nil
	}
	body, err := io.ReadAll(io.LimitReader(c.Request().Body, operLogBodyLimit+1))
	if err != nil {
		return nil
	}
	c.Request().Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), c.Request().Body))
	if len(body) > operLogBodyLimit {
		return nil
	}
	return body
}

func extractParams(c *echo.Context, body []byte) string {
	// 检查是否为文件上传请求
	contentType := c.Request().Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// 对于文件上传，记录表单字段而不是二进制内容
		return extractMultipartParams(c)
	}

	if len(body) > 0 {
		return sanitizeJSON(body)
	}
	return sanitizeQuery(c.Request().URL.Query()).Encode()
}

// extractMultipartParams 提取 multipart/form-data 请求的参数摘要
func extractMultipartParams(c *echo.Context) string {
	// 解析 multipart 表单
	if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
		return "文件上传请求（解析失败）"
	}

	var params []string

	// 记录普通表单字段
	if c.Request().MultipartForm != nil && c.Request().MultipartForm.Value != nil {
		for key, values := range c.Request().MultipartForm.Value {
			for _, value := range values {
				if isSensitiveField(key) {
					value = redactedValue
				}
				params = append(params, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	// 记录文件信息（文件名和大小，不包含内容）
	if c.Request().MultipartForm != nil && c.Request().MultipartForm.File != nil {
		for fieldName, files := range c.Request().MultipartForm.File {
			for _, file := range files {
				params = append(params, fmt.Sprintf("%s=[文件: %s, 大小: %d bytes]",
					fieldName, file.Filename, file.Size))
			}
		}
	}

	if len(params) == 0 {
		return "文件上传请求"
	}

	return strings.Join(params, ", ")
}

func resolveStatus(c *echo.Context) (string, string) {
	if value := c.Get("handler_error"); value != nil {
		return "1", fmt.Sprint(value)
	}
	status := responseStatus(c.Response())
	if status >= http.StatusBadRequest {
		return "1", http.StatusText(status)
	}
	return "0", ""
}

func responseStatus(writer http.ResponseWriter) int {
	for writer != nil {
		if response, ok := writer.(*echo.Response); ok {
			if response.Status == 0 {
				return http.StatusOK
			}
			return response.Status
		}
		unwrapper, ok := writer.(interface{ Unwrap() http.ResponseWriter })
		if !ok {
			break
		}
		writer = unwrapper.Unwrap()
	}
	return http.StatusOK
}

func sanitizeJSON(body []byte) string {
	var value any
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return "[请求体不是有效 JSON]"
	}
	redactValue(value)
	encoded, err := json.Marshal(value)
	if err != nil {
		return "[请求体脱敏失败]"
	}
	return string(encoded)
}

func redactValue(value any) {
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			if isSensitiveField(key) {
				current[key] = redactedValue
				continue
			}
			redactValue(child)
		}
	case []any:
		for _, child := range current {
			redactValue(child)
		}
	}
}

func sanitizeQuery(values url.Values) url.Values {
	clean := make(url.Values, len(values))
	for key, items := range values {
		if isSensitiveField(key) {
			clean[key] = []string{redactedValue}
			continue
		}
		clean[key] = append([]string(nil), items...)
	}
	return clean
}

func isSensitiveField(field string) bool {
	normalized := strings.NewReplacer("_", "", "-", "").Replace(strings.ToLower(field))
	switch normalized {
	case "password", "oldpassword", "newpassword", "accesstoken", "refreshtoken", "token",
		"authorization", "secret", "clientsecret", "accesskey", "accesskeyid", "accesskeysecret",
		"secretkey", "secretaccesskey", "code", "smscode":
		return true
	default:
		return false
	}
}

func getStringValue(c *echo.Context, key string) string {
	if v := c.Get(key); v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

func truncate(val string, max int) string {
	if len(val) <= max {
		return val
	}
	return val[:max]
}

// getOperLogTitle 获取操作日志标题
// 优先从上下文获取，否则根据路径自动生成
func getOperLogTitle(c *echo.Context) string {
	// 优先从上下文获取自定义标题
	if title := getStringValue(c, OperLogTitleKey); title != "" {
		return title
	}

	// 根据路径自动生成标题
	path := c.Request().URL.Path
	return generateTitleFromPath(path)
}

// getBusinessType 获取业务类型
// 优先从上下文获取，否则根据请求方法和路径自动推断
func getBusinessType(c *echo.Context) string {
	// 优先从上下文获取自定义业务类型
	if businessType := getStringValue(c, OperLogBusinessTypeKey); businessType != "" {
		return businessType
	}

	// 根据请求方法和路径自动推断业务类型
	return inferBusinessType(c)
}

// generateTitleFromPath 根据路径生成标题
func generateTitleFromPath(path string) string {
	// 移除 /api/v1 前缀
	if len(path) > 7 && path[:7] == "/api/v1" {
		path = path[7:]
	}

	// 提取资源名称（第一个路径段）
	parts := splitPath(path)
	if len(parts) == 0 {
		return "未知操作"
	}

	resource := parts[0]

	// 资源名称映射
	titleMap := map[string]string{
		"user":       "用户",
		"role":       "角色",
		"menu":       "菜单",
		"org":        "组织",
		"dict":       "字典",
		"config":     "配置",
		"loginLog":   "登录日志",
		"operLog":    "操作日志",
		"attachment": "附件",
	}

	if title, ok := titleMap[resource]; ok {
		return title
	}

	return resource
}

// inferBusinessType 推断业务类型
func inferBusinessType(c *echo.Context) string {
	method := c.Request().Method
	path := c.Request().URL.Path

	// 特殊路径判断
	if containsSegment(path, "/export") {
		return "EXPORT"
	}
	if containsSegment(path, "/import") {
		return "IMPORT"
	}
	if containsSegment(path, "/grant") || containsSegment(path, "/permission") {
		return "GRANT"
	}
	if containsSegment(path, "/clean") {
		return "CLEAN"
	}
	if containsSegment(path, "/page") || containsSegment(path, "/list") {
		return "QUERY"
	}

	// 根据 HTTP 方法判断
	switch method {
	case http.MethodGet:
		return "QUERY"
	case http.MethodPost:
		return "CREATE"
	case http.MethodPut, http.MethodPatch:
		return "UPDATE"
	case http.MethodDelete:
		return "DELETE"
	default:
		return "OTHER"
	}
}

// splitPath 分割路径
func splitPath(path string) []string {
	var parts []string
	current := ""

	for _, ch := range path {
		if ch == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// containsSegment 检查路径是否包含指定段
func containsSegment(path, segment string) bool {
	parts := splitPath(path)
	for _, part := range parts {
		if part == segment[1:] { // 移除前导斜杠
			return true
		}
	}
	return false
}

// formatOperName 格式化操作者名称为 "用户ID-用户名称"
func formatOperName(userId, userName string) string {
	if userId == "" && userName == "" {
		return "-"
	}
	if userId == "" {
		return userName
	}
	if userName == "" {
		return userId
	}
	return fmt.Sprintf("%s-%s", userId, userName)
}

// shouldSkipLogging 判断是否应该跳过日志记录
func shouldSkipLogging(path string) bool {
	if path == "/metrics" || path == "/health" || strings.HasPrefix(path, "/health/") {
		return true
	}

	// 跳过操作日志相关的接口
	skipPaths := []string{
		"/api/v1/operLog",
		"/api/v1/oper-log",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}
