# Kratos 工程说明

## 目标

在 `kratos/` 目录下基于 Kratos 框架实现一套全新的后端，同时保持以下原则：

- `native/` 作为业务基线和只读参考，不直接修改。
- `web-react/` 尽量保持零改动可用。
- `sys-api` 作为对外 HTTP 服务。
- `sys-rpc` 作为对内 gRPC 服务，承接核心业务逻辑和数据库访问。

## 工程结构

当前重写方案采用 monorepo 结构，但每个服务内部仍保持标准 Kratos 分层。

```text
kratos/
├── api/
│   └── system/v1/
├── application/
│   ├── sys-api/
│   │   ├── cmd/server
│   │   ├── configs
│   │   └── internal/{biz,data,server,service,conf}
│   └── sys-rpc/
│       ├── cmd/server
│       ├── configs
│       └── internal/{biz,data,server,service,conf}
├── cmd/
├── docs/
├── pkg/
├── third_party/
├── Makefile
└── go.mod
```

## 服务边界

### sys-api

职责：

- 向 `web-react/` 暴露 HTTP 接口。
- 解析请求参数。
- 把响应包装成兼容 `native` 的 HTTP 返回结构。
- 处理 token、鉴权、缓存、HTTP 中间件等外层能力。
- 通过 gRPC 调用 `sys-rpc`。

非职责：

- 不直接访问数据库。
- 不承担 system 域的核心持久化业务逻辑。

### sys-rpc

职责：

- 对内暴露 gRPC 接口。
- 承担 system 域核心业务逻辑。
- 承担 PostgreSQL 访问。
- 维护 Ent schema、repo、事务和底层资源。

## 契约策略

使用 proto 作为单一事实来源。

- `api/system/v1/*.proto` 同时定义 HTTP 和 gRPC 契约。
- HTTP 路径、方法、请求字段、返回字段都要尽量与 `native/` 对齐。
- 除非 `web-react/` 本身偏离 `native/`，否则不应要求前端为 Kratos 版本单独适配。

## ORM 策略

当前使用 `ent + PostgreSQL`。

- Ent schema 位于 `application/sys-rpc/ent/schema`
- 数据访问 repo 位于 `application/sys-rpc/internal/data`
- 业务编排位于 `application/sys-rpc/internal/biz`

## 开发阶段

### 第一阶段

- 创建 `kratos/` 基础工程结构。
- 创建 `sys-api` 与 `sys-rpc` 两个服务。
- 建立共享 proto 目录结构。
- 建立 `Makefile` 和生成/构建命令。
- 先保证整个工程可编译。

### 第二阶段

- 增加第一批 proto 契约，例如 `health / auth`。
- 生成 HTTP / gRPC 代码。
- 先形成可运行的最小骨架。

### 后续阶段

- `auth / captcha / me / menu`
- `user / role / org`
- `dict / config`
- `loginlog / operlog / storage-env / attachment`
- 最后对照 `native/` 与 `web-react/` 做完整兼容性回归。

## 当前实现约束

1. 不修改 `native/`
2. 每个服务内部遵守 Kratos 推荐分层
3. Kratos 支持生成的部分尽量通过生成完成
4. `sys-api` 不直接查库
5. 契约保持强类型，不做泛化 JSON RPC 包装

## 当前工程状态

当前仓库已经具备一套可编译的 Kratos 工程骨架，并完成了主要基础设施：

- `sys-api` 对外提供 HTTP 能力
- `sys-api` 通过 gRPC 调 `sys-rpc`
- `sys-rpc` 承担数据库访问、Redis、JWT、注册发现、存储、指标等底层能力
- PostgreSQL 底层驱动当前使用 `pgx/v5/stdlib`，上层通过 `database/sql + ent` 接入
- `application/sys-rpc/ent/schema` 已完成 Ent 生成链路
- `Wire` 已规范为标准 `ProviderSet` 风格
- OpenAPI 已按服务/版本输出到 `api/system/v1.openapi.yaml`

## 常用 Make 命令

常用命令包括：

- `make init`
- `make conf`
- `make proto-all`
- `make ent`
- `make wire`
- `make fmt`
- `make build-all`
- `make test`

更详细的变量和目标说明见：

- [docs/makefile-reference.md](docs/makefile-reference.md)

## 已抽取到 pkg 的通用能力

当前已经沉淀到 `pkg/`，后续可复用于其他微服务的能力有：

- `pkg/configx`
  - 共享配置加载能力
  - 共享时长解析辅助函数
- `pkg/grpcx`
  - 共享 gRPC client 建连逻辑
  - 同时支持 `direct` 与 `discovery` 两种模式
- `pkg/registryx`
  - 共享 etcd client 创建
  - 共享服务注册创建
  - 共享服务发现创建
- `pkg/metrics`
  - 共享 Prometheus 指标定义与辅助函数

## 后续适合继续抽取到 pkg 的候选能力

下面这些能力已经有较明显的复用价值，但目前先记录在文档中，不立即抽取。

### 优先级 1

- `pkg/httpx`
  - 目标：
    - 统一 HTTP 响应编码
    - 统一 HTTP 错误编码
    - 统一分页/数据包装方式
  - 当前来源：
    - `application/sys-api/internal/server/codec.go`

- `pkg/authx`
  - 目标：
    - 当前用户上下文辅助函数
    - token/header 提取
    - client IP 和 user-agent 提取
    - 可复用的认证元信息透传能力
  - 当前来源：
    - `application/sys-api/internal/data/context.go`

### 优先级 2

- `pkg/observabilityx`
  - 目标：
    - 统一 HTTP metrics 中间件
    - 统一慢 SQL / 慢 Redis 日志辅助能力
    - 统一可观测性初始化入口
  - 当前来源：
    - `application/sys-api/internal/server/metrics_handler.go`
    - `application/sys-rpc/internal/data/db_observability.go`
    - `application/sys-rpc/internal/data/redis_observability.go`

- `pkg/wsx`
  - 目标：
    - 可复用的 websocket hub
    - 可复用的连接注册表
    - 可复用的心跳与广播原语
  - 当前来源：
    - `application/sys-api/internal/server/websocket_hub.go`
    - `application/sys-api/internal/server/websocket_handler.go`

### 优先级 3

- `pkg/permx`
  - 目标：
    - 可复用的 operation 到 permission 映射辅助能力
    - 可复用的权限中间件骨架
  - 说明：
    - 只有在多个服务的权限命名约定稳定后，才值得抽取
  - 当前来源：
    - `application/sys-api/internal/server/middleware.go`
    - `application/sys-api/internal/server/permissions.go`

- `pkg/operlogx`
  - 目标：
    - 可复用的操作日志中间件
    - 可复用的请求/响应裁剪和归一化辅助能力
  - 说明：
    - 只有在日志 DTO 与审计约定稳定后，才值得抽取
  - 当前来源：
    - `application/sys-api/internal/server/operlog_middleware.go`

## 未来抽取到 pkg 的判断规则

只有同时满足以下条件，才建议把代码从服务内部抽取到 `pkg/`：

1. 这段逻辑不依赖某个服务的业务语义。
2. 至少有两个服务可以合理复用。
3. 抽象已经足够稳定，不会频繁返工。
4. 抽出的包不需要反向依赖服务内部的 `biz` 或 `data` 代码。

如果逻辑仍然带有明显的 system 业务假设、权限命名假设、DTO 耦合，应该先继续留在服务内部。

## 文档索引

如果你要继续了解这套工程，建议按下面顺序阅读：

1. proto 到 HTTP / gRPC 串联流程：
   - [docs/proto-http-grpc-flow.md](docs/proto-http-grpc-flow.md)
2. Ent 使用说明：
   - [docs/ent-usage.md](docs/ent-usage.md)
3. Wire 使用说明：
   - [docs/wire-usage.md](docs/wire-usage.md)
4. HTTP 手工路由与 proto 路由边界：
   - [docs/http-routing-boundary.md](docs/http-routing-boundary.md)
5. Makefile 详细说明：
   - [docs/makefile-reference.md](docs/makefile-reference.md)
