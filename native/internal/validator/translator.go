package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator

// Init 初始化中文翻译器
func Init() {
	Register()
	// 创建中文翻译器
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)

	var found bool
	trans, found = uni.GetTranslator("zh")
	if !found {
		return
	}

	// 注册翻译器
	_ = engine.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0}不能为空", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", getFieldName(fe.Field()))
		return t
	})

	_ = engine.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0}长度不能少于{1}个字符", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", getFieldName(fe.Field()), fe.Param())
		return t
	})

	_ = engine.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0}长度不能超过{1}个字符", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", getFieldName(fe.Field()), fe.Param())
		return t
	})

	_ = engine.RegisterTranslation("len", trans, func(ut ut.Translator) error {
		return ut.Add("len", "{0}长度必须为{1}位", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("len", getFieldName(fe.Field()), fe.Param())
		return t
	})

	_ = engine.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0}格式不正确", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", getFieldName(fe.Field()))
		return t
	})

	_ = engine.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0}的值必须是以下之一: {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", getFieldName(fe.Field()), fe.Param())
		return t
	})

	// 注册 json tag 获取函数，使错误信息中的字段名使用 json tag
	engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 1)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
}

// Translate 翻译错误信息
func Translate(err error) string {
	if err == nil {
		return ""
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, e.Translate(trans))
		}
		return strings.Join(messages, "；")
	}

	return err.Error()
}

// TranslateValidationError 翻译验证错误为友好的中文提示（兼容旧版本）
func TranslateValidationError(err error) string {
	return Translate(err)
}

// fieldNameMap 字段中文名映射
var fieldNameMap = map[string]string{
	"userName":    "用户名",
	"nickName":    "昵称",
	"password":    "密码",
	"email":       "邮箱",
	"phonenumber": "手机号",
	"avatar":      "头像",
	"status":      "状态",
	"remark":      "备注",
	"roleKey":     "角色标识",
	"roleName":    "角色名称",
	"sort":        "排序",
	"userType":    "用户类型",
	"sex":         "性别",
	"orgName":     "组织名称",
	"orgCode":     "组织编码",
	"orgType":     "组织类型",
	"leader":      "负责人",
	"phone":       "联系电话",
	"menuName":    "菜单名称",
	"path":        "路由地址",
	"component":   "组件路径",
	"permission":  "权限标识",
	"title":       "标题",
	"content":     "内容",
	"code":        "编码",
	"name":        "名称",
	"type":        "类型",
	"value":       "值",
	"parentId":    "父级ID",
	"isSystem":    "是否系统内置",
	"dataScope":   "数据范围",
	"oldPassword": "旧密码",
	"newPassword": "新密码",
	"clientId":    "客户端ID",
	"deviceType":  "设备类型",
	// 保留 error_translator.go 中的映射
	"UserId":      "用户ID",
	"UserName":    "用户名",
	"NickName":    "昵称",
	"Password":    "密码",
	"UserType":    "用户类型",
	"Email":       "邮箱",
	"Phonenumber": "手机号",
	"Sex":         "性别",
	"Avatar":      "头像",
	"Status":      "状态",
	"Remark":      "备注",
	"NewPassword": "新密码",
	"UserIds":     "用户ID列表",
	"SortOrder":   "显示顺序",
}

// getFieldName 获取字段的中文名称
func getFieldName(field string) string {
	if name, ok := fieldNameMap[field]; ok {
		return name
	}
	return field
}

// TranslateWithMsg 使用结构体中的 msg 标签翻译错误
// 支持在结构体字段上定义 msg 标签来自定义错误信息：
//
//	type CreateUserRequest struct {
//	    Email string `json:"email" binding:"required,email" msg:"请输入有效的邮箱地址"`
//	}
func TranslateWithMsg(err error, obj interface{}) string {
	if err == nil {
		return ""
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err.Error()
	}

	// 获取结构体的 msg 标签映射（key 为结构体字段名）
	msgMap := extractMsgMap(obj)

	var messages []string
	for _, e := range validationErrors {
		// e.Field() 返回结构体字段名（不是 json tag）
		// 优先使用 msg 标签定义的错误信息
		if msg, ok := msgMap[e.Field()]; ok && msg != "" {
			messages = append(messages, msg)
			continue
		}
		// 回退到翻译器
		messages = append(messages, e.Translate(trans))
	}

	return strings.Join(messages, "；")
}

// extractMsgMap 提取结构体的 msg 标签映射
// key 为 json tag（与 validator.FieldError.Field() 返回一致）
func extractMsgMap(obj interface{}) map[string]string {
	msgMap := make(map[string]string)
	if obj == nil {
		return msgMap
	}

	typ := reflect.TypeOf(obj)

	// 处理指针
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// 处理 interface{}（如果传入的是 interface{} 包装的结构体）
	if typ.Kind() == reflect.Interface {
		// 无法从 interface{} 类型获取字段信息，返回空 map
		// 调用方应该传入具体类型的指针
		return msgMap
	}

	if typ.Kind() != reflect.Struct {
		return msgMap
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if msg := field.Tag.Get("msg"); msg != "" {
			// 使用 json tag 作为 key（与 validator.FieldError.Field() 一致）
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 1)[0]
			if jsonName != "" && jsonName != "-" {
				msgMap[jsonName] = msg
			} else {
				msgMap[field.Name] = msg
			}
		}
	}

	return msgMap
}
