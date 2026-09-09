package errors

/*
结构化错误处理使用示例

================================================================================
1. 业务逻辑错误（可以直接向用户展示）
================================================================================

示例场景：Service 层根据实际参数查询数据库发现用户没有关联信息不可以操作

// Service 层代码
func (s *userService) SomeBusinessLogic(ctx context.Context, userId int64) error {
    // 查询用户信息
    user, err := s.db.FindByID(userId)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 业务逻辑错误：用户不存在
            return apperrors.NewBusiness(apperrors.CodeUserNotFound, "用户不存在")
        }
        // 数据库错误属于基础设施错误
        return apperrors.NewInfrastructure(apperrors.CodeDatabaseError, "数据库查询失败", err)
    }

    // 检查业务规则
    if user.RelatedInfo == nil {
        // 业务逻辑错误：返回自定义消息
        return apperrors.NewBusiness(apperrors.CodeNoPermission, "你没有关联信息，无法操作")
    }

    return nil
}

// Controller 层代码
func (c *userController) SomeAction(ctx *echo.Context) {
    err := c.userService.SomeBusinessLogic(ctx.Request().Context(), userId)
    if err != nil {
        response.Error(ctx, err)  // 统一错误处理
        return
    }
    response.Success(ctx, data)
}

响应结果：
HTTP Status: 200
{
    "code": 10007,
    "msg": "你没有关联信息，无法操作"
}

前端处理：直接弹框提示 msg 即可


================================================================================
2. 基础设施错误（需要记录日志，不向用户暴露细节）
================================================================================

示例场景：数据库连接超时、Redis 连接失败、S3 上传失败等

// Service 层代码
func (s *userService) CreateUser(ctx context.Context, req *request.CreateUserRequest) error {
    user := &model.User{...}

    // 数据库操作失败
    if err := s.db.Create(user).Error; err != nil {
        // 基础设施错误：记录详细日志，向用户返回统一文案
        return apperrors.NewInfrastructure(
            apperrors.CodeDatabaseError,
            "创建用户时数据库写入失败",  // 内部日志消息
            err,                          // 原始错误
        )
    }

    // Redis 缓存失败
    if err := s.redis.Set(ctx, key, value, ttl).Err(); err != nil {
        return apperrors.NewInfrastructure(
            apperrors.CodeRedisError,
            "缓存写入失败",
            err,
        )
    }

    return nil
}

// Controller 层代码
func (c *userController) Create(ctx *echo.Context) {
    err := c.userService.CreateUser(ctx.Request().Context(), &req)
    if err != nil {
        response.Error(ctx, err)  // 统一错误处理
        return
    }
    response.Success(ctx, data)
}

响应结果：
HTTP Status: 500
{
    "code": 20001,
    "msg": "服务暂时不可用，请稍后重试"  // 统一的用户友好文案
}

日志记录（自动记录）：
level=ERROR msg="基础设施错误" code=20001 internal_message="创建用户时数据库写入失败"
error="pq: duplicate key value violates unique constraint" path="/api/v1/user"
method="POST" params={"username":"test"}


================================================================================
3. 系统级错误（如 panic、数组越界，需要记录堆栈）
================================================================================

重要说明：系统错误是无法预知的（如数组越界、空指针等），Go 语言通过全局中间件
自动捕莹 panic，无需开发者在每个函数中手动 defer recover。

这类似于 Java Spring Boot 的 @ControllerAdvice，在 HTTP 层面统一拦截异常。

示例场景：数组越界、空指针等系统异常

// Service 层代码（无需手动处理 panic）
func (s *someService) DangerousOperation(ctx context.Context) error {
    // 开发者可以直接写业务代码，不用担心 panic
    arr := []int{1, 2, 3}
    _ = arr[10]  // panic: index out of range

    // 或者其他可能 panic 的操作
    var user *model.User
    _ = user.Name  // panic: nil pointer dereference

    return nil
}

// Controller 层代码（无需特殊处理）
func (c *someController) DangerousAction(ctx *echo.Context) {
    err := c.someService.DangerousOperation(ctx.Request().Context())
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx, data)
}

全局中间件自动捕莹（已在 middleware.Recovery 中实现）：
- 自动捕莹所有未处理的 panic
- 自动记录详细堆栈信息
- 自动返回统一的用户友好文案
- 无需开发者在每个函数中写 defer recover

响应结果：
HTTP Status: 500
{
    "code": 30001,
    "msg": "系统异常，请联系管理员"  // 统一文案
}

日志记录（由 Recovery 中间件自动记录）：
level=ERROR msg="系统 Panic 捕获" error="runtime error: index out of range [10] with length 3"
path="/api/v1/dangerous" method="POST" params={}
stack="goroutine 123 [running]:\nruntime/debug.Stack()\n..."


================================================================================
4. 错误码规划
================================================================================

业务模块错误码（10000-19999）：可直接向用户展示
  - 用户模块：10xxx
    10001: 用户不存在
    10002: 密码错误
    10003: 用户已禁用
    10004: 用户名已存在
    10005: 手机号已存在
    10006: 邮箱已存在
    10007: 无权限操作

  - 角色模块：11xxx
    11001: 角色不存在
    11002: 角色名已存在

  - 菜单模块：12xxx
    12001: 菜单不存在

基础设施错误码（20000-29999）：需记录日志，返回统一文案
  - 数据库：200xx
    20001: 数据库错误
    20002: 数据库超时
    20003: 数据库连接断开

  - 缓存：201xx
    20101: Redis 错误
    20102: Redis 超时

  - 存储：202xx
    20201: S3 错误
    20202: S3 上传失败

系统级错误码（30000-39999）：需记录堆栈，返回统一文案
  30001: Panic 错误
  30002: 数组越界
  30003: 空指针异常


================================================================================
5. 最佳实践
================================================================================

1. Service 层抛出精确的错误：
   - 业务逻辑问题 → NewBusiness()
   - 基础设施问题 → NewInfrastructure()
   - 系统异常（panic）→ 无需手动处理，由中间件自动捕获

2. Controller 层统一使用 response.Error(c, err)

3. 不要在 Controller 层判断错误类型，让 Error() 函数自动处理

4. 基础设施错误的 message 参数用于内部日志，不会暴露给用户

5. 确保在 main.go 或 router 中注册了 Recovery 中间件

6. 可以在 context 中设置 logger：c.Set("logger", logger)

7. 不要在每个函数中写 defer recover，这是反模式！
   系统级 panic 由全局中间件统一处理

8. 只在极少数情况下需要手动 recover：
   - goroutine 内部（中间件无法捕获）
   - 需要优雅降级而不是返回错误的场景

   例如：
   go func() {
       defer func() {
           if r := recover() {
               logger.Error("goroutine panic", zap.Any("error", r))
           }
       }()
       // goroutine 业务代码
   }()
*/
