# Validator 使用指南

## 问题：同名字段不同场景需要不同提示语

例如：
- 用户注册时：`email` 是必填的，提示"请输入邮箱"
- 用户更新时：`email` 是选填的，提示"邮箱格式不正确"

---

## 方案一：使用结构体标签定义提示语（推荐）

使用自定义标签 `msg` 在结构体上直接定义错误信息：

```go
type CreateUserRequest struct {
    // 注册场景：邮箱必填
    Email string `json:"email" binding:"required,email" msg:"请输入有效的邮箱地址"`
}

type UpdateUserRequest struct {
    // 更新场景：邮箱选填
    Email string `json:"email" binding:"omitempty,email" msg:"邮箱格式不正确，请检查"`
}
```

实现方式：

```go
// RegisterTagNameFunc 获取 msg 标签作为错误信息
func init() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterTagNameFunc(func(fld reflect.StructField) string {
            // 优先使用 msg 标签
            if msg := fld.Tag.Get("msg"); msg != "" {
                return msg
            }
            // 其次使用 json 标签
            name := strings.SplitN(fld.Tag.Get("json"), ",", 1)[0]
            if name != "" && name != "-" {
                return name
            }
            return fld.Name
        })
    }
}
```

---

## 方案二：按场景分组映射

```go
package validator

// Scene 定义验证场景
type Scene string

const (
    SceneCreate Scene = "create"  // 创建场景
    SceneUpdate Scene = "update"  // 更新场景
    SceneLogin  Scene = "login"   // 登录场景
)

// sceneFieldMessages 按场景分组的字段提示
type sceneMessages map[Scene]map[string]string

var sceneFieldMessages = sceneMessages{
    SceneCreate: {
        "Email":       "请输入有效的邮箱地址",
        "Phonenumber": "请输入11位手机号",
        "Password":    "密码长度不能少于6位",
    },
    SceneUpdate: {
        "Email":       "邮箱格式不正确",
        "Phonenumber": "手机号格式不正确",
        "Password":    "新密码长度不能少于6位",
    },
    SceneLogin: {
        "UserName": "请输入用户名",
        "Password": "请输入密码",
    },
}

// TranslateWithScene 按场景翻译错误
func TranslateWithScene(err error, scene Scene) string {
    if err == nil {
        return ""
    }

    messages, ok := sceneFieldMessages[scene]
    if !ok {
        return Translate(err) // 回退到默认翻译
    }

    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        var errs []string
        for _, e := range validationErrors {
            field := e.Field()
            if msg, ok := messages[field]; ok {
                errs = append(errs, msg)
            } else {
                errs = append(errs, e.Translate(trans))
            }
        }
        return strings.Join(errs, "；")
    }

    return err.Error()
}
```

使用方式：

```go
// 创建用户场景
var createReq request.CreateUserRequest
if err := c.Bind(&createReq); err != nil {
    response.FailCode(c, 400, validator.TranslateWithScene(err, validator.SceneCreate))
    return
}

// 更新用户场景
var updateReq request.UpdateUserRequest
if err := c.Bind(&updateReq); err != nil {
    response.FailCode(c, 400, validator.TranslateWithScene(err, validator.SceneUpdate))
    return
}
```

---

## 方案三：使用嵌套结构体 + 接口（最灵活）

```go
// Validator 定义可自定义错误信息的接口
type Validator interface {
    Validate() error
    ErrorMessages() map[string]string
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    UserName string `json:"userName" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
}

// ErrorMessages 返回创建场景的错误提示
func (r CreateUserRequest) ErrorMessages() map[string]string {
    return map[string]string{
        "UserName": "用户名必须是3-20个字符",
        "Email":    "请输入有效的邮箱地址",
    }
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
    Email string `json:"email" binding:"omitempty,email"`
}

// ErrorMessages 返回更新场景的错误提示
func (r UpdateUserRequest) ErrorMessages() map[string]string {
    return map[string]string{
        "Email": "邮箱格式不正确，请检查",
    }
}
```

---

## 方案四：使用字段路径 + 类型名（精细控制）

```go
// fieldPathMessages 使用结构体类型名+字段路径作为 key
var fieldPathMessages = map[string]string{
    "CreateUserRequest.UserName": "注册时用户名必须是3-20个字符",
    "CreateUserRequest.Email":    "注册时请输入有效的邮箱",
    "UpdateUserRequest.Email":    "更新时邮箱格式不正确",
    "LoginRequest.UserName":      "登录时请填写用户名",
    "LoginRequest.Password":      "登录时请填写密码",
}

// TranslateWithType 根据结构体类型翻译错误
func TranslateWithType(err error, typeName string) string {
    if err == nil {
        return ""
    }

    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        var errs []string
        for _, e := range validationErrors {
            key := typeName + "." + e.Field()
            if msg, ok := fieldPathMessages[key]; ok {
                errs = append(errs, msg)
            } else {
                errs = append(errs, e.Translate(trans))
            }
        }
        return strings.Join(errs, "；")
    }

    return err.Error()
}
```

使用方式：

```go
var req request.CreateUserRequest
if err := c.Bind(&req); err != nil {
    response.FailCode(c, 400, validator.TranslateWithType(err, "CreateUserRequest"))
    return
}
```

---

## 推荐方案

| 场景 | 推荐方案 | 理由 |
|------|---------|------|
| 简单项目 | 方案一（msg 标签） | 简单直观，与字段定义在一起 |
| 多场景复用结构体 | 方案二（场景分组） | 结构体复用，提示语分离 |
| 大型项目 | 方案三（接口） | 最灵活，可扩展性强 |
| 精确控制 | 方案四（类型+字段路径） | 最精细，可针对每个结构体定制 |

---

## 当前项目采用方案

当前采用**方案一（msg 标签）+ 默认翻译**，如需更精细控制，可切换到方案二或方案四。
