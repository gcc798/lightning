# AGENTS.md

本文件供 Codex、Claude Code 和其他代码代理快速理解仓库结构与开发约定。

## 项目定位

`microservice-kit` 是一个渐进式微服务后台脚手架 Monorepo。它不是单一后端工程，而是同一套业务语义的三套 Go 实现，加上一个 React 前端工程：

- `native/`：不依赖完整微服务框架的原生 Go 渐进式微服务实现，是主要演进方向、业务语义和接口行为基线。
- `kratos/`：基于 Kratos 开源微服务框架的跟进实现，开发进度可能阶段性落后于 `native/`。
- `gozero/`：基于 go-zero 开源微服务框架的跟进实现，开发进度可能阶段性落后于 `native/`。
- `web-react/`：唯一的前端工程，应尽量无缝对接三套后端实现。

“渐进式”是指先把需求归入现有领域服务，只有出现清晰、独立的数据所有权和部署需求时才新增服务；它不表示同时维护单体和微服务两种运行形态。`native` 当前只保留真实微服务拓扑。做需求或修 bug 时，先理解 `native` 的业务语义，再把同等行为落到目标框架实现中。Kratos 和 go-zero 允许阶段性滞后，但不得以滞后实现反向定义或限制 `native` 的业务语义。

## 核心原则

1. `native/` 是业务基线。接口路径、请求参数、响应结构、错误语义和业务规则优先参考它。
2. `kratos/` 和 `gozero/` 是框架化实现，不应发明不同的业务语义。
3. `web-react/` 是同一套前端，通过 HTTP 契约与后端交互；不应为了不同后端实现让前端做兼容分支。
4. 三套后端可以有不同工程分层和框架代码生成方式，但对外业务能力应保持一致。
5. 只要 `web-react` 能和 `native` 正常交互，`kratos` 和 `gozero` 也必须提供兼容的 HTTP 契约。
6. 修改生成文件时优先使用对应框架命令重新生成，不要只做字符串硬改。
7. 本仓库是全新的试验工程，任何改动默认不考虑旧代码、旧配置、旧接口或旧数据的向后兼容；只有开发者明确提出兼容要求时，才实现兼容逻辑。
8. `native/` 只保留微服务部署形态。`application/` 下一级目录必须对应真实进程；每个领域服务在自己的 `internal/` 中拥有业务代码和数据模型，根 `internal/` 只放跨进程通用技术设施。不得重新引入共享 API 应用或单体装配入口。
9. 渐进拆分沿事务、数据所有权和独立部署边界进行，不按 controller 数量拆分，也不为假设中的未来需求预建服务、共享业务层或分布式事务。

## 根目录结构

```text
microservice-kit/
├── native/      # 原生 Go 渐进式微服务实现，业务与 HTTP 契约基线
├── kratos/      # Kratos 微服务框架版本
├── gozero/      # go-zero 微服务框架版本
├── web-react/   # React 前端
├── README.md
├── AGENTS.md
└── CLAUDE.md
```

## 后端实现关系

### native

`native/` 是当前主要演进的后端实现。它用 Go 标准能力和小型基础库组合出可独立部署的微服务能力，同时避免绑定完整微服务框架。

常见入口与边界：

- `native/application/gateway/`：唯一对外 HTTP/WS 入口，负责服务发现和反向代理，不拥有业务数据。
- `native/application/{iam,sys,resource}/internal/`：各领域服务私有的 controller、router、DTO、domain、model 和 Goose 迁移。
- `native/application/scheduler/`：Scheduler 进程与异步 Job，不执行数据库迁移。
- `native/cmd/usermgr/`：管理员维护工具，不是 `application/` 部署单元。
- `native/api/<service>/v1/`：Proto 源文件。
- `native/internal/api/<service>/v1/`：生成的 gRPC 契约和薄适配。
- `native/internal/httpserver/`：共享 Echo 生命周期，不注册任何业务路由。
- `native/internal/transport/`：gRPC server、注册生命周期和动态 client pool。
- `native/internal/registry/`：Consul、etcd、in-process 注册发现与实例选择。
- `native/internal/openapi/`：Gateway 提供的统一 Swagger 生成文件。
- `native/internal/database/`：共享数据库连接与 GORM 插件，不负责执行迁移。
- `native/internal/`：多个进程共享、但通过 Go `internal` 规则禁止仓库外导入的技术能力；不得放领域业务实现。

每个进程在自己的目录版本化 `conf.example.yaml` 和 `zaplogger.example.yaml`；开发者通过 `make init-config` 基于模板创建被 Git 忽略的 `*.dev.yaml` 和 `*.prod.yaml`。程序只按 `MS_K_APP_ENV=dev|prod` 读取对应运行配置，不直接读取 example 文件；部署环境变量可覆盖其中的地址和敏感值。
每个服务的 YAML 模板只声明自身实际依赖，不得为了复用配置结构加入未使用字段。

IAM、SYS、Resource 分别拥有 `application/<service>/internal/migrations/sql/`，启动时只执行自己的迁移，并使用独立的 `goose_<service>_version` 表。Scheduler 不执行迁移。

`native` 不使用顶层 `pkg/` 存放普通共享代码。仅当某个包明确作为稳定 API 供当前 Go Module 之外的工程导入时，才考虑新增 `pkg/`；仓库内多应用共享代码应保留在 `internal/`。

权限模型长期约束：

- API 权限是后端安全边界，菜单只负责前端路由和按钮显隐，两者不得合并为一棵权限树。
- `s_menu.perms` 不参与后端鉴权；菜单通过 `m_menu_api_permission` 关联一个或多个 API 权限。
- 角色和用户 API 权限的 `source=0` 表示手工授权，`source=1` 表示菜单派生授权；保存或撤销任一来源不得覆盖另一来源。
- 路由必须显式声明 `resource` 和 `action`，不得根据权限名称或路径猜测 action。
- 当前 SQL 权限查询已满足需求，不引入 Casbin；出现不绑定用户的真实机器客户端后再设计 `client_credentials` 权限。

常用命令：

```bash
cd native
make verify
```

### kratos

`kratos/` 是基于 Kratos 的微服务实现，采用 `sys-api` + `sys-rpc` 分层。

常见入口：

- `kratos/api/system/v1/`：proto 契约。
- `kratos/application/sys-api/`：对外 HTTP 服务。
- `kratos/application/sys-rpc/`：内部 gRPC 服务和核心数据访问。
- `kratos/application/sys-rpc/ent/schema/`：Ent schema。
- `kratos/pkg/`：Kratos 版本沉淀的通用能力。
- `kratos/Makefile`：代码生成、构建和测试入口。

常用命令：

```bash
cd kratos
make conf
make proto
make ent
make wire
go test ./...
```

注意：proto 相关 Go 代码包含编码后的 descriptor。修改 `go_package` 或 module 路径后，要用 `make conf` / `make proto` 重新生成，不能只做文本替换。

### gozero

`gozero/` 是基于 go-zero 的微服务实现，保留 `sys-api` / `sys-rpc` 分层。

常见入口：

- `gozero/application/sys-api/`：对外 API 服务。
- `gozero/application/sys-rpc/`：内部 RPC 服务。
- `gozero/application/sys-api/sys.api`：API 描述文件。
- `gozero/Makefile`：构建入口。

常用命令：

```bash
cd gozero
go test ./...
```

## 前端工程

### web-react

`web-react/` 是唯一的前端工程。它只应依赖统一 HTTP 契约，不应因为后端实现不同而写兼容逻辑或分支判断。

前端契约原则：

- 前端只面向一套接口语义。
- `native` 是前端契约基线。
- 如果 `web-react` 可以和 `native` 正常交互，`kratos` 和 `gozero` 也应该无缝可用。
- 发现某个框架后端无法被同一套前端调用时，优先修后端契约，而不是改前端兼容。

常用命令：

```bash
cd web-react
pnpm install
pnpm dev
pnpm build
```

默认开发服务端口参考 `web-react/README.md`。

## 开发建议

- 做业务改动前，先确定 IAM、SYS 或 Resource 的所有权，再在对应的 `application/<service>/internal/` 中修改 controller/domain/model/request/response。
- 若目标是 `kratos` 或 `gozero`，实现时保持与 `native` 的 HTTP 契约和接口语义一致。
- 若改动会影响前端接口，优先确认是否破坏了 `native` 契约；不要让 `web-react` 为不同后端实现做特殊兼容。
- 每个 Go 子工程单独运行测试：`native`、`kratos`、`gozero` 各自都有自己的 `go.mod`。
- 不要把 `native`、`kratos`、`gozero` 当成互相引用的包；它们是同一业务的不同实现。
- 新需求默认归入 IAM、SYS 或 Resource；只有确认独立数据所有权和部署需求后，才新增 `application/<service>`。
- 不要手工修改 `native/internal/openapi/`，应通过 `make swagger` 重新生成。
- 不要在根 `native/internal/` 创建 IAM、SYS 或 Resource 的业务包，也不要创建聚合业务 HTTP 代码的 `application/api`。
- 每个 `native/application/<service>/main.go` 显式加载配置、初始化日志并构建自身依赖；不要增加 `app.Base` 一类只转发启动步骤的包装层。无测试复用需求时，启动流程直接写在 `main` 中，不额外封装 `run`。

## 快速判断应该看哪里

- 想知道业务原始行为：看 `native/`。
- 想改 Kratos 微服务版本：看 `kratos/api` 和 `kratos/application`。
- 想改 go-zero 微服务版本：看 `gozero/application`。
- 想看前端当前对接方式：看 `web-react/src`。
