# Native 配置与环境变量

本文记录 `native` 各进程实际读取的配置和可用环境变量。配置事实以 `internal/config/config.go` 和各服务的 `conf.example.yaml` 为准。

## 加载规则

1. 每个进程都必须通过环境变量显式设置 `LIGHTNING_APP_ENV=dev` 或 `LIGHTNING_APP_ENV=prod`。
2. 程序据此读取当前服务目录的 `conf.dev.yaml` 或 `conf.prod.yaml`，不会读取 `conf.example.yaml`。
3. Git 只管理 `conf.example.yaml`。执行 `make init-config` 可创建缺失的 dev/prod 文件，已有文件不会被覆盖。
4. `conf.<env>.yaml` 提供完整配置；同名且非空的 `LIGHTNING_*` 环境变量覆盖 YAML，适合容器地址和密钥注入。空环境变量不会清空 YAML 值。
5. 配置没有代码默认值。服务拥有的每个键必须在 YAML 或环境变量中出现；允许空值的键必须在 YAML 中显式声明。
6. 环境变量布尔值使用 `true`/`false`，整数使用十进制数字。

`zaplogger.dev.yaml` 和 `zaplogger.prod.yaml` 只从文件加载，目前不支持字段级环境变量覆盖。对日志配置而言，`LIGHTNING_APP_ENV` 只负责选择文件。

## 公共变量

以下变量被多个服务使用，具体归属以各服务章节为准。

| YAML 键 | 环境变量 | 说明 |
| --- | --- | --- |
| `registry.driver` | `LIGHTNING_REGISTRY_DRIVER` | 注册中心驱动：`consul`、`etcd` 或 `inprocess`。跨进程部署使用 Consul 或 etcd。 |
| `registry.address` | `LIGHTNING_REGISTRY_ADDRESS` | 注册中心 HTTP 地址，例如 `http://consul:8500`。 |
| `registry.prefix` | `LIGHTNING_REGISTRY_PREFIX` | etcd 服务注册 key 前缀；仍需为其他驱动显式配置。 |
| `database.dsn` | `LIGHTNING_DATABASE_DSN` | PostgreSQL DSN，包含主机、账号、数据库、端口等连接信息。 |
| `database.maxOpenConns` | `LIGHTNING_DATABASE_MAX_OPEN_CONNS` | 最大打开连接数，必须大于 0。 |
| `database.maxIdleConns` | `LIGHTNING_DATABASE_MAX_IDLE_CONNS` | 最大空闲连接数，范围为 0 到最大打开连接数。 |
| `database.connMaxLifetimeMinutes` | `LIGHTNING_DATABASE_CONN_MAX_LIFETIME_MINUTES` | 单连接最大存活时间，单位为分钟，必须大于 0。 |
| `database.slowThreshold` | `LIGHTNING_DATABASE_SLOW_THRESHOLD` | 慢 SQL 阈值，单位为毫秒，必须大于 0。 |

## Gateway

Gateway 只负责 HTTP/WS 接入、鉴权代理和服务发现，不连接数据库。

| YAML 键 | 环境变量 | 说明 |
| --- | --- | --- |
| `server.port` | `LIGHTNING_SERVER_PORT` | Gateway HTTP 监听端口，范围 1-65535。 |
| `server.tlsCertFile` | `LIGHTNING_SERVER_TLS_CERT_FILE` | TLS 证书路径；不启用 TLS 时显式设为空字符串。 |
| `server.tlsKeyFile` | `LIGHTNING_SERVER_TLS_KEY_FILE` | TLS 私钥路径；必须和证书同时为空或同时非空。 |
| `registry.driver` | `LIGHTNING_REGISTRY_DRIVER` | 用于发现 IAM、SYS、Resource。 |
| `registry.address` | `LIGHTNING_REGISTRY_ADDRESS` | 注册中心地址。 |
| `registry.prefix` | `LIGHTNING_REGISTRY_PREFIX` | 注册中心 key 前缀。 |
| `service.id` | `LIGHTNING_SERVICE_ID` | Gateway 注册中心实例 ID；允许显式留空后自动生成。 |
| `service.advertiseHost` | `LIGHTNING_SERVICE_ADVERTISE_HOST` | Gateway 注册的 HTTP 地址主机名或 IP；Docker 镜像未显式注入时使用当前容器 IP。 |
| `gateway.rateLimitPerMinute` | `LIGHTNING_GATEWAY_RATE_LIMIT_PER_MINUTE` | 单来源 IP 每分钟请求上限；`0` 表示不限制，不能为负数。 |
| `cors.enabled` | `LIGHTNING_CORS_ENABLED` | 是否允许跨域；生产环境必须为 `false`。 |

## IAM

IAM 提供认证授权 HTTP/gRPC 接口，拥有 IAM 数据、Redis 会话、JWT 和 WebSocket 配置。

| YAML 键 | 环境变量 | 说明 |
| --- | --- | --- |
| `server.port` | `LIGHTNING_SERVER_PORT` | IAM HTTP 监听端口。 |
| `grpc.port` | `LIGHTNING_GRPC_PORT` | IAM gRPC 监听端口。 |
| `service.id` | `LIGHTNING_SERVICE_ID` | 注册中心实例 ID；可显式留空，由配置加载器按服务名和主机名生成。 |
| `service.advertiseHost` | `LIGHTNING_SERVICE_ADVERTISE_HOST` | 注册给其他进程访问的主机名或 IP，不是监听地址；Docker 镜像未显式注入时使用当前容器 IP。 |
| `registry.driver` | `LIGHTNING_REGISTRY_DRIVER` | 注册自身并发现 SYS。 |
| `registry.address` | `LIGHTNING_REGISTRY_ADDRESS` | 注册中心地址。 |
| `registry.prefix` | `LIGHTNING_REGISTRY_PREFIX` | 注册中心 key 前缀。 |
| `database.dsn` | `LIGHTNING_DATABASE_DSN` | IAM PostgreSQL DSN。 |
| `database.maxOpenConns` | `LIGHTNING_DATABASE_MAX_OPEN_CONNS` | IAM 数据库最大打开连接数。 |
| `database.maxIdleConns` | `LIGHTNING_DATABASE_MAX_IDLE_CONNS` | IAM 数据库最大空闲连接数。 |
| `database.connMaxLifetimeMinutes` | `LIGHTNING_DATABASE_CONN_MAX_LIFETIME_MINUTES` | 数据库连接最大存活分钟数。 |
| `database.slowThreshold` | `LIGHTNING_DATABASE_SLOW_THRESHOLD` | 慢 SQL 阈值，单位毫秒。 |
| `redis.addr` | `LIGHTNING_REDIS_ADDR` | Redis 地址，例如 `redis:6379`。 |
| `redis.password` | `LIGHTNING_REDIS_PASSWORD` | Redis 密码；Redis 无密码时显式设为空字符串。 |
| `redis.db` | `LIGHTNING_REDIS_DB` | Redis DB 编号，不能为负数。 |
| `jwt.secret` | `LIGHTNING_JWT_SECRET` | JWT 签名密钥，至少 32 个字符；生产环境应从 Secret 注入。 |
| `jwt.expire` | `LIGHTNING_JWT_EXPIRE` | Access Token 有效期，单位秒，必须大于 0。 |
| `auth.tokenHeader` | `LIGHTNING_AUTH_TOKEN_HEADER` | 读取访问令牌的 HTTP Header 名，例如 `Authorization`。 |
| `auth.allowConcurrent` | `LIGHTNING_AUTH_ALLOW_CONCURRENT` | 是否允许同一用户保留多个并发登录会话。 |
| `cors.enabled` | `LIGHTNING_CORS_ENABLED` | 是否允许跨域；生产环境必须为 `false`。 |
| `websocket.enabled` | `LIGHTNING_WEBSOCKET_ENABLED` | 是否启用 IAM WebSocket 模块。 |
| `websocket.timeoutEnabled` | `LIGHTNING_WEBSOCKET_TIMEOUT_ENABLED` | 是否启用 WebSocket 读写超时控制。 |
| `websocket.readTimeoutSeconds` | `LIGHTNING_WEBSOCKET_READ_TIMEOUT_SECONDS` | WebSocket 读超时秒数，必须大于 0。 |
| `websocket.writeTimeoutSeconds` | `LIGHTNING_WEBSOCKET_WRITE_TIMEOUT_SECONDS` | WebSocket 写超时秒数，必须大于 0。 |
| `websocket.heartbeatEnabled` | `LIGHTNING_WEBSOCKET_HEARTBEAT_ENABLED` | 是否启用 WebSocket 心跳。 |
| `websocket.maxReadTimeouts` | `LIGHTNING_WEBSOCKET_MAX_READ_TIMEOUTS` | 连续读超时上限，必须大于 0。 |

## SYS

SYS 提供系统配置和日志 HTTP/gRPC 接口，依赖 IAM 鉴权、PostgreSQL 和 Redis。

| YAML 键 | 环境变量 | 说明 |
| --- | --- | --- |
| `server.port` | `LIGHTNING_SERVER_PORT` | SYS HTTP 监听端口。 |
| `grpc.port` | `LIGHTNING_GRPC_PORT` | SYS gRPC 监听端口。 |
| `service.id` | `LIGHTNING_SERVICE_ID` | 注册中心实例 ID；允许显式留空后自动生成。 |
| `service.advertiseHost` | `LIGHTNING_SERVICE_ADVERTISE_HOST` | 其他进程访问 SYS 的主机名或 IP；Docker 镜像未显式注入时使用当前容器 IP。 |
| `registry.driver` | `LIGHTNING_REGISTRY_DRIVER` | 注册自身并发现 IAM。 |
| `registry.address` | `LIGHTNING_REGISTRY_ADDRESS` | 注册中心地址。 |
| `registry.prefix` | `LIGHTNING_REGISTRY_PREFIX` | 注册中心 key 前缀。 |
| `database.dsn` | `LIGHTNING_DATABASE_DSN` | SYS PostgreSQL DSN。 |
| `database.maxOpenConns` | `LIGHTNING_DATABASE_MAX_OPEN_CONNS` | SYS 数据库最大打开连接数。 |
| `database.maxIdleConns` | `LIGHTNING_DATABASE_MAX_IDLE_CONNS` | SYS 数据库最大空闲连接数。 |
| `database.connMaxLifetimeMinutes` | `LIGHTNING_DATABASE_CONN_MAX_LIFETIME_MINUTES` | 数据库连接最大存活分钟数。 |
| `database.slowThreshold` | `LIGHTNING_DATABASE_SLOW_THRESHOLD` | 慢 SQL 阈值，单位毫秒。 |
| `redis.addr` | `LIGHTNING_REDIS_ADDR` | Redis 地址。 |
| `redis.password` | `LIGHTNING_REDIS_PASSWORD` | Redis 密码；无密码时显式设为空字符串。 |
| `redis.db` | `LIGHTNING_REDIS_DB` | Redis DB 编号，不能为负数。 |
| `auth.tokenHeader` | `LIGHTNING_AUTH_TOKEN_HEADER` | 读取访问令牌的 HTTP Header 名。 |
| `cors.enabled` | `LIGHTNING_CORS_ENABLED` | 是否允许跨域；生产环境必须为 `false`。 |

## Resource

Resource 提供资源 HTTP/gRPC 接口，依赖 IAM 鉴权、PostgreSQL 和 S3 兼容对象存储。

| YAML 键 | 环境变量 | 说明 |
| --- | --- | --- |
| `server.port` | `LIGHTNING_SERVER_PORT` | Resource HTTP 监听端口。 |
| `grpc.port` | `LIGHTNING_GRPC_PORT` | Resource gRPC 监听端口。 |
| `service.id` | `LIGHTNING_SERVICE_ID` | 注册中心实例 ID；允许显式留空后自动生成。 |
| `service.advertiseHost` | `LIGHTNING_SERVICE_ADVERTISE_HOST` | 其他进程访问 Resource 的主机名或 IP；Docker 镜像未显式注入时使用当前容器 IP。 |
| `registry.driver` | `LIGHTNING_REGISTRY_DRIVER` | 注册自身并发现 IAM、SYS。 |
| `registry.address` | `LIGHTNING_REGISTRY_ADDRESS` | 注册中心地址。 |
| `registry.prefix` | `LIGHTNING_REGISTRY_PREFIX` | 注册中心 key 前缀。 |
| `database.dsn` | `LIGHTNING_DATABASE_DSN` | Resource PostgreSQL DSN。 |
| `database.maxOpenConns` | `LIGHTNING_DATABASE_MAX_OPEN_CONNS` | Resource 数据库最大打开连接数。 |
| `database.maxIdleConns` | `LIGHTNING_DATABASE_MAX_IDLE_CONNS` | Resource 数据库最大空闲连接数。 |
| `database.connMaxLifetimeMinutes` | `LIGHTNING_DATABASE_CONN_MAX_LIFETIME_MINUTES` | 数据库连接最大存活分钟数。 |
| `database.slowThreshold` | `LIGHTNING_DATABASE_SLOW_THRESHOLD` | 慢 SQL 阈值，单位毫秒。 |
| `storage.endpoint` | `LIGHTNING_STORAGE_ENDPOINT` | S3 兼容服务地址，不包含协议，例如 `rustfs:9000`。 |
| `storage.accessKey` | `LIGHTNING_STORAGE_ACCESS_KEY` | 对象存储 Access Key；生产环境应从 Secret 注入。 |
| `storage.secretKey` | `LIGHTNING_STORAGE_SECRET_KEY` | 对象存储 Secret Key；生产环境应从 Secret 注入。 |
| `storage.region` | `LIGHTNING_STORAGE_REGION` | 对象存储区域，例如 `us-east-1`。 |
| `storage.bucket` | `LIGHTNING_STORAGE_BUCKET` | 附件 Bucket 名称。 |
| `storage.useSSL` | `LIGHTNING_STORAGE_USE_SSL` | 是否使用 TLS 连接对象存储。 |
| `auth.tokenHeader` | `LIGHTNING_AUTH_TOKEN_HEADER` | 读取访问令牌的 HTTP Header 名。 |
| `cors.enabled` | `LIGHTNING_CORS_ENABLED` | 是否允许跨域；生产环境必须为 `false`。 |

## Scheduler

Scheduler 不监听 HTTP/gRPC，但会注册进程实例，便于观察存活拓扑。它从注册中心发现 SYS、Resource，并使用数据库读取调度配置。

| YAML 键 | 环境变量 | 说明 |
| --- | --- | --- |
| `database.dsn` | `LIGHTNING_DATABASE_DSN` | 调度配置所在的 PostgreSQL DSN。 |
| `database.maxOpenConns` | `LIGHTNING_DATABASE_MAX_OPEN_CONNS` | Scheduler 数据库最大打开连接数。 |
| `database.maxIdleConns` | `LIGHTNING_DATABASE_MAX_IDLE_CONNS` | Scheduler 数据库最大空闲连接数。 |
| `database.connMaxLifetimeMinutes` | `LIGHTNING_DATABASE_CONN_MAX_LIFETIME_MINUTES` | 数据库连接最大存活分钟数。 |
| `database.slowThreshold` | `LIGHTNING_DATABASE_SLOW_THRESHOLD` | 慢 SQL 阈值，单位毫秒。 |
| `registry.driver` | `LIGHTNING_REGISTRY_DRIVER` | 用于发现 SYS、Resource。 |
| `registry.address` | `LIGHTNING_REGISTRY_ADDRESS` | 注册中心地址。 |
| `registry.prefix` | `LIGHTNING_REGISTRY_PREFIX` | 注册中心 key 前缀。 |
| `service.id` | `LIGHTNING_SERVICE_ID` | Scheduler 注册中心实例 ID；允许显式留空后自动生成。 |

Scheduler 没有 `LIGHTNING_SERVER_*`、`LIGHTNING_GRPC_*` 或 `LIGHTNING_SERVICE_ADVERTISE_HOST`，因为它不提供网络 endpoint。

## User Manager

`cmd/usermgr` 是一次性管理工具，不是服务。默认读取 IAM 配置目录，但实际只使用数据库 DSN：

| 来源 | 变量 | 说明 |
| --- | --- | --- |
| 配置选择 | `LIGHTNING_APP_ENV` | 选择 IAM 目录中的 `conf.dev.yaml` 或 `conf.prod.yaml`。 |
| YAML 覆盖 | `LIGHTNING_DATABASE_DSN` | 覆盖工具连接的 PostgreSQL DSN。 |
| 工具输入 | `LIGHTNING_USERMGR_PASSWORD` | 创建或重置的用户密码，至少 8 个字符；只能通过环境变量传入。 |

## 容器注入示例

```yaml
environment:
  LIGHTNING_APP_ENV: prod
  LIGHTNING_SERVER_PORT: "9010"
  LIGHTNING_GRPC_PORT: "9110"
  LIGHTNING_SERVICE_ADVERTISE_HOST: iam
  LIGHTNING_REGISTRY_DRIVER: consul
  LIGHTNING_REGISTRY_ADDRESS: http://consul:8500
  LIGHTNING_DATABASE_DSN: "host=postgres user=postgres password=... dbname=lightning port=5432 sslmode=disable"
  LIGHTNING_REDIS_ADDR: redis:6379
  LIGHTNING_JWT_SECRET: "从 Secret 注入至少 32 个字符的随机值"
```

环境变量只覆盖当前进程使用的配置。不要向 Gateway 注入数据库变量，也不要向 Scheduler 注入 HTTP/gRPC 端口；未归属该进程的变量不会形成有效的服务配置。

## OpenTelemetry 标准变量

Gateway、IAM、SYS、Resource、Scheduler 均初始化 OpenTelemetry Trace。这部分使用 OpenTelemetry SDK 标准环境变量，不属于 `LIGHTNING_*` YAML 配置：

| 环境变量 | 说明 |
| --- | --- |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Collector 通用 OTLP 地址；不设置时不创建 exporter。 |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` | 仅 Trace 使用的 OTLP 地址，优先级高于通用地址。 |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | OTLP 协议；当前 exporter 使用 gRPC。 |
| `OTEL_EXPORTER_OTLP_HEADERS` | Collector 要求认证时使用的请求 Header。 |
| `OTEL_TRACES_SAMPLER` | Trace 采样器，例如 `parentbased_traceidratio`。 |
| `OTEL_TRACES_SAMPLER_ARG` | 比例采样参数，例如 `0.1`。 |
| `OTEL_TRACES_EXPORTER` | 设置为 `none` 时强制关闭 Trace 导出。 |

服务名、实例 ID 和 `dev`/`prod` 环境由各进程根据已有服务配置写入 OTel Resource，不通过 `OTEL_SERVICE_NAME` 重复配置。完整插桩边界和业务 Span 写法见 [`opentelemetry.md`](opentelemetry.md)。
