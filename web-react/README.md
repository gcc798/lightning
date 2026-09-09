# web-react

`web-react` 是对现有 `web` 管理后台的 React + TypeScript 重写版本，同时也作为 React 学习工程使用。

详细开发约束见：[docs/react-dev-spec.md](docs/react-dev-spec.md)

如果你想把这个工程当成 React 学习案例来看，建议从这里开始：

- [docs/learn-react/README.md](docs/learn-react/README.md)

## 运行环境

- Node.js：建议 `>= 18`
- 包管理器：`pnpm`
- 后端服务：`sys-rpc`、`sys-api`
- 基础依赖：PostgreSQL、Redis

## 首次安装

在项目根目录执行：

```bash
cd web-react
pnpm install
```

## 如何启动

### 1. 启动后端依赖

`sys-api` 当前配置文件在 [sys-api.yaml](../gozero/application/sys-api/etc/sys-api.yaml)，默认依赖如下：

- `sys-api`：`http://localhost:9009`
- `sys-rpc`：`127.0.0.1:9002`
- PostgreSQL：`127.0.0.1:5433`
- Redis：`127.0.0.1:6379`

先启动 `sys-rpc`：

```bash
cd gozero/application/sys-rpc
GOCACHE=/tmp/lightning-go-build go run sys.go -f etc/sys-rpc.yaml
```

再启动 `sys-api`：

```bash
cd gozero/application/sys-api
GOCACHE=/tmp/lightning-go-build go run sys.go -f etc/sys-api.yaml
```

### 2. 启动前端开发服务

```bash
cd web-react
pnpm dev
```

默认访问地址：

- 前端开发服务器：`http://localhost:3001`
- 后端接口地址：`http://localhost:9009`

### 3. 生产构建

```bash
cd web-react
pnpm build
```

构建产物输出到 `web-react/dist`。

### 4. 本地预览构建结果

```bash
cd web-react
pnpm preview
```

## 如何修改配置文件

### 1. 修改前端环境变量

开发环境配置文件：

- [web-react/.env.development](.env.development)

生产环境配置文件：

- [web-react/.env.production](.env.production)

当前支持的主要配置项：

- `VITE_APP_TITLE`
  - 前端应用标题
- `VITE_API_BASE_URL`
  - 后端接口基础地址
  - 当前前端请求层直接使用这个地址发请求，不走本地代理
- `VITE_CLIENT_KEY`
  - 登录相关客户端标识
- `VITE_CLIENT_SECRET`
  - 登录相关客户端密钥

示例：

```env
VITE_APP_TITLE=lightning Admin React
VITE_API_BASE_URL=http://localhost:9009
VITE_CLIENT_KEY=web-admin
VITE_CLIENT_SECRET=web-secret-2024
```

修改环境变量后，需要重新执行 `pnpm dev` 或重新构建，Vite 才会生效。

### 2. 修改前端开发端口

开发端口配置在 [web-react/vite.config.ts](vite.config.ts)：

```ts
server: {
  port: 3001,
}
```

如果你想改成其他端口，例如 `3002`，直接修改这里即可。修改后需要重启开发服务器。

### 3. 修改后端服务配置

后端主配置文件：

- [sys-api.yaml](../gozero/application/sys-api/etc/sys-api.yaml)

这个文件里可以调整：

- `Port`
- `SysRpc.Target`
- `Postgres.Dsn`
- `Redis.Addr`
- `Redis.Password`
- `Jwt.Secret`
- `Captcha` 开关

如果你改了后端端口或域名，前端的 `VITE_API_BASE_URL` 也要一起改。

## 常见启动问题

### 1. 页面能打开，但登录或列表接口报错

优先检查：

- `sys-rpc` 是否已经启动
- `sys-api` 是否已经启动
- PostgreSQL 和 Redis 是否与 [sys-api.yaml](../gozero/application/sys-api/etc/sys-api.yaml) 一致
- `VITE_API_BASE_URL` 是否指向正确地址

### 2. 修改了 `.env` 但页面没变化

原因通常是 Vite 已经启动，环境变量不会热更新。处理方式：

```bash
Ctrl + C
pnpm dev
```

### 3. 修改了后端配置但前端仍然请求旧地址

检查两个地方：

- 是否真的修改了 [web-react/.env.development](.env.development) 或 [web-react/.env.production](.env.production)
- 前端开发服务是否已经重启

## 目录说明

- `src/app`：应用入口
- `src/router`：静态路由和动态路由装配
- `src/store`：登录态、权限、主题、应用状态
- `src/api`：按后端模块拆分的接口定义
- `src/pages`：页面实现
- `src/components`：可复用组件
- `src/types`：公共类型定义
- `src/utils`：请求层、权限工具、菜单工具
- `docs`：开发规范和补充文档

## 开发说明

- 接口请求以 `gozero/application/sys-api/api/*.api` 为主契约来源
- 返回值结构优先参考后端 `internal/logic` 实际组装
- 仅在无法确认真实返回值时，才参考现有 `web` 工程
- 本工程要求良好的中文注释，重点解释复杂逻辑和设计原因
- 当前不要求编写前端单元测试
