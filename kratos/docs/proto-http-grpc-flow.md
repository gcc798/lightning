# Kratos 中 Proto 到 HTTP / gRPC 串联流程说明

这份文档用 `登录接口` 作为例子，说明当前 `kratos/` 工程里，Kratos 是如何通过同一份 `proto`：

1. 定义 HTTP 接口
2. 定义 gRPC 接口
3. 在 `sys-api` 暴露 HTTP
4. 在 `sys-api` 内部通过 gRPC 调 `sys-rpc`
5. 在 `sys-rpc` 内部执行业务逻辑

---

## 1. 核心思路

在当前工程里，`proto` 是单一契约源：

- `api/system/v1/*.proto`

同一份 `proto` 同时承担两件事：

1. 定义 HTTP 路由
2. 定义 gRPC Service / Message

也就是说，Kratos 不是像 go-zero 那样分成：

- `.api` 负责 HTTP
- `.proto` 负责 RPC

而是直接用一份 `proto` 统一描述两种接口。

---

## 2. 登录接口的 proto 定义

登录接口定义在：

- [../api/system/v1/auth.proto](../api/system/v1/auth.proto)

关键定义：

```proto
service AuthService {
  rpc Login (LoginRequest) returns (LoginReply) {
    option (google.api.http) = {
      post: "/login"
      body: "*"
    };
  }
}
```

这里有两层含义：

1. `rpc Login(LoginRequest) returns (LoginReply)`
   - 定义了一个 gRPC 方法
   - 方法名是 `Login`
   - 请求结构是 `LoginRequest`
   - 返回结构是 `LoginReply`

2. `option (google.api.http)`
   - 告诉 Kratos：这个 gRPC 方法同时映射成 HTTP 接口
   - 当前映射结果是：
     - `POST /login`

所以，**一份定义同时告诉系统：**

- 有一个 gRPC 方法 `AuthService.Login`
- 这个方法还要暴露成 HTTP 的 `POST /login`

---

## 3. 代码生成后会得到什么

执行：

```bash
make proto-all
```

后，会基于 `auth.proto` 生成几类代码。

### 3.1 HTTP 绑定代码

生成文件：

- [../api/system/v1/auth_http.pb.go](../api/system/v1/auth_http.pb.go)

你可以直接看到：

```go
func RegisterAuthServiceHTTPServer(s *http.Server, srv AuthServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/login", _AuthService_Login0_HTTP_Handler(srv))
	r.POST("/logout", _AuthService_Logout0_HTTP_Handler(srv))
	r.POST("/auth/refresh", _AuthService_RefreshToken0_HTTP_Handler(srv))
	r.GET("/me", _AuthService_Me0_HTTP_Handler(srv))
}
```

这说明 `proto` 里的：

```proto
post: "/login"
```

已经被转成了真正的 HTTP 路由注册代码。

### 3.2 gRPC client / server 代码

生成文件：

- [../api/system/v1/auth_grpc.pb.go](../api/system/v1/auth_grpc.pb.go)

这里会生成：

1. gRPC Client 接口

```go
type AuthServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
}
```

2. gRPC Server 接口

```go
type AuthServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginReply, error)
}
```

这就是后面 `sys-api` 和 `sys-rpc` 对接的桥梁：

- `sys-api` 用 `AuthServiceClient`
- `sys-rpc` 实现 `AuthServiceServer`

---

## 4. sys-api 如何把 proto 变成 HTTP 接口

### 4.1 HTTP Server 注册 auth service

注册位置：

- [../application/sys-api/internal/server/http.go](../application/sys-api/internal/server/http.go)

关键代码：

```go
v1.RegisterAuthServiceHTTPServer(srv, authSvc)
```

这里的 `authSvc` 是 `sys-api` 的一个 service 实现。

当 Kratos 启动 HTTP Server 时，这行代码就会把：

- `POST /login`

绑定到：

- `authSvc.Login(...)`

### 4.2 sys-api 的 service 层

文件：

- [../application/sys-api/internal/service/authservice.go](../application/sys-api/internal/service/authservice.go)

关键代码：

```go
func (s *AuthServiceService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return s.uc.Login(ctx, req)
}
```

这里很薄，职责只是：

1. 接住 HTTP 转换后的请求对象 `LoginRequest`
2. 调用 usecase

也就是说，`sys-api service` 在这个项目里主要做的是：

- 适配 proto 接口
- 调用业务层

### 4.3 sys-api 的 biz 层

文件：

- [../application/sys-api/internal/biz/auth.go](../application/sys-api/internal/biz/auth.go)

关键代码：

```go
func (uc *AuthUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return uc.repo.Login(ctx, req)
}
```

这层当前也比较薄，主要职责是：

- 把请求继续交给 repo

在 `sys-api` 里，这个 repo 不是数据库 repo，而是 **RPC repo**。

### 4.4 sys-api 的 data 层：真正调用 sys-rpc

文件：

- [../application/sys-api/internal/data/auth.go](../application/sys-api/internal/data/auth.go)

关键代码：

```go
type AuthRepo struct {
	client v1.AuthServiceClient
}

func NewAuthRepo(clients *RPCClientSet) *AuthRepo {
	return &AuthRepo{client: clients.Auth}
}

func (r *AuthRepo) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return r.client.Login(ctx, req)
}
```

这里就是 `sys-api -> sys-rpc` 的真正调用点。

注意这一点：

- `client` 类型是 `v1.AuthServiceClient`
- 它来自 `auth_grpc.pb.go` 生成代码

所以 `sys-api` 并不是自己实现登录逻辑，而是：

1. 持有一个 gRPC client
2. 直接调用 `client.Login(...)`

---

## 5. sys-api 的 gRPC client 是怎么创建的

文件：

- [../application/sys-api/internal/data/rpc_client.go](../application/sys-api/internal/data/rpc_client.go)

### 5.1 统一建连

关键代码：

```go
func NewRPCClientSet(cfg conf.RPC) (*RPCClientSet, error) {
	conn, cleanup, err := newRPCConn(cfg)
	...
	return &RPCClientSet{
		Auth: v1.NewAuthServiceClient(conn),
		...
	}, nil
}
```

这里做了两件事：

1. 先建立一条统一的 gRPC 连接
2. 再基于这条连接构造多个 typed client

其中 auth 对应的是：

```go
v1.NewAuthServiceClient(conn)
```

这也是 `auth_grpc.pb.go` 生成出来的代码。

### 5.2 连接模式

`newRPCConn(cfg)` 支持两种模式：

1. `direct`
   - 直接连指定地址
2. `discovery`
   - 通过注册中心发现服务

也就是说，`sys-api` 本身并不关心登录逻辑，它只负责：

- 连上 `sys-rpc`
- 调用 `AuthServiceClient.Login(...)`

---

## 6. sys-api 这一侧的完整调用链

到这里，登录接口在 `sys-api` 里的完整链路是：

```text
POST /login
  -> proto 生成的 HTTP Handler
  -> sys-api/internal/service/authservice.go: Login
  -> sys-api/internal/biz/auth.go: Login
  -> sys-api/internal/data/auth.go: Login
  -> AuthServiceClient.Login(ctx, req)
  -> gRPC 调 sys-rpc
```

---

## 7. sys-rpc 如何接住这个 gRPC 调用

### 7.1 gRPC Server 注册 auth service

`sys-rpc` 启动时会注册 gRPC 服务。

装配位置：

- [../application/sys-rpc/cmd/server/wire.go](../application/sys-rpc/cmd/server/wire.go)

这里会创建：

- `service.Services`
- gRPC Server

### 7.2 sys-rpc 的 service 层

文件：

- [../application/sys-rpc/internal/service/authservice.go](../application/sys-rpc/internal/service/authservice.go)

关键代码：

```go
func (s *AuthServiceService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return s.uc.Login(ctx, req)
}
```

它实现的就是 `auth_grpc.pb.go` 里生成的：

- `AuthServiceServer`

所以当 `sys-api` 发起：

```go
client.Login(ctx, req)
```

时，`sys-rpc` 这边最终会落到：

- `AuthServiceService.Login(...)`

### 7.3 sys-rpc 的 biz 层

文件：

- [../application/sys-rpc/internal/biz/auth.go](../application/sys-rpc/internal/biz/auth.go)

这个文件里就是登录的核心业务逻辑。

当前登录流程大致是：

1. 校验验证码
2. 校验 clientKey / clientSecret / grantType
3. 按登录方式查用户
4. 校验用户状态
5. 校验密码
6. 签发 access token / refresh token
7. 更新登录状态
8. 写登录日志
9. 返回 `LoginReply`

你可以在这里看到这些关键调用：

```go
uc.res.AuthenticateClient(...)
uc.res.FindUserByAccount(...)
uc.res.ResolveXcxUser(...)
uc.res.IssueSession(...)
uc.res.UpdateUserLoginState(...)
uc.res.CreateLoginLogEntry(...)
```

也就是说，真正的业务规则都在 `sys-rpc biz`。

### 7.4 sys-rpc 的 data / resources 层

`sys-rpc biz` 依赖的是：

- `data.Resources`

它负责真正访问：

- PostgreSQL / Ent
- Redis
- JWT
- 存储
- 微信

这部分资源初始化在：

- [../application/sys-rpc/cmd/server/wire.go](../application/sys-rpc/cmd/server/wire.go)

并由：

- `newResources(cfg)`

创建出来。

所以 `sys-rpc` 的结构是：

```text
gRPC service
  -> biz
  -> data / resources
  -> DB / Redis / JWT / 外部系统
```

---

## 8. sys-rpc 这一侧的完整调用链

登录请求进到 `sys-rpc` 后，链路是：

```text
AuthServiceClient.Login(ctx, req)
  -> proto 生成的 gRPC handler
  -> sys-rpc/internal/service/authservice.go: Login
  -> sys-rpc/internal/biz/auth.go: Login
  -> data.Resources 的各类方法
  -> PostgreSQL / Redis / JWT / 日志
  -> 返回 LoginReply
```

---

## 9. 最完整的一条串行链路

把两边合起来，登录接口完整链路就是：

```text
前端 POST /login
  -> auth.proto 中 Login 的 http option
  -> auth_http.pb.go 生成的 HTTP Handler
  -> sys-api/internal/service/authservice.go
  -> sys-api/internal/biz/auth.go
  -> sys-api/internal/data/auth.go
  -> v1.AuthServiceClient.Login(ctx, req)
  -> gRPC 网络调用
  -> auth_grpc.pb.go 生成的 gRPC handler
  -> sys-rpc/internal/service/authservice.go
  -> sys-rpc/internal/biz/auth.go
  -> sys-rpc/internal/data / Resources
  -> DB / Redis / JWT / 外部依赖
  -> 返回 LoginReply
  -> sys-api 原样拿到结果
  -> HTTP ResponseEncoder 包装后返回前端
```

---

## 10. Wire 在这里起什么作用

当前项目使用 `wire` 做依赖注入。

### 10.1 sys-api

文件：

- [../application/sys-api/cmd/server/wire.go](../application/sys-api/cmd/server/wire.go)

这里的装配顺序是：

```text
Repositories
  -> Usecases
  -> Services
  -> HTTP Server
  -> Kratos App
```

对应：

1. `newRepositories(cfg)` 创建 RPC repo
2. `biz.NewUsecases(...)`
3. `service.NewServices(...)`
4. `newHTTPServer(...)`
5. `newApp(...)`

### 10.2 sys-rpc

文件：

- [../application/sys-rpc/cmd/server/wire.go](../application/sys-rpc/cmd/server/wire.go)

这里的装配顺序是：

```text
Resources
  -> Usecases
  -> Services
  -> gRPC Server
  -> Kratos App
```

对应：

1. `newResources(cfg)` 创建 DB / Redis / JWT 等资源
2. `biz.NewUsecases(...)`
3. `service.NewServices(...)`
4. `newGRPCServer(...)`
5. `newApp(...)`

所以 `wire` 的职责是：

- 把依赖关系串起来
- 不让 `main.go` 手写一大堆 `NewXxx(...)`

---

## 11. 你可以怎么理解这套结构

如果只抓住一句话，可以这样理解：

### 在 Kratos 里：

- `proto` 负责定义契约
- 生成代码负责把契约变成 HTTP / gRPC 桥接层
- `sys-api` 负责对外 HTTP 和转发 RPC
- `sys-rpc` 负责真正业务逻辑和数据访问

### 登录接口就是：

```text
proto 定义一次
  -> 自动生成 HTTP 路由
  -> 自动生成 gRPC client / server
  -> sys-api 调 sys-rpc
  -> sys-rpc 处理业务
```

---

## 12. 后续看别的模块时怎么套这个模板

你后面看 `menu / user / role / dict / config`，都可以按同一个模板看：

1. 先看 `api/system/v1/*.proto`
2. 再看 `api/system/v1/*_http.pb.go`
3. 看 `sys-api/internal/service/*.go`
4. 看 `sys-api/internal/biz/*.go`
5. 看 `sys-api/internal/data/*.go`
6. 看 `sys-rpc/internal/service/*.go`
7. 看 `sys-rpc/internal/biz/*.go`
8. 最后看 `sys-rpc/internal/data/*.go`

只要按这个顺序看，Kratos 这套“proto 驱动 HTTP + gRPC”的结构就会很清楚。
