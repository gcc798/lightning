# native 微服务架构重构结果

## 1. 文档状态

本文记录 `native/` 已落地的最终架构和后续开发必须遵守的边界。原先规划的“单机与微服务双形态”已经取消；`native` 只保留真实微服务部署形态。

权限体系结果见 `NATIVE_PERMISSION_PLAN.md`。

## 2. 不可破坏的边界

1. `application/` 下一级目录就是一个可独立构建、启动和部署的进程。
2. IAM、SYS、Resource 的 controller、router、DTO、领域逻辑、model 和迁移必须位于各自的 `application/<service>/internal/`。
3. 根 `internal/` 只承载多个进程共同使用的技术设施和内部 RPC 契约，不承载 IAM、SYS、Resource 的业务实现或 GORM model。
4. 不允许重新创建聚合所有业务 HTTP 源码的 `application/api`，也不提供单体装配入口。
5. `gateway` 是唯一对外 HTTP/WS 入口；IAM、SYS、Resource 提供自己的 HTTP 与 gRPC 服务，并通过注册中心发布实例。
6. 每个领域服务拥有自己的数据库迁移和 Goose 版本表；Scheduler 不执行迁移。
7. `cmd/usermgr` 是一次性管理工具，不属于 `application/` 部署单元，可以直接保留自己的最小持久化结构。

## 3. 当前目录结构

```text
native/
├── application/
│   ├── gateway/                 # HTTP/WS 入口、鉴权、限流、服务发现与反向代理
│   ├── iam/
│   │   └── internal/            # IAM controller/router/DTO/domain/model/migrations
│   ├── sys/
│   │   └── internal/            # SYS controller/router/DTO/domain/model/migrations
│   ├── resource/
│   │   └── internal/            # Resource controller/router/DTO/domain/migrations
│   └── scheduler/               # Scheduler 进程与 Job
├── cmd/usermgr/                 # 管理员维护工具
├── api/<service>/v1/            # Proto 源文件
├── internal/
│   ├── api/<service>/v1/        # 生成的 gRPC 契约和薄适配
│   ├── httpserver/              # Echo 生命周期，不包含业务路由
│   ├── transport/               # gRPC server、注册生命周期和动态 client pool
│   ├── registry/                # Consul、etcd、in-process 实现和 selector
│   ├── openapi/                 # 统一生成的 Swagger 文档
│   ├── database/                # 共享数据库连接和 GORM 插件
│   └── platform/...             # JWT、Redis、存储、WebSocket 等技术设施
└── Makefile
```

旧目录 `application/api`、`application/server`、`application/runtime`、`application/monolith` 已删除。根 `internal/domain/{iam,sys,resource}` 也已删除，业务所有权不再跨服务共享。

各服务入口显式执行配置加载、日志初始化和依赖容器构建，不通过 `internal/app.Base` 或无必要的 `run` 函数隐藏装配顺序。

## 4. 服务所有权

| 进程 | 业务所有权 | 运行依赖 |
| --- | --- | --- |
| `gateway` | 无业务数据；统一 HTTP/WS 接入 | Registry、IAM gRPC |
| `iam` | user、role、menu、org、api_permission、auth、client | PostgreSQL、Redis、Registry |
| `sys` | dict、config、login_log、oper_log | PostgreSQL、Redis、Registry |
| `resource` | attachment、对象存储 | PostgreSQL、S3、Registry |
| `scheduler` | 定时任务和后台作业 | PostgreSQL、Registry、目标服务 gRPC |
| `cmd/usermgr` | 创建管理员、重置密码 | PostgreSQL |

服务划分沿事务和数据所有权，而不是沿 controller 文件划分。IAM 内相互依赖的用户、角色、菜单和权限表继续使用本地事务；本轮不引入分布式事务。

## 5. 通信与部署

南北向流量使用 HTTP：

```text
web-react -> gateway -> IAM / SYS / Resource HTTP
```

东西向调用使用 gRPC：

```text
Scheduler / domain client -> Registry -> gRPC instance
```

- Proto 位于 `api/<service>/v1`，运行 `make proto` 生成 `internal/api/<service>/v1`。
- IAM、SYS、Resource 启动 gRPC server，并将 HTTP/gRPC endpoint 注册到 Consul 或 etcd。
- Gateway 按路由前缀解析服务实例并反向代理，路径不改写。
- `internal/registry/inprocess.go` 只用于注册中心接口的单元测试和轻量本地组合，不是单体部署模式。
- 不做 HTTP 到 gRPC transcoding，也不让 Gateway 承担 BFF 聚合或业务编排。

当前路由所有权：

```text
/login /logout /auth/* /captcha/*
/api/v1/{user,role,menu,org,api-permission}/*
/system/user/* /resource/websocket                    -> iam

/api/v1/{config,dict,loginLog,operLog}/*              -> sys
/api/v1/attachment/*                                  -> resource
/swagger/* /health                                    -> gateway
```

## 6. 数据库迁移所有权

每个领域服务启动时只执行自己嵌入的 Goose 迁移：

| 服务 | SQL 目录 | 版本表 |
| --- | --- | --- |
| IAM | `application/iam/internal/migrations/sql` | `goose_iam_version` |
| SYS | `application/sys/internal/migrations/sql` | `goose_sys_version` |
| Resource | `application/resource/internal/migrations/sql` | `goose_resource_version` |

工程不使用 GORM `AutoMigrate`。这是全新试验工程，Schema 直接以当前目标状态建立，不为旧表和旧迁移保留兼容路径。

## 7. 已完成阶段

### 阶段一：清理

- 删除未使用的 RabbitMQ、`storage_env` 和重复存储管理逻辑。
- 对象存储收敛为 Resource 服务拥有的单一 S3 兼容配置。
- 删除空 bootstrap 和无调用者的容器接口。

### 阶段二：生命周期收敛

- 统一 Module 生命周期，删除重复的 Component 抽象。
- 公共 HTTP 生命周期进入 `internal/httpserver`。
- gRPC 和注册生命周期进入 `internal/transport`。

### 阶段三：服务业务归位

- IAM、SYS、Resource 的业务代码和 model 全部进入各自服务 `internal`。
- 中间件、HTTP response、validator 等跨服务技术能力进入根 `internal`。
- `cmd/usermgr` 和 `runtimeconfig` 使用自己的最小 persistence struct，不导入服务私有包。

### 阶段四：真实 RPC 与注册中心

- IAM、SYS、Resource 已提供真实 gRPC server。
- 已实现 Consul、etcd、in-process Registry 和动态 gRPC client pool。
- Scheduler 通过注册中心发现目标服务。

### 阶段五：Gateway 与部署闭环

- Gateway 已按 Registry 中的 HTTP endpoint 代理到 IAM、SYS、Resource。
- Docker Compose 和 Kubernetes 清单只保留微服务部署单元。
- Swagger 由 Gateway 自己提供，不依赖或转发到 IAM。
- 各服务配置只声明自身实际依赖。

## 8. 验收基线

每次修改至少执行：

```bash
cd native
make verify
```

涉及部署、注册或迁移时还必须验证：

1. IAM、SYS、Resource 在注册中心可见且健康。
2. Gateway 能真实转发到目标服务。
3. 三个服务可从空数据库并发完成各自迁移，重复执行幂等。
4. `web-react` 构建通过且不包含后端实现分支。
5. `application/` 下一级仍只有真实进程目录。

## 9. 明确不做

- 聚合业务代码的共享 HTTP 应用。
- 单体装配和双部署形态。
- 为未来需求预建共享业务层或分布式事务。
- HTTP 到 gRPC transcoding。
- Gateway 中的业务聚合、字段裁剪和流程编排。
- 重新引入没有真实消费者的消息队列。
