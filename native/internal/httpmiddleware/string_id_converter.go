package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
)

// StringIDConverter 字符串ID转换中间件
// 用于处理前端传递的字符串类型大数ID，将其转换为int64
// 这样可以避免前端大数精度丢失问题
//
// 支持两种场景：
// 1. 路径参数中的ID（如 /api/v1/user/:id）
// 2. JSON请求体中的ID字段（如 {"id": "123456", "ids": ["1", "2", "3"]}）
func StringIDConverter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// 1. 处理路径参数中的 id
			if idStr := c.Param("id"); idStr != "" {
				// 尝试将字符串转换为 int64
				if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					// 转换成功，将 int64 值存储到上下文中
					c.Set("parsed_id", id)
				}
				// 如果转换失败，保持原样，让后续处理器处理错误
			}

			// 2. 处理 JSON 请求体中的 ID 字段
			if c.Request().Method == "POST" || c.Request().Method == "PUT" || c.Request().Method == "PATCH" {
				if strings.HasPrefix(c.Request().Header.Get("Content-Type"), "application/json") {
					// 读取原始请求体
					bodyBytes, err := io.ReadAll(c.Request().Body)
					if err != nil {
						return next(c)
					}

					// 恢复请求体供后续使用
					c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					// 如果请求体为空，直接跳过
					if len(bodyBytes) == 0 {
						return next(c)
					}

					// 尝试转换 JSON 中的字符串 ID
					convertedBody, converted := convertStringIDsInJSON(bodyBytes)
					if converted {
						// 如果有转换，使用转换后的请求体
						c.Request().Body = io.NopCloser(bytes.NewBuffer(convertedBody))
						c.Request().ContentLength = int64(len(convertedBody))
					}
				}
			}

			return next(c)
		}
	}
}

// convertStringIDsInJSON 递归转换 JSON 中的字符串 ID 为数字
// 支持的字段名：id, userId, orgId, roleId, menuId, configId, dictId, ids, userIds, roleIds, menuIds 等
func convertStringIDsInJSON(data []byte) ([]byte, bool) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return data, false
	}

	converted := convertValue(obj)
	if !converted {
		return data, false
	}

	result, err := json.Marshal(obj)
	if err != nil {
		return data, false
	}

	return result, true
}

// convertValue 递归转换值中的字符串 ID
func convertValue(v interface{}) bool {
	converted := false

	switch val := v.(type) {
	case map[string]interface{}:
		for key, value := range val {
			// 检查是否是 ID 相关字段
			if isIDField(key) {
				switch fieldVal := value.(type) {
				case string:
					// 字符串转数字
					if num, err := strconv.ParseInt(fieldVal, 10, 64); err == nil {
						val[key] = num
						converted = true
					}
				case []interface{}:
					// 字符串数组转数字数组
					if convertStringArray(fieldVal) {
						converted = true
					}
				}
			} else {
				// 递归处理嵌套对象
				if convertValue(value) {
					converted = true
				}
			}
		}

	case []interface{}:
		for _, item := range val {
			if convertValue(item) {
				converted = true
			}
		}
	}

	return converted
}

// convertStringArray 转换字符串数组为数字数组
func convertStringArray(arr []interface{}) bool {
	converted := false
	for i, item := range arr {
		if str, ok := item.(string); ok {
			if num, err := strconv.ParseInt(str, 10, 64); err == nil {
				arr[i] = num
				converted = true
			}
		}
	}
	return converted
}

// isIDField 判断字段名是否是 ID 相关字段
func isIDField(fieldName string) bool {
	idFields := []string{
		"id", "ids",
		"userId", "userIds",
		"orgId", "orgIds",
		"roleId", "roleIds",
		"menuId", "menuIds",
		"configId", "configIds",
		"dictId", "dictIds",
		"envId", "envIds",
		"parentId", "parentIds",
		"storageId", "storageIds",
		"attachmentId", "attachmentIds",
		"clientId", "clientIds",
		"createBy", "updateBy",
	}

	for _, field := range idFields {
		if fieldName == field {
			return true
		}
	}
	return false
}
