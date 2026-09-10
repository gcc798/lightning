# Native Proto/gRPC 开发约定

本文是 `native` 的 Proto 契约说明，供开发者和代码 agent 修改内部 gRPC 接口时使用。

## 目录边界

- `api/<service>/v1/*.proto`：手写的稳定契约，只放服务对外提供的 gRPC service、RPC、message 和字段。
- `internal/api/<service>/v1/*.pb.go`：`protoc-gen-go` 生成的 message 代码。
- `internal/api/<service>/v1/*_grpc.pb.go`：`protoc-gen-go-grpc` 生成的 client、server interface 和注册描述。
- `internal/api/<service>/v1/contracts.go`：手写的 Go 适配层。`Remote` 通过服务发现获取连接，`API` 是业务代码依赖的最小接口。
- `application/<service>/internal/domain/rpc.go`：服务端把领域 API 适配为生成的 gRPC server interface。

当前契约归属如下：

| Proto | 服务名 | 实现位置 |
| --- | --- | --- |
| `api/iam/v1/iam.proto` | `iam` | `application/iam/internal/domain/rpc.go` |
| `api/sys/v1/sys.proto` | `sys` | `application/sys/internal/domain/rpc.go` |
| `api/resource/v1/resource.proto` | `resource` | `application/resource/internal/domain/rpc.go` |

## 编写规则

1. 文件使用 `syntax = "proto3"`，包名使用 `microservice-kit.<service>.v1`。
2. `go_package` 必须指向 `github.com/gcc798/microservice-kit/internal/api/<service>/v1`，并使用稳定的别名，例如 `iamv1`。
3. 每个 `service`、每个 `rpc`、每个 `message` 以及每个 message 字段都必须紧邻写一行 `//` 注释。注释说明业务含义；有单位、状态值或约束时一并写明。
4. 字段编号一旦发布就不能复用或改变含义。删除字段使用 `reserved` 保留编号和名称。
5. 优先新增字段保持向后兼容；不要为了重命名、换类型或调整语义直接复用旧编号。破坏性变更应新建 `v2` 包。
6. Proto 只描述跨进程契约，不放数据库模型、HTTP DTO 或实现细节。
7. 不要手工编辑 `*.pb.go` 和 `*_grpc.pb.go`。

示例：

```proto
// IAMService 提供身份认证和权限校验的内部 gRPC 能力。
service IAMService {
  // CheckPermission 校验用户是否拥有指定资源和动作的权限。
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
}

// CheckPermissionRequest 是权限校验请求。
message CheckPermissionRequest {
  // user_id 是待校验权限的用户 ID。
  int64 user_id = 1;
}
```

## 生成代码

在 `native` 目录执行：

```bash
make proto
```

该命令等价于：

```bash
protoc -I . \
  --go_out=. --go_opt=module=github.com/gcc798/microservice-kit \
  --go-grpc_out=. --go-grpc_opt=module=github.com/gcc798/microservice-kit \
  api/iam/v1/iam.proto api/sys/v1/sys.proto api/resource/v1/resource.proto
```

生成文件落在 `native/internal/api/<service>/v1/`。新增或修改 Proto 文件时，要同步把文件加入 `native/Makefile` 的 `proto` 和 `proto-check` 目标；生成器版本以本机安装的 `protoc`、`protoc-gen-go`、`protoc-gen-go-grpc` 为准。

## 服务端注册

服务入口先创建注册中心管理的 gRPC server，再注册生成的实现：

```go
grpcServer, err := transport.StartRegisteredGRPC(ctx, cfg, iamv1.ServiceName, func(server *grpc.Server) {
	iamv1.RegisterIAMServiceServer(server, iam.NewGRPCServer(security))
})
if err != nil {
	return err
}
defer grpcServer.Stop(shutdown)
```

对应入口是 `application/iam/main.go`、`application/sys/main.go` 和 `application/resource/main.go`。`StartRegisteredGRPC` 会监听 `grpc.port`，把 HTTP/gRPC endpoint 注册到配置的 Consul、etcd 或 in-process 注册中心，并在退出时注销。

## 客户端调用

业务代码依赖 `contracts.go` 暴露的 `API`，不直接保存连接，也不绕过服务发现：

```go
reg, err := registry.New(cfg.Registry.Driver, cfg.Registry.Address, cfg.Registry.Prefix)
if err != nil {
	return err
}
pool := transport.NewClientPool(reg)
defer pool.Close()

systemAPI := sysv1.NewRemote(pool)
response, err := systemAPI.CleanLogs(ctx, 30)
```

`Remote` 内部调用 `pool.Conn(ctx, service)`，从注册中心选择健康实例并缓存 gRPC connection；生成的 `NewSystemServiceClient`、`NewIAMServiceClient` 和 `NewResourceServiceClient` 只在适配层使用。

- Gateway 通过 `iamv1.NewRemote(pool)` 校验访问令牌。
- Scheduler 通过 `sysv1.NewRemote(pool)` 和 `resourcev1.NewRemote(pool)` 执行日志、资源清理任务。
- SYS、Resource 服务也通过同样的 IAM Remote 做权限校验。

## 修改流程和闭环检查

1. 先修改 `api/<service>/v1/*.proto`，为所有新增声明写注释，并确认字段编号没有冲突。
2. 执行 `make proto` 重新生成代码；不要直接修改生成文件。
3. 在 `contracts.go` 更新业务适配接口，在目标服务的 `internal/domain/rpc.go` 更新 server 实现和注册逻辑。
4. 执行以下检查：

```bash
make proto-check  # 生成结果与已提交文件一致，同时检查生成器输出
go test ./...     # 包含 Proto 注释规则测试
make verify       # test、vet、race、proto-check、swagger-check
```

5. 查看 `git diff --check`，确认只包含预期的 Proto、手写适配和生成文件变化。

### Agent 修改清单

- [ ] 找到正确的服务 Proto 和对应的 domain RPC 实现。
- [ ] service、rpc、message、字段都有紧邻 `//` 注释。
- [ ] 未复用已发布字段编号，未手改生成文件。
- [ ] 已执行 `make proto`、`make proto-check` 和 `go test ./...`。
- [ ] 客户端通过 `Remote`/`ClientPool` 使用服务发现，未硬编码服务地址。
