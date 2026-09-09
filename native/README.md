# lightning native

`native` 是 lightning 的原生 Go 后端，也是三套后端实现的业务与 HTTP 契约基线。

## 本地启动

依赖 Go 1.26.5 和 PostgreSQL 16；IAM/SYS 另外使用 Redis 7。先从 Git 托管的模板创建本地配置，再选择运行环境：

```bash
cd native
make init-config
export LIGHTNING_APP_ENV=dev
```

每个进程的配置位于自己的 `application/<service>/` 目录。修改对应服务的 `conf.dev.yaml` 后，分别启动微服务：

```bash
go run ./application/iam
```

网关负责统一 HTTP/WS 接入，IAM/SYS/Resource 各自提供 HTTP 与 gRPC；Scheduler、定时任务和常驻后台协程由独立进程运行：

```bash
cd native
go run ./application/scheduler
```

生产环境只提供 `gateway + iam + sys + resource` 微服务部署；Scheduler 默认运行一个副本。

本地启动微服务时先启动 PostgreSQL、Redis、RustFS 与 Consul，然后分别运行以下进程。默认 HTTP 端口为 gateway `9009`、iam `9010`、sys `9011`、resource `9012`，gRPC 端口依次为 `9110`、`9111`、`9112`。

```bash
go run ./application/iam
go run ./application/sys
go run ./application/resource
go run ./application/gateway
go run ./application/scheduler
```

服务间契约在 `api/<domain>/v1/*.proto`；修改后运行 `make proto`，`make verify` 会检查生成文件是否最新。默认注册中心为 Consul `http://127.0.0.1:8500`。切换 etcd v3 HTTP gateway 时设置 `LIGHTNING_REGISTRY_DRIVER=etcd` 和 `LIGHTNING_REGISTRY_ADDRESS=http://127.0.0.1:2379`；注册 key 前缀可用 `LIGHTNING_REGISTRY_PREFIX` 调整。

`LIGHTNING_APP_ENV` 必须显式设置为 `dev` 或 `prod`，用于在具体服务目录中选择 `conf.dev.yaml` 或 `conf.prod.yaml`。程序不会读取 `*.example.yaml`。配置没有代码默认值；每个服务使用的键必须在 YAML 或对应的 `LIGHTNING_*` 环境变量中显式出现，环境变量优先于 YAML。

各服务支持的环境变量、对应 YAML 键和约束见 [`docs/configuration.md`](docs/configuration.md)。日志配置只从 `zaplogger.<env>.yaml` 读取，不支持字段级环境变量覆盖。

Git 只管理每个服务的 `conf.example.yaml` 和 `zaplogger.example.yaml`。`make init-config` 会在文件不存在时把模板分别复制为 `*.dev.yaml` 和 `*.prod.yaml`，不会覆盖已有配置；这些实际运行配置已被 Git 忽略。模板只声明该进程实际使用的配置段，例如 Gateway 只配置接入与注册中心，Resource 配置数据库和对象存储，Scheduler 配置数据库和注册中心。创建后应按环境修改地址和凭据；生产敏感值也可以通过对应的 `LIGHTNING_*` 环境变量注入。

IAM、SYS、Resource 的 `service.id` 是注册中心中的实例唯一标识；留空时程序按“服务名 + 主机名 + HTTP 端口”自动生成，只有需要固定实例 ID 时才填写。`service.advertiseHost` 是写入注册中心、供 Gateway 和其他服务访问该实例的地址；本机开发使用 `127.0.0.1`，Docker Compose 使用服务名 `iam`、`sys`、`resource`，Kubernetes 使用可被其他 Pod 解析的 Service DNS 或 Pod 地址。它不是监听地址，HTTP/gRPC 仍由 `server.port` 和 `grpc.port` 监听。

生产环境默认关闭 CORS，前后端通过 Nginx 同源代理；开发环境可在 `conf.dev.yaml` 中开启。

验证码、微信、短信、邮件和 Scheduler 属于运行期模块配置，不在 YAML 或环境变量中定义。它们持久化在 `s_config`；IAM 的运行期模块通过 Redis 共享读取，缓存缺失时在分布式锁内从数据库加载并回填。单实例 Scheduler 直接周期读取数据库中的调度配置，不依赖 Redis。每个配置编码只保留一条最新记录，不使用配置版本号和 Redis Pub/Sub。

## 数据库迁移与初始管理员

IAM、SYS、Resource 启动时分别执行自己嵌入二进制的 Goose 迁移，SQL 位于 `application/<service>/internal/migrations/sql/`，版本表分别为 `goose_iam_version`、`goose_sys_version`、`goose_resource_version`。数据库结构不使用 GORM AutoMigrate；Scheduler 不执行迁移。

迁移会创建逻辑客户端 `web-admin` 和内置角色，但不会写入默认管理员或默认密码。创建首个管理员：

```bash
export LIGHTNING_USERMGR_PASSWORD='replace-with-a-strong-password'
go run ./cmd/usermgr --operation=create --username=admin --nickname=管理员 --role=super_admin
```

重置密码：

```bash
export LIGHTNING_USERMGR_PASSWORD='replace-with-a-new-strong-password'
go run ./cmd/usermgr --operation=reset --username=admin
```

密码不支持命令行参数，工具也不会回显密码。

## 认证

登录接口为 `POST /login`，客户端 ID 为 `web-admin`，支持 `password`、`email`、`sms`、`wechat`。微信登录只接收小程序 `wxCode`，服务端通过微信接口换取 OpenID/UnionID，不接收客户端提交的 OpenID。

Access Token 与 Refresh Token 都关联服务端 Redis 会话；刷新令牌单次使用并在刷新后轮换，登出会立即撤销当前会话。

## 质量检查

```bash
make verify # go test、go vet、race、Swagger freshness
make ci     # verify + Docker build
```

`make build` 构建 gateway、iam、sys、resource 和 scheduler 五个进程。

涉及通用工具、配置、认证基础设施和中间件的改动必须补充单元测试。controller 与 `application/<service>/internal/domain` 中的具体业务逻辑不强制单测，可按风险补充集成或契约测试。

## Docker

```bash
make init-config
export LIGHTNING_JWT_SECRET='replace-with-at-least-32-random-characters'
docker compose up -d --build
```

所有进程使用根目录唯一的 `Dockerfile`，通过 `TARGET=gateway|iam|sys|resource|scheduler` 选择构建入口。例如：

```bash
docker build --build-arg TARGET=scheduler -t lightning-native-scheduler .
```

Compose 默认启动 Consul、gateway、iam、sys、resource、scheduler、PostgreSQL、Redis 和 RustFS，即完整微服务形态。Consul UI/API 位于 `8500`，所有前端 HTTP 请求统一进入 gateway 的 `9009`。

PostgreSQL、Redis 和 RustFS 带有本地默认值；需要覆盖时使用 `LIGHTNING_POSTGRES_PASSWORD`、`LIGHTNING_REDIS_PASSWORD` 和 `LIGHTNING_RUSTFS_*`。RustFS 提供 S3 兼容对象存储；resource 启动时会检查并按需创建 `lightning` Bucket。

Kubernetes 清单位于 `k8s/`，包含上述五个服务进程、Scheduler、持久化单节点 Consul、ClusterIP、Ingress 和 gateway HPA。部署前需替换镜像名并基于 `k8s/secret.yaml.example` 创建 Secret；PostgreSQL、Redis 和 RustFS 服务地址按现有 ConfigMap 接入。

当前 WebSocket Hub 位于 IAM 进程内，因此 IAM 清单固定为一个副本；gateway 可独立水平扩容。需要 IAM 多副本时，先增加 Redis Pub/Sub 跨实例广播，再调整 IAM 副本数。

## 代码边界

```text
application/gateway    HTTP/WS 反向代理、鉴权和服务发现入口
application/{iam,sys,resource}  独立领域服务入口及各自 internal 业务代码
application/scheduler  Scheduler 入口与 Job 定义
cmd/usermgr            一次性管理工具
api                    Proto 源文件
internal/api           gRPC 契约生成代码与薄适配
internal/httpserver    共享 HTTP 生命周期，不包含业务路由
internal/registry      in-process、Consul、etcd 注册发现
internal/transport     gRPC server 与动态 client 连接池
internal               多个进程共享、但禁止仓库外导入的技术能力；不放领域业务代码
```

项目不使用顶层 `pkg`。当前共享代码并不是提供给其他 Go Module 使用的公共 SDK，因此使用 Go 编译器强制约束的 `internal` 更准确；只有未来出现明确、稳定且需要被仓库外工程导入的 API 时，才新增 `pkg`。
