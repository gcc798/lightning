# 密码工具使用指南

## 概述

`password.go` 提供了安全的密码哈希和验证功能，基于 bcrypt 算法实现。

## 特性

- ✅ 使用 bcrypt 算法（行业标准）
- ✅ 自动处理盐值（每次生成的哈希都不同）
- ✅ 可配置的成本因子
- ✅ 完整的错误处理
- ✅ 100% 测试覆盖率
- ✅ 线程安全

## 快速开始

### 基本使用（推荐）

```go
package main

import (
    "fmt"
    "github.com/gcc798/microservice-kit/internal/utils"
)

func main() {
    // 1. 生成密码哈希
    password := "admin123"
    hash, err := utils.HashPassword(password)
    if err != nil {
        panic(err)
    }
    fmt.Println("哈希:", hash)
    
    // 2. 验证密码
    err = utils.VerifyPassword(hash, password)
    if err != nil {
        fmt.Println("密码错误")
    } else {
        fmt.Println("密码正确")
    }
}
```

### 使用自定义哈希器

```go
package main

import (
    "github.com/gcc798/microservice-kit/internal/utils"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    // 创建自定义成本因子的哈希器
    hasher := utils.NewBcryptHasher(12) // 更高的安全性，但更慢
    
    hash, err := hasher.HashPassword("mypassword")
    if err != nil {
        panic(err)
    }
    
    err = hasher.VerifyPassword(hash, "mypassword")
    if err != nil {
        panic(err)
    }
}
```

## API 文档

### 全局函数

#### HashPassword

```go
func HashPassword(password string) (string, error)
```

生成密码的 bcrypt 哈希值。

**参数**：
- `password`: 明文密码

**返回**：
- `string`: bcrypt 哈希值（60个字符）
- `error`: 错误信息

**示例**：
```go
hash, err := utils.HashPassword("admin123")
if err != nil {
    // 处理错误
}
// hash: $2a$10$...
```

#### VerifyPassword

```go
func VerifyPassword(hashedPassword, password string) error
```

验证密码是否匹配哈希值。

**参数**：
- `hashedPassword`: bcrypt 哈希值
- `password`: 明文密码

**返回**：
- `error`: 如果匹配返回 nil，否则返回错误

**示例**：
```go
err := utils.VerifyPassword(hash, "admin123")
if err != nil {
    if err == utils.ErrPasswordMismatch {
        fmt.Println("密码错误")
    }
}
```

#### MustHashPassword

```go
func MustHashPassword(password string) string
```

生成密码哈希，失败则 panic。仅用于测试或初始化场景。

**示例**：
```go
hash := utils.MustHashPassword("admin123")
```

### BcryptHasher 类型

#### NewBcryptHasher

```go
func NewBcryptHasher(cost int) *BcryptHasher
```

创建指定成本因子的哈希器。

**参数**：
- `cost`: 成本因子（4-31），推荐 10-14

**成本因子说明**：
- `4`: 最快，最不安全（仅测试用）
- `10`: 默认值，平衡性能和安全性
- `12`: 更安全，稍慢
- `14`: 非常安全，较慢

#### NewDefaultBcryptHasher

```go
func NewDefaultBcryptHasher() *BcryptHasher
```

创建默认成本因子（10）的哈希器。

## 错误处理

### 错误类型

```go
var (
    ErrPasswordEmpty    = errors.New("password cannot be empty")
    ErrHashEmpty        = errors.New("hash cannot be empty")
    ErrPasswordMismatch = errors.New("password mismatch")
)
```

### 错误处理示例

```go
err := utils.VerifyPassword(hash, password)
switch err {
case nil:
    fmt.Println("登录成功")
case utils.ErrPasswordMismatch:
    fmt.Println("密码错误")
case utils.ErrPasswordEmpty:
    fmt.Println("密码不能为空")
case utils.ErrHashEmpty:
    fmt.Println("哈希值不能为空")
default:
    fmt.Printf("验证失败: %v\n", err)
}
```

## 实际应用场景

### 场景 1：用户注册

```go
func RegisterUser(username, password string) error {
    // 1. 验证密码强度（略）
    
    // 2. 生成密码哈希
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        return fmt.Errorf("密码加密失败: %w", err)
    }
    
    // 3. 存储到数据库
    user := &User{
        Username: username,
        Password: hashedPassword,
    }
    return db.Create(user).Error
}
```

### 场景 2：用户登录

```go
func Login(username, password string) (*User, error) {
    // 1. 从数据库查询用户
    var user User
    err := db.Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, fmt.Errorf("用户不存在")
    }
    
    // 2. 验证密码
    err = utils.VerifyPassword(user.Password, password)
    if err != nil {
        if err == utils.ErrPasswordMismatch {
            return nil, fmt.Errorf("密码错误")
        }
        return nil, fmt.Errorf("验证失败: %w", err)
    }
    
    // 3. 登录成功
    return &user, nil
}
```

### 场景 3：修改密码

```go
func ChangePassword(userId int64, oldPassword, newPassword string) error {
    // 1. 查询用户
    var user User
    err := db.First(&user, userId).Error
    if err != nil {
        return fmt.Errorf("用户不存在")
    }
    
    // 2. 验证旧密码
    err = utils.VerifyPassword(user.Password, oldPassword)
    if err != nil {
        return fmt.Errorf("旧密码错误")
    }
    
    // 3. 生成新密码哈希
    newHash, err := utils.HashPassword(newPassword)
    if err != nil {
        return fmt.Errorf("密码加密失败: %w", err)
    }
    
    // 4. 更新数据库
    return db.Model(&user).Update("password", newHash).Error
}
```

### 场景 4：重置密码

```go
func ResetPassword(userId int64, newPassword string) error {
    // 1. 生成新密码哈希
    newHash, err := utils.HashPassword(newPassword)
    if err != nil {
        return fmt.Errorf("密码加密失败: %w", err)
    }
    
    // 2. 更新数据库
    return db.Model(&User{}).
        Where("user_id = ?", userId).
        Update("password", newHash).Error
}
```

## 性能考虑

### 基准测试结果

```
BenchmarkHashPassword-10           25    46ms/op    5251 B/op    11 allocs/op
BenchmarkVerifyPassword-10         25    47ms/op    5212 B/op    12 allocs/op
```

### 性能建议

1. **成本因子选择**：
   - 开发环境：使用 `bcrypt.MinCost` (4) 加快测试速度
   - 生产环境：使用 `bcrypt.DefaultCost` (10) 或更高

2. **异步处理**：
   - 密码哈希和验证都是 CPU 密集型操作
   - 建议在 goroutine 中处理，避免阻塞主线程

3. **缓存考虑**：
   - 不要缓存密码哈希结果
   - 每次都应该重新生成（因为包含随机盐值）

## 安全最佳实践

1. **永远不要存储明文密码**
2. **使用足够的成本因子**（至少 10）
3. **不要在日志中记录密码或哈希值**
4. **实施密码强度策略**
5. **考虑添加登录失败次数限制**
6. **定期更新 bcrypt 库**

## 测试

运行单元测试：

```bash
go test -v ./internal/utils
```

运行基准测试：

```bash
go test -bench=. -benchmem ./internal/utils
```

查看测试覆盖率：

```bash
go test -cover ./internal/utils
```

## 常见问题

### Q: 为什么每次生成的哈希都不同？

A: bcrypt 会自动生成随机盐值并包含在哈希中，这是正常的安全特性。虽然哈希不同，但都能正确验证原始密码。

### Q: 哈希值的长度是多少？

A: bcrypt 哈希值固定为 60 个字符。

### Q: 如何选择成本因子？

A: 
- 开发/测试：4-6
- 生产环境：10-12
- 高安全需求：13-14

成本每增加 1，计算时间翻倍。

### Q: 可以使用其他哈希算法吗？

A: 可以实现 `PasswordHasher` 接口来使用其他算法（如 Argon2），但 bcrypt 已经足够安全。

## 相关资源

- [bcrypt 官方文档](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP 密码存储指南](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
