# HTTP 路由边界说明

## 目标

`sys-api` 的 HTTP 路由遵循一条简单规则：

- 标准业务接口，优先使用 **proto + google.api.http** 定义。
- 只有在 proto 生成的 HTTP Handler 不适合表达时，才保留**手工注册路由**。

这样做的目的有两个：

1. 大多数业务接口继续走 Kratos 默认的 proto-first 开发方式。
2. 少数必须依赖具体传输协议的接口，也能被明确隔离出来，避免代码边界混乱。

---

## Proto 路由

这类路由通过生成代码里的 `Register*HTTPServer(...)` 完成注册。

当前入口文件：
- [../application/sys-api/internal/server/http_routes.go](../application/sys-api/internal/server/http_routes.go)

对应函数：
- `registerProtoHTTPRoutes(...)`

当前走 proto 注册的模块有：
- `health`
- `auth`
- `captcha`
- `menu`
- `user`
- `role`
- `org`
- `config`
- `dict`
- `loginlog`
- `operlog`
- `storage-env`

这些接口适合放进 proto 的原因是：
- 请求和响应基本都是标准 JSON
- 既需要 HTTP 契约，也需要 gRPC 契约
- Kratos 生成代码可以直接覆盖这类场景

---

## 手工路由

这类路由保留手工注册，不强行塞进 proto HTTP 生成链路。

当前入口文件：
- [../application/sys-api/internal/server/http_routes.go](../application/sys-api/internal/server/http_routes.go)

对应函数：
- `registerManualHTTPRoutes(...)`

当前保留手工注册的场景有：
- 附件上传：`multipart/form-data`
- 附件下载：二进制文件流、下载响应头
- metrics 指标接口
- 健康探针接口
- swagger 文档接口
- websocket 接口

这些接口共同特点是：
- 强依赖 HTTP 传输细节
- 不适合用标准 JSON body 表达
- 或者根本就不属于业务领域契约的一部分

---

## 路由归类规则

后续新增接口时，按下面规则判断：

1. 如果是标准 JSON 业务接口，并且也需要 gRPC 契约，放进 proto。
2. 如果依赖具体 HTTP 传输行为，就保留手工注册。

典型的“适合手工注册”的场景：
- multipart 上传
- 文件流下载
- websocket 升级
- metrics 原始暴露
- swagger 静态页面/文档入口
- 不属于业务 service 的健康探针接口

---

## 为什么不强行统一成一种方式

如果把所有路由都强行塞进 proto：
- multipart、stream、websocket 这类接口会变得别扭
- 代码可读性会下降
- 后续维护成本更高

如果把所有路由都改成手工注册：
- 会丢掉 Kratos proto-first 的核心优势
- HTTP/gRPC 契约会重新分裂
- 生成代码和约束能力也会被削弱

所以当前项目的标准是：

- 业务接口：**proto-first**
- 传输特化接口：**手工注册，但要显式隔离**

---

## 当前相关代码位置

主 HTTP Server 装配：
- [../application/sys-api/internal/server/http.go](../application/sys-api/internal/server/http.go)

手工 Handler 相关文件：
- [../application/sys-api/internal/server/attachment_handler.go](../application/sys-api/internal/server/attachment_handler.go)
- [../application/sys-api/internal/server/metrics_handler.go](../application/sys-api/internal/server/metrics_handler.go)
- [../application/sys-api/internal/server/health_handler.go](../application/sys-api/internal/server/health_handler.go)
- [../application/sys-api/internal/server/swagger_handler.go](../application/sys-api/internal/server/swagger_handler.go)
- [../application/sys-api/internal/server/websocket_handler.go](../application/sys-api/internal/server/websocket_handler.go)
