package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ParseInt64Param 解析和校验 int64 类型的路径参数
// 参数：
//   - ctx: Echo 上下文
//   - paramName: 参数名称
//   - rules: 可选的校验规则（如 "required,gt=0"）
//
// 返回：
//   - int64: 解析后的值
//   - error: 错误信息
func ParseInt64Param(ctx *echo.Context, paramName string, rules ...string) (int64, error) {
	paramStr := ctx.Param(paramName)
	if paramStr == "" {
		return 0, fmt.Errorf("%s不能为空", paramName)
	}

	value, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s格式错误", paramName)
	}

	// 如果提供了校验规则，则进行校验
	if len(rules) > 0 && rules[0] != "" {
		if err := validate.Var(value, rules[0]); err != nil {
			return 0, fmt.Errorf("%s校验失败: %v", paramName, err)
		}
	}

	return value, nil
}

// ParseIntParam 解析和校验 int 类型的路径参数
func ParseIntParam(ctx *echo.Context, paramName string, rules ...string) (int, error) {
	paramStr := ctx.Param(paramName)
	if paramStr == "" {
		return 0, fmt.Errorf("%s不能为空", paramName)
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, fmt.Errorf("%s格式错误", paramName)
	}

	// 如果提供了校验规则，则进行校验
	if len(rules) > 0 && rules[0] != "" {
		if err := validate.Var(value, rules[0]); err != nil {
			return 0, fmt.Errorf("%s校验失败: %v", paramName, err)
		}
	}

	return value, nil
}

// ParseBoolParam 解析和校验 bool 类型的参数
func ParseBoolParam(ctx *echo.Context, paramName string, defaultValue bool) bool {
	paramStr := ctx.QueryParam(paramName)
	if paramStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(paramStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// ParseInt64Query 解析和校验 int64 类型的 Query 参数
func ParseInt64Query(ctx *echo.Context, paramName string, defaultValue int64, rules ...string) (int64, error) {
	paramStr := ctx.QueryParam(paramName)
	if paramStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s格式错误", paramName)
	}

	// 如果提供了校验规则，则进行校验
	if len(rules) > 0 && rules[0] != "" {
		if err := validate.Var(value, rules[0]); err != nil {
			return 0, fmt.Errorf("%s校验失败: %v", paramName, err)
		}
	}

	return value, nil
}

// ParseIntQuery 解析和校验 int 类型的 Query 参数
func ParseIntQuery(ctx *echo.Context, paramName string, defaultValue int, rules ...string) (int, error) {
	paramStr := ctx.QueryParam(paramName)
	if paramStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, fmt.Errorf("%s格式错误", paramName)
	}

	// 如果提供了校验规则，则进行校验
	if len(rules) > 0 && rules[0] != "" {
		if err := validate.Var(value, rules[0]); err != nil {
			return 0, fmt.Errorf("%s校验失败: %v", paramName, err)
		}
	}

	return value, nil
}

// ValidateVar 使用 Validator API 校验单个变量
func ValidateVar(value interface{}, rules string) error {
	return validate.Var(value, rules)
}

// ValidateStruct 校验结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func BindJSONWithTypeCasting(ctx *echo.Context, obj interface{}) error {
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if len(bodyBytes) == 0 {
		return ctx.Bind(obj)
	}
	var data interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return err
	}
	castingFields := buildTypeCastingFieldMap(obj)
	if len(castingFields) == 0 {
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return ctx.Bind(obj)
	}
	if converted := convertJSONValueByCastingFields(data, castingFields); converted {
		convertedBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(convertedBytes))
	} else {
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return ctx.Bind(obj)
}

func buildTypeCastingFieldMap(obj interface{}) map[string]string {
	result := make(map[string]string)
	t := reflect.TypeOf(obj)
	if t == nil {
		return result
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return result
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			fieldType := f.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldType.Kind() == reflect.Struct {
				sub := buildTypeCastingFieldMap(reflect.New(fieldType).Interface())
				for k, v := range sub {
					result[k] = v
				}
			}
		}
		tag, ok := f.Tag.Lookup("typecast")
		if !ok {
			continue
		}
		mode := strings.Split(tag, ",")[0]
		switch mode {
		case "stringInt64", "toTime":
		default:
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		name := strings.Split(jsonTag, ",")[0]
		if name == "" || name == "-" {
			continue
		}
		result[name] = mode
	}
	return result
}

func convertJSONValueByCastingFields(v interface{}, fields map[string]string) bool {
	converted := false
	switch val := v.(type) {
	case map[string]interface{}:
		for key, value := range val {
			if mode, ok := fields[key]; ok {
				newValue, changed := applyTypeCasting(value, mode)
				if changed {
					val[key] = newValue
					converted = true
				}
			} else {
				if convertJSONValueByCastingFields(value, fields) {
					converted = true
				}
			}
		}
	case []interface{}:
		for _, item := range val {
			if convertJSONValueByCastingFields(item, fields) {
				converted = true
			}
		}
	}
	return converted
}

func applyTypeCasting(v interface{}, mode string) (interface{}, bool) {
	switch mode {
	case "stringInt64":
		return castStringInt64(v)
	case "toTime":
		return castToTime(v)
	default:
		return v, false
	}
}

func castStringInt64(v interface{}) (interface{}, bool) {
	switch val := v.(type) {
	case string:
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			return num, true
		}
	case []interface{}:
		changed := false
		for i, item := range val {
			if s, ok := item.(string); ok {
				if num, err := strconv.ParseInt(s, 10, 64); err == nil {
					val[i] = num
					changed = true
				}
			}
		}
		return val, changed
	}
	return v, false
}

func castToTime(v interface{}) (interface{}, bool) {
	switch val := v.(type) {
	case float64:
		sec := int64(val)
		return time.Unix(sec, 0), true
	case string:
		if sec, err := strconv.ParseInt(val, 10, 64); err == nil {
			return time.Unix(sec, 0), true
		}
	case []interface{}:
		changed := false
		for i, item := range val {
			switch it := item.(type) {
			case float64:
				sec := int64(it)
				val[i] = time.Unix(sec, 0)
				changed = true
			case string:
				if sec, err := strconv.ParseInt(it, 10, 64); err == nil {
					val[i] = time.Unix(sec, 0)
					changed = true
				}
			}
		}
		return val, changed
	}
	return v, false
}
