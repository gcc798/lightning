package validator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

// 测试结构体 - 带 msg 标签
type TestCreateRequest struct {
	UserName string `json:"userName" binding:"required,min=3" msg:"用户名必须是3-20个字符"`
	Email    string `json:"email" binding:"required,email" msg:"请输入有效的企业邮箱"`
}

// 测试结构体 - 不带 msg 标签
type TestUpdateRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

func TestExtractMsgMap(t *testing.T) {
	req := TestCreateRequest{}

	typ := reflect.TypeOf(req)
	msgMap := make(map[string]string)

	t.Log("=== 结构体字段信息 ===")
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		msg := field.Tag.Get("msg")
		jsonTag := field.Tag.Get("json")
		jsonName := strings.SplitN(jsonTag, ",", 1)[0]
		t.Logf("Field.Name: %s, json: %s, jsonName: %s, msg: %s", field.Name, jsonTag, jsonName, msg)

		if msg != "" {
			if jsonName != "" && jsonName != "-" {
				msgMap[jsonName] = msg
			} else {
				msgMap[field.Name] = msg
			}
		}
	}

	t.Logf("msgMap: %+v", msgMap)

	// 验证 map 内容
	if msgMap["userName"] != "用户名必须是3-20个字符" {
		t.Errorf("msgMap[userName] = %q, want %q", msgMap["userName"], "用户名必须是3-20个字符")
	}
	if msgMap["email"] != "请输入有效的企业邮箱" {
		t.Errorf("msgMap[email] = %q, want %q", msgMap["email"], "请输入有效的企业邮箱")
	}
}

func TestFieldErrorField(t *testing.T) {
	Init()

	req := TestCreateRequest{UserName: "a", Email: "bad"}
	err := engine.Struct(&req)

	if err == nil {
		t.Fatal("expected error")
	}

	t.Log("=== 验证错误信息 ===")
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			t.Logf("Field(): %s, Tag(): %s", e.Field(), e.Tag())
		}
	}
}

func TestTranslate(t *testing.T) {
	Init()

	tests := []struct {
		name     string
		req      interface{}
		wantErr  bool
		contains []string
	}{
		{
			name: "使用默认翻译",
			req: struct {
				UserName string `json:"userName" binding:"required,min=3,max=20"`
			}{
				UserName: "a",
			},
			wantErr:  true,
			contains: []string{"用户名长度不能少于3个字符"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				msg := Translate(err)
				t.Logf("翻译结果: %s", msg)
				for _, want := range tt.contains {
					if !strings.Contains(msg, want) {
						t.Errorf("Translate() 结果不包含 %q, 实际结果: %s", want, msg)
					}
				}
			}
		})
	}
}

func TestTranslateWithMsg(t *testing.T) {
	Init()

	tests := []struct {
		name     string
		req      interface{}
		wantErr  bool
		contains []string
	}{
		{
			name: "使用 msg 标签 - 用户名错误",
			req: TestCreateRequest{
				UserName: "a",
				Email:    "test@example.com",
			},
			wantErr:  true,
			contains: []string{"用户名必须是3-20个字符"},
		},
		{
			name: "使用 msg 标签 - 邮箱错误",
			req: TestCreateRequest{
				UserName: "testuser",
				Email:    "invalid-email",
			},
			wantErr:  true,
			contains: []string{"请输入有效的企业邮箱"},
		},
		{
			name: "使用 msg 标签 - 多个错误",
			req: TestCreateRequest{
				UserName: "a",
				Email:    "invalid-email",
			},
			wantErr: true,
			contains: []string{
				"用户名必须是3-20个字符",
				"请输入有效的企业邮箱",
			},
		},
		{
			name: "无 msg 标签 - 回退到默认翻译",
			req: TestUpdateRequest{
				Email: "invalid-email",
			},
			wantErr:  true,
			contains: []string{"邮箱格式不正确"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用类型断言获取具体类型
			var msg string
			switch req := tt.req.(type) {
			case TestCreateRequest:
				err := engine.Struct(&req)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err != nil {
					msg = TranslateWithMsg(err, &req)
				}
			case TestUpdateRequest:
				err := engine.Struct(&req)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err != nil {
					msg = TranslateWithMsg(err, &req)
				}
			default:
				t.Fatalf("未知的请求类型: %T", tt.req)
			}

			if msg != "" {
				t.Logf("翻译结果: %s", msg)
				for _, want := range tt.contains {
					if !strings.Contains(msg, want) {
						t.Errorf("TranslateWithMsg() 结果不包含 %q, 实际结果: %s", want, msg)
					}
				}
			}
		})
	}
}
