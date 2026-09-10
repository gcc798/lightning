# Native OpenTelemetry 链路追踪

Native 使用 OpenTelemetry 生成和传播 Trace。应用未配置 OTLP 接收端时仍会生成 `trace_id`、`span_id` 并写入关联日志，但不会导出或持久化 Span，也不会影响业务启动。

## 自动插桩边界

- Gateway HTTP 入口与反向代理：自动创建并传播 HTTP Span。
- IAM、SYS、Resource HTTP 入口：Echo 中间件自动提取 W3C `traceparent` 并创建 Server Span。
- gRPC：`internal/transport` 为全部客户端和服务端自动创建、传播 Span。
- PostgreSQL：`internal/container` 在 `database/sql` 驱动层创建查询 Span；业务查询必须继续使用 `db.WithContext(ctx)`。
- Redis：`internal/container` 注册 go-redis tracing hook；业务调用必须传入当前 `ctx`。
- Scheduler：调度器为每次任务执行创建 `scheduler.job` Root Span，并把 `ctx` 传入 Job。
- `/health/*` 和 `/metrics` 不创建 HTTP Span，避免健康检查噪声。

业务 controller、domain 和 model 不需要为 HTTP、gRPC、PostgreSQL、Redis 重复创建 Span。不得使用 `context.Background()` 替换入口传下来的 `ctx`。

## 业务 Span 示例

只有框架 Span 无法表达、且排障确实关心的业务阶段才手工创建 Span。登录流程 `application/iam/internal/domain/auth.go` 中的 `auth.login` 是参考实现：

```go
ctx, span := otel.Tracer("github.com/gcc798/microservice-kit/application/iam").Start(ctx, "auth.login")
defer span.End()
```

错误需要同时记录到 Span：

```go
span.RecordError(err)
span.SetStatus(codes.Error, err.Error())
```

Span 属性只记录排障需要的非敏感信息。禁止记录密码、Token、验证码、Cookie、完整请求体或 SQL 参数。

## 日志关联

需要关联当前调用链的日志通过 `logger.WithContext` 输出：

```go
logging.WithContext(ctx, log).Info("login succeeded", zap.Int64("user_id", userID))
```

有效 Span 上下文会产生：

```json
{"msg":"login succeeded","trace_id":"4bf92f3577b34da6a3ce929d0e0e4736","span_id":"00f067aa0ba902b7"}
```

`trace_id` 标识完整的跨服务请求，`span_id` 标识其中一个具体操作。HTTP 访问日志、登录示例日志和 Scheduler Job 日志已接入该关联方式。

## 可选 OTLP 导出

不部署采集组件时无需设置任何 `OTEL_*` 变量。需要导出到 OpenTelemetry Collector 时，至少显式设置：

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=http://127.0.0.1:4317
export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
```

容器中通常把 Endpoint 设置为 Collector 的 Compose 服务地址。SDK 使用批量异步导出，进程退出时在 `main` 中执行 flush。可用 `OTEL_TRACES_EXPORTER=none` 强制关闭导出；Native 只有检测到 `OTEL_EXPORTER_OTLP_ENDPOINT` 或 `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` 时才创建 exporter。

部署 Collector、Tempo、Loki、Grafana 不属于业务源码要求，可在需要集中存储和展示时独立增加。
