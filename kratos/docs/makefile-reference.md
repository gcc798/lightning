# Makefile 说明文档

这份文档用于说明当前 `kratos` 工程中的 `Makefile`：

- 每个变量表示什么
- 每个 target 负责什么工作
- 为什么生成链路被拆成多段
- 平时开发时应该怎么使用

主文件：
- [../Makefile](../Makefile)

---

## 一、变量说明

### `KRATOS_ROOT`
- 含义：当前 `kratos` 工程根目录
- 作用：其他路径变量都基于它展开
- 目的：让 Makefile 不依赖手写绝对路径，保持工程内部可迁移

### `API_DIR`
- 值：`$(KRATOS_ROOT)/api`
- 含义：业务 proto 的根目录
- 作用：`protoc` 在这里作为 import root 工作
- 结果：proto 里可以写这种导入：
  - `import "system/v1/auth.proto"`
- 而不需要写：
  - `import "api/system/v1/auth.proto"`

### `THIRD_PARTY_DIR`
- 值：`$(KRATOS_ROOT)/third_party`
- 含义：外部 proto 依赖目录
- 当前主要包含：
  - `google/api/*.proto`
  - `google/protobuf/*.proto`
  - `validate/validate.proto`
- 作用：
  - 给 `protoc` 提供外部导包
  - 给 `kratos proto client` 提供依赖
  - 给 IDE 提供稳定的本地索引目录

### `SYS_API_DIR`
- 值：`$(KRATOS_ROOT)/application/sys-api`
- 含义：`sys-api` 服务根目录

### `SYS_RPC_DIR`
- 值：`$(KRATOS_ROOT)/application/sys-rpc`
- 含义：`sys-rpc` 服务根目录

### `SYS_API_SERVICE_GEN_DIR`
- 值：`$(SYS_API_DIR)/internal/servicegen`
- 含义：`sys-api` 的生成骨架目录
- 作用：存放 `kratos proto server` 生成的 skeleton 文件
- 为什么不直接生成到 `internal/service`：
  - 避免反复生成时覆盖真实实现
  - 保留骨架参考，但真实代码仍然写在 `internal/service`

### `SYS_RPC_SERVICE_GEN_DIR`
- 值：`$(SYS_RPC_DIR)/internal/servicegen`
- 含义：`sys-rpc` 的生成骨架目录
- 逻辑和 `SYS_API_SERVICE_GEN_DIR` 一样

### `PROTO_FILES`
- 含义：`api/` 目录下全部业务 proto 文件
- 生成方式：
  1. 进入 `api/`
  2. 找出所有 `*.proto`
  3. 转成相对于 `api/` 的路径
  4. 排序
- 结果示例：
  - `system/v1/auth.proto`
  - `system/v1/user.proto`
- 作用：作为 `proto`、`openapi`、`proto-client`、`proto-server-*` 的统一输入集合

### `CONF_PROTO_FILES`
- 含义：配置 proto 文件列表
- 当前包含：
  - `application/sys-api/internal/conf/conf.proto`
  - `application/sys-rpc/internal/conf/conf.proto`
- 作用：配置体系单独生成，不和业务 proto 混在一起

### `PROTOBUF_INCLUDE`
- 含义：本机 protobuf 标准 include 目录
- 作用：从这里复制标准 protobuf 文件到 `third_party/google/protobuf`

### `PATH`
- 作用：优先使用本地工具目录 `/Users/guoc/dev/code_go/bin`
- 里面通常会放：
  - `kratos`
  - `protoc-gen-openapi`
  - 其他生成工具

### `GOCACHE`
- 值：`/tmp/lightning-kratos-go-build`
- 作用：给当前项目单独指定 Go 编译缓存目录
- 好处：
  - 不污染仓库
  - 减少重复编译开销
  - 更容易排查构建问题

---

## 二、Target 说明

### `help`
- 打印当前支持的 Make target 列表
- 适合作为快速入口

### `init`
- 执行：
  - `prepare-third-party`
  - `go mod tidy`
- 作用：初始化一个新 checkout 的工程环境

### `prepare-third-party`
- 创建 `third_party` 目录结构
- 从本地依赖里复制外部 proto 文件
- 当前复制来源包括：
  - `grpc-gateway` 的 `google/api` proto
  - 本机 protobuf include 里的标准 proto
  - `protoc-gen-validate` 的 proto
- 作用：
  - 保证 proto 生成时不依赖远程资源
  - 保证 IDE 可以稳定识别导包

### `clean-proto`
- 清理业务 proto 生成物：
  - `*.pb.go`
  - `*_grpc.pb.go`
  - `*_http.pb.go`
  - `*_pb.validate.go`
  - `*_client.go`
- 清理 OpenAPI 生成物：
  - 所有 `openapi.yaml`
  - 所有 `*.openapi.yaml`
- 也会顺手删掉历史残留的根目录 `openapi.yaml`
- 作用：确保重新生成时环境干净

### `clean-conf`
- 清理配置 proto 生成物：
  - `application/sys-api/internal/conf/conf.pb.go`
  - `application/sys-rpc/internal/conf/conf.pb.go`

### `conf`
- 用 `protoc` 单独生成配置 protobuf 代码
- 输出：
  - `sys-api/internal/conf/conf.pb.go`
  - `sys-rpc/internal/conf/conf.pb.go`
- 为什么单独做：
  - 配置 proto 不属于 `api/` 下的业务契约
  - 配置和业务 proto 的生成场景不同

### `proto-add-health`
- 在 `api/` 目录中执行 `kratos proto add`
- 用途：快速新增一个标准 Kratos proto 骨架
- 当前示例：
  - `system/v1/health.proto`

### `proto`
- 在 `api/` 目录下执行 `protoc`
- include root：
  - `-I .`，即 `api/` 自身
  - `-I third_party`
- 生成内容：
  - protobuf message
  - gRPC 代码
  - HTTP binding 代码
  - error 代码
  - validate 代码
- 这是业务契约生成链路的核心步骤

### `openapi`
- 当前不是直接“一次生成一个总文档”
- 实际流程是：
  1. 每个 proto 先各自生成一份临时 openapi
  2. 再按服务目录做聚合
  3. 输出成每个服务/版本一份文档
- 当前输出规则：
  - `api/<service>/<version>.openapi.yaml`
- 当前 system 输出示例：
  - [../api/system/v1.openapi.yaml](../api/system/v1.openapi.yaml)
- 为什么需要聚合：
  - `protoc-gen-openapi` 直接对我们当前这套 proto 目录执行时，不能稳定产出完整的 service 级单文档
  - 之前会只保留最后一个 proto 的结果
- 所以现在引入了聚合工具：
  - [../cmd/openapi-merge/main.go](../cmd/openapi-merge/main.go)

### `proto-all`
- 全量生成入口
- 当前执行顺序：
  1. `conf`
  2. `clean-proto`
  3. `prepare-third-party`
  4. `proto`
  5. `openapi`
  6. `proto-client`
  7. `proto-server-sys-api`
  8. `proto-server-sys-rpc`
- 这是修改 proto 之后最常用的一条命令

### `proto-client`
- 对每个业务 proto 执行 `kratos proto client`
- 作用：生成 Kratos 风格的 client 代码
- 为什么它和 `proto` 分开：
  - `proto` 负责标准 protobuf/grpc/http 生成
  - `kratos proto client` 负责 Kratos 的 client 封装

### `proto-server-sys-api`
- 对每个业务 proto 执行 `kratos proto server`
- 输出目录：
  - `application/sys-api/internal/servicegen`
- 作用：
  - 生成骨架
  - 供开发时参考
  - 不覆盖真实 service 实现

### `proto-server-sys-rpc`
- 同 `proto-server-sys-api`
- 输出目录：
  - `application/sys-rpc/internal/servicegen`

### `ent`
- 执行 Ent schema 生成
- 当前只针对：
  - `application/sys-rpc/ent/schema`
- 原因：
  - `sys-rpc` 负责数据库访问
  - `sys-api` 不负责持久化模型生成

### `wire`
- 执行 Wire 生成
- 当前入口：
  - `application/sys-api/cmd/server`
  - `application/sys-rpc/cmd/server`
- 输出：
  - 两个服务各自的 `wire_gen.go`
- 作用：生成依赖注入装配代码

### `fmt`
- 执行 `go fmt ./...`
- 作用：统一格式

### `build-sys-api`
- 构建 `sys-api` 二进制
- 输出位置：
  - `bin/sys-api`

### `build-sys-rpc`
- 构建 `sys-rpc` 二进制
- 输出位置：
  - `bin/sys-rpc`

### `build-all`
- 同时构建两个二进制

### `test`
- 执行 `go test ./...`
- 用于全项目验证

---

## 三、为什么不写成一条超长命令

当前 Makefile 是有意拆开的，不是随意分散。

### 1. 配置生成和业务 proto 生成是两回事
因为：
- 配置 proto 在 `application/.../internal/conf`
- 业务 proto 在 `api/...`
- import root、输出语义、使用场景都不同

### 2. `protoc` 和 `kratos proto` 解决的问题不同
- `protoc`：
  - protobuf
  - grpc
  - http
  - validate
  - error
- `kratos proto client`：
  - 生成 Kratos client
- `kratos proto server`：
  - 生成 service skeleton

### 3. OpenAPI 需要独立聚合
因为当前生成器不能直接给出我们想要的“每个微服务一份完整文档”，所以要先拆、再合。

---

## 四、当前 Makefile 反映出的工程约定

### 1. `servicegen` 长期保留
这是已经确定的约定：
- `internal/servicegen` 保留生成骨架
- `internal/service` 保留真实实现
- 不让生成覆盖开发者代码

### 2. proto import root 统一是 `api/`
所以内部导包现在写成：
- `system/v1/auth.proto`
- `system/v1/user.proto`

这比把仓库目录名 `api/` 暴露进契约路径更干净。

### 3. OpenAPI 输出按服务/版本组织
当前规则：
- proto 源码目录：`api/system/v1/`
- OpenAPI 输出：`api/system/v1.openapi.yaml`

这样做的好处：
- 源码目录和生成产物分开
- 多版本更清楚
- 多微服务扩展也更自然

---

## 五、推荐使用方式

### 修改 proto 之后
建议执行：
```bash
make proto-all
make wire
make test
```

### 修改 Ent schema 之后
建议执行：
```bash
make ent
make test
```

### 修改依赖注入之后
建议执行：
```bash
make wire
make test
```

### 提交前
建议执行：
```bash
make fmt
make test
make build-all
```

---

## 六、最常受影响的文件区域

### Proto 相关
- [../api/system/v1](../api/system/v1)
- [../application/sys-api/internal/servicegen](../application/sys-api/internal/servicegen)
- [../application/sys-rpc/internal/servicegen](../application/sys-rpc/internal/servicegen)

### 配置相关
- [../application/sys-api/internal/conf](../application/sys-api/internal/conf)
- [../application/sys-rpc/internal/conf](../application/sys-rpc/internal/conf)

### OpenAPI 相关
- [../api/system/v1.openapi.yaml](../api/system/v1.openapi.yaml)
- [../cmd/openapi-merge/main.go](../cmd/openapi-merge/main.go)

### Wire 相关
- [../application/sys-api/cmd/server/wire.go](../application/sys-api/cmd/server/wire.go)
- [../application/sys-api/cmd/server/wire_gen.go](../application/sys-api/cmd/server/wire_gen.go)
- [../application/sys-rpc/cmd/server/wire.go](../application/sys-rpc/cmd/server/wire.go)
- [../application/sys-rpc/cmd/server/wire_gen.go](../application/sys-rpc/cmd/server/wire_gen.go)
