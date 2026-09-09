# Ent 使用说明

这份文档分两部分：

1. 先讲 `ent` 从 0 到 1 的基本使用方式。
2. 再讲当前 `kratos` 工程里，`ent` 是如何落地和被调用的。

如果你想先建立整体理解，建议先看这份文档，再去看：

- [../application/sys-rpc/ent/schema](../application/sys-rpc/ent/schema)
- [../application/sys-rpc/internal/data/resources.go](../application/sys-rpc/internal/data/resources.go)
- [../application/sys-rpc/internal/data](../application/sys-rpc/internal/data)

---

## 一、什么是 Ent

`ent` 是一个 Go 生态里的强类型 ORM / 实体代码生成工具。

你可以把它理解成：

- 先写 schema
- 再生成 Go 代码
- 然后通过生成出来的 `Client`、`Query`、`Create`、`Update`、`Tx` 来访问数据库

它和“运行时靠反射拼模型”的 ORM 不太一样，`ent` 更强调：

- schema 明确
- 代码生成
- 编译期类型安全

---

## 二、Ent 从 0 到 1 的最小流程

### 1. 定义 schema

`ent` 的起点不是直接写 repo，而是先定义 schema。

例如在一个最小工程里，你会有这样的目录：

```text
ent/
└── schema/
    └── user.go
```

一个典型 schema 会包含：

- 字段
- 索引
- 边（关系）
- 表名注解
- mixin

比如当前工程中的用户 schema：

- [../application/sys-rpc/ent/schema/user.go](../application/sys-rpc/ent/schema/user.go)

这里做了几件事：

1. 定义表字段
2. 通过 `Annotations()` 指定真实表名 `s_user`
3. 通过 `Indexes()` 定义索引
4. 通过 `Mixin()` 复用公共审计字段

---

### 2. 执行代码生成

定义好 schema 之后，要执行 `ent generate`。

当前工程对应命令在：

- [../Makefile](../Makefile)

对应 target：

```make
ent:
	ent generate ./application/sys-rpc/ent/schema
```

执行：

```bash
make ent
```

生成后会得到：

- `application/sys-rpc/ent/client.go`
- `application/sys-rpc/ent/<entity>.go`
- `application/sys-rpc/ent/<entity>_create.go`
- `application/sys-rpc/ent/<entity>_query.go`
- `application/sys-rpc/ent/<entity>_update.go`
- `application/sys-rpc/ent/migrate/*`

这些生成代码是 `ent` 的运行基础。

---

### 3. 初始化 Client

生成代码之后，要在程序启动时创建 `*ent.Client`。

最常见的方式有两种：

1. `ent.Open(driver, dsn)`
2. 自己先准备 `*sql.DB`，再用 `ent.NewClient(ent.Driver(...))`

当前工程两种都支持，但 PostgreSQL 走了更完整的一层封装。

代码在：

- [../application/sys-rpc/internal/data/resources.go](../application/sys-rpc/internal/data/resources.go)

里面的关键逻辑是：

1. 先读 `conf.Data`
2. 初始化底层 `*sql.DB`
3. 包上观测能力
4. 创建 `*ent.Client`

当前 PostgreSQL 分支是：

```go
sqlDB, err = openObservedPostgresDB(...)
client = ent.NewClient(ent.Driver(entsql.OpenDB(..., sqlDB)))
```

这说明当前工程不是裸用 `ent.Open`，而是把：

- `pgx/v5/stdlib`
- 慢 SQL
- DB metrics
- 连接池采样

先挂到底层 `sql.DB` 上，再交给 `ent`。

当前工程里 PostgreSQL 的实际底层驱动已经是：

- `github.com/jackc/pgx/v5/stdlib`

也就是说，当前链路是：

```text
pgx/v5/stdlib
  -> database/sql.(*sql.DB)
  -> ent.Client
```

这意味着：

1. `ent` 底层复用了 `SQLDB`
2. 慢 SQL 日志和连接池指标是挂在 `SQLDB` 这一层的
3. repo 层日常查询仍然通过 `Ent` 完成，不直接操作 `SQLDB`

---

### 4. 在业务代码里使用 Client

有了 `*ent.Client` 后，最常见的操作就是：

#### 查询

```go
item, err := client.User.Get(ctx, id)
items, err := client.User.Query().All(ctx)
```

#### 创建

```go
item, err := client.User.Create().
	SetID(id).
	SetUserName(name).
	Save(ctx)
```

#### 更新

```go
_, err := client.User.UpdateOneID(id).
	SetNickName(nick).
	Save(ctx)
```

#### 删除

```go
_, err := client.User.Delete().
	Where(user.ID(id)).
	Exec(ctx)
```

当前工程里大量代码都在这样用，只是又包了一层 repo 语义。

例如：

- [../application/sys-rpc/internal/data/user_repo.go](../application/sys-rpc/internal/data/user_repo.go)
- [../application/sys-rpc/internal/data/role_repo.go](../application/sys-rpc/internal/data/role_repo.go)

---

### 5. 事务怎么用

`ent` 的事务是基于一个 `*ent.Client` 开出来的。

标准方式是：

```go
tx, err := client.Tx(ctx)
```

然后在事务里用：

```go
tx.User.Create()...
tx.Role.UpdateOneID(...)...
```

最后：

- `tx.Commit()`
- 或 `tx.Rollback()`

当前工程已经把这个模式封装到：

- [../application/sys-rpc/internal/data/resources.go](../application/sys-rpc/internal/data/resources.go)

函数：

- `func (r *Resources) withTx(ctx context.Context, fn func(tx *ent.Tx) error) error`

这层封装的好处是：

1. repo 不用反复手写 `Tx / Rollback / Commit`
2. panic 时也会回滚
3. rollback 失败会补充错误信息

这是当前工程推荐的事务使用方式。

---

## 三、当前工程中 Ent 的目录结构

当前 `sys-rpc` 使用 `ent`，目录分成三层理解最清楚。

### 1. schema 定义层

目录：

- [../application/sys-rpc/ent/schema](../application/sys-rpc/ent/schema)

这里是开发者手写的地方。

当前主要 schema 有：

- `user.go`
- `role.go`
- `org.go`
- `menu.go`
- `config.go`
- `dict_data.go`
- `login_log.go`
- `oper_log.go`
- `storage_env.go`
- `attachment.go`
- `auth_client.go`
- `casbin_rule.go`

公共 mixin：

- [../application/sys-rpc/ent/schema/audit_mixin.go](../application/sys-rpc/ent/schema/audit_mixin.go)

它统一提供：

- `create_by`
- `update_by`
- `created_time`
- `updated_time`
- `deleted_at`

这也是当前工程里“软删除”的基础之一。

---

### 2. generated code 层

目录：

- [../application/sys-rpc/ent](../application/sys-rpc/ent)

这里基本都是 `ent generate` 自动生成的代码。

例如：

- `client.go`
- `user.go`
- `user_create.go`
- `user_query.go`
- `user_update.go`
- `migrate/*`

这层代码：

- 允许阅读
- 允许用于理解 API
- **不要手改**

如果 schema 有变更，应该改 `ent/schema/*.go`，然后重新执行：

```bash
make ent
```

---

### 3. 项目自己的 data/repo 层

目录：

- [../application/sys-rpc/internal/data](../application/sys-rpc/internal/data)

这是当前工程把 `ent` 真正接入业务的地方。

这里做的事情包括：

1. 初始化 `ent.Client`
2. 封装事务
3. 做 DTO 转换
4. 组合多表查询
5. 做业务约束校验
6. 和 Redis / Storage / JWT 等其他资源协同

也就是说：

- `ent` 负责数据库访问能力
- `internal/data` 负责把数据库能力变成项目实际可用的 repo 行为

---

## 四、当前工程里 Ent 是怎么启动起来的

### 1. 配置来源

配置来自：

- [../application/sys-rpc/internal/conf/conf.proto](../application/sys-rpc/internal/conf/conf.proto)
- [../application/sys-rpc/configs/config.yaml](../application/sys-rpc/configs/config.yaml)

对应 `Data.Database` 配置里至少要有：

- `driver`
- `dsn`

---

### 2. Wire 注入入口

`ent` 相关资源是从 `internal/data` 注入出来的。

入口：

- [../application/sys-rpc/internal/data/data.go](../application/sys-rpc/internal/data/data.go)

这里的 `ProviderSet` 最终会把：

- `*conf.Data`
- `*conf.Auth`
- `*conf.JWT`
- `*conf.Observability`

交给：

- `NewData(...)`

而 `NewData(...)` 内部会继续调用：

- `NewResources(...)`

最后把 `*Resources` 提供给后面的 biz / service 层。

---

### 3. Resources 是当前工程的 Ent 根资源

核心文件：

- [../application/sys-rpc/internal/data/resources.go](../application/sys-rpc/internal/data/resources.go)

这里的 `Resources` 里包含：

- `Ent *ent.Client`
- `SQLDB *sql.DB`
- `Redis *redis.Client`
- `Storage *StorageManager`
- `WeChat *WeChatManager`
- `JWT *JWTManager`

也就是说，当前工程不是把 `ent.Client` 裸着到处传，而是把它放进统一资源容器里。

这两个数据库相关字段的分工可以这样理解：

- `SQLDB`
  - 是原生 `database/sql` 连接池
  - 当前 PostgreSQL 底层驱动是 `pgx/v5/stdlib`
  - 主要用于承载底层连接、慢 SQL、DB metrics、连接池采样
- `Ent`
  - 是基于 `SQLDB` 构建出来的 `ent.Client`
  - 主要用于 repo 层的查询、创建、更新、删除和事务

关系就是：

```text
SQLDB (*sql.DB)
   ↓
Ent (*ent.Client)
   ↓
repo / biz 使用
```

这也是当前工程的一条设计约定：

- repo 方法大多挂在 `*Resources` 上
- 而不是一个表一个 repo struct

这和很多教材里的写法不完全一样，但在当前工程里是成立的。

---

## 五、当前工程里如何基于 Ent 写业务查询

### 1. 先把 ent entity 转成 proto DTO

当前工程不会把 `ent` 生成的 entity 直接往 RPC 外返回，而是先做转换。

例如用户：

- [../application/sys-rpc/internal/data/user_repo.go](../application/sys-rpc/internal/data/user_repo.go)

函数：

- `userEntityToItem(item *entpkg.User) *v1.UserItem`

这个函数的作用是把：

- `*ent.User`

转换成：

- `*api.system.v1.UserItem`

这一步很重要，因为：

1. 数据库模型不应该直接暴露给 RPC
2. proto 字段名和数据库字段名不一定一样
3. 时间、状态、空值有时需要额外处理

---

### 2. 典型查询

例如：

```go
item, err := r.Ent.User.Get(ctx, id)
```

这表示按主键查询一条用户记录。

当前工程里又补了一层软删除判断：

```go
if item.DeletedAt != nil {
	return nil, nil
}
```

这说明当前工程不是数据库物理删除优先，而是：

- 逻辑删除优先
- 查询时主动过滤 `deleted_at`

---

### 3. 典型创建

例如：

```go
item, err := r.Ent.Role.Create().
	SetID(nextID()).
	SetRoleKey(req.GetRoleKey()).
	SetRoleName(req.GetRoleName()).
	...
	Save(ctx)
```

这就是典型的 Ent Create Builder 风格。

你可以把它理解成：

1. 选中实体：`Role`
2. 进入 `Create()`
3. 设置字段
4. `Save(ctx)`

---

### 4. 典型更新

例如：

```go
_, err := r.Ent.User.UpdateOneID(req.GetUserId()).
	SetNickName(req.GetNickName()).
	SetStatus(req.GetStatus()).
	Save(ctx)
```

如果你已经知道主键，当前工程通常优先用：

- `UpdateOneID(...)`

而不是先查一遍再更新。

---

### 5. 典型事务

例如当前工程删除角色时，会同时处理：

- 角色自身
- 角色菜单关系

代码在：

- [../application/sys-rpc/internal/data/role_repo.go](../application/sys-rpc/internal/data/role_repo.go)

这里走的是：

```go
return r.withTx(ctx, func(tx *entpkg.Tx) error {
    ...
})
```

事务内统一使用：

- `tx.Role`
- `tx.RoleMenu`
- `tx.UserRole`

而不是 `r.Ent.Role`

这是当前工程里最重要的一条事务约定：

- **一旦进入事务，事务内所有数据库操作都必须走 `tx`，不要再走外部 `r.Ent`。**

---

## 六、如果你要新增一张表，当前工程怎么做

这里给你一个最贴近当前工程的顺序。

### 第一步：写 schema

在：

- [../application/sys-rpc/ent/schema](../application/sys-rpc/ent/schema)

新增一个文件，比如：

- `tenant.go`

定义：

1. 结构体
2. `Fields()`
3. `Indexes()`
4. `Annotations()`
5. 如果需要，复用 `AuditMixin`

---

### 第二步：执行生成

执行：

```bash
make ent
```

然后就会生成：

- `application/sys-rpc/ent/tenant*.go`

---

### 第三步：补 data 层

在：

- [../application/sys-rpc/internal/data](../application/sys-rpc/internal/data)

新增对应 repo 文件，例如：

- `tenant_repo.go`

通常要做：

1. `entity -> proto DTO` 转换函数
2. `GetTenant`
3. `CreateTenant`
4. `UpdateTenant`
5. `DeleteTenant`
6. `PageTenants`

---

### 第四步：补 biz 层

在：

- [../application/sys-rpc/internal/biz](../application/sys-rpc/internal/biz)

增加 usecase 方法，把 data 层暴露出去。

---

### 第五步：补 service 层

在：

- [../application/sys-rpc/internal/service](../application/sys-rpc/internal/service)

把 proto server 接口实现补齐。

---

### 第六步：如果改了 proto，再补 sys-api

如果你这张新表对应新的 HTTP / gRPC 接口，那么还要：

1. 改 `api/system/v1/*.proto`
2. 执行：

```bash
make proto-all
make wire
```

3. 再补 `sys-api` 的：
   - `internal/data`
   - `internal/biz`
   - `internal/service`

---

## 七、当前工程使用 Ent 的几个重要约定

### 1. `ent/schema` 是手写区，`ent/*` 是生成区

你应该改：

- `application/sys-rpc/ent/schema/*.go`

不要手改：

- `application/sys-rpc/ent/*.go`

---

### 2. 当前工程用 `deleted_at` 做软删除

这个规则来自：

- [../application/sys-rpc/ent/schema/audit_mixin.go](../application/sys-rpc/ent/schema/audit_mixin.go)

所以很多 repo 里会出现：

1. 删除时：
   - 更新 `deleted_at`
2. 查询时：
   - 过滤掉 `deleted_at != nil`

这点你后面新增表时要保持一致。

---

### 3. 当前工程的分页很多还是“查全量后内存过滤”

例如：

- [`activeUsers(...)`](../application/sys-rpc/internal/data/user_repo.go)
- [`activeRoles(...)`](../application/sys-rpc/internal/data/role_repo.go)

这是当前工程的实际状态，不是 Ent 的限制。

如果后面数据量变大，更合理的优化方向是：

- 把过滤条件尽量下推到 `ent.Query()`
- 把分页改成数据库侧分页

也就是说：

- 当前工程已经能工作
- 但部分查询实现还有继续优化空间

---

### 4. 当前工程只把 Ent 用在 `sys-rpc`

原因很明确：

- `sys-api` 不直接查库
- `sys-rpc` 承担数据库访问

所以 `ent` 只存在于：

- `application/sys-rpc`

这和当前整体分层是匹配的。

---

## 八、当前工程里 Ent 还没有替你做的事

有几件事要特别注意。

### 1. `make ent` 只负责生成代码，不负责自动改数据库

当前工程里：

```bash
make ent
```

本质上只是：

```bash
ent generate ./application/sys-rpc/ent/schema
```

它会生成 Go 代码，但**不会自动把数据库表结构改掉**。

如果 schema 变了，数据库迁移还需要你额外处理。

当前工程里虽然已经生成了：

- `application/sys-rpc/ent/migrate`

但并没有把“自动执行迁移”封装成一个标准命令链路。

所以你要区分两件事：

1. 生成 Ent Go 代码
2. 真正更新数据库结构

这不是一回事。

---

### 2. Ent 不会自动替你完成业务校验

例如这些逻辑：

- 用户名唯一性检查
- 角色是否系统角色
- 角色是否仍被用户引用
- 删除前是否允许操作

这些都是当前工程在 repo 层额外写的业务逻辑，不是 `ent` 自动赠送的。

例如：

- [../application/sys-rpc/internal/data/user_repo.go](../application/sys-rpc/internal/data/user_repo.go)
- [../application/sys-rpc/internal/data/role_repo.go](../application/sys-rpc/internal/data/role_repo.go)

---

### 3. Ent 事务只覆盖一个 `ent.Client`

这一点很重要。

当前工程只有一个主数据库 client，所以没问题。

但如果以后你扩到多个独立数据库：

- 2 个 PostgreSQL
- 1 个 MySQL

那 `ent` 本身并不会把多个 client 自动变成一个本地事务。

它能很好处理的是：

- 一个 `ent.Client`
- 一个 `ent.Tx`

跨多个独立数据源，就不是 `ent` 单独能解决的。

---

## 九、推荐你阅读源码时的顺序

如果你想靠当前工程理解 Ent，最推荐的顺序是：

1. 看 schema
   - [../application/sys-rpc/ent/schema/user.go](../application/sys-rpc/ent/schema/user.go)
   - [../application/sys-rpc/ent/schema/audit_mixin.go](../application/sys-rpc/ent/schema/audit_mixin.go)
2. 看资源初始化
   - [../application/sys-rpc/internal/data/resources.go](../application/sys-rpc/internal/data/resources.go)
3. 看一个简单 repo
   - [../application/sys-rpc/internal/data/user_repo.go](../application/sys-rpc/internal/data/user_repo.go)
4. 看一个带事务的 repo
   - [../application/sys-rpc/internal/data/role_repo.go](../application/sys-rpc/internal/data/role_repo.go)
5. 再看 biz 怎么调用 data
   - [../application/sys-rpc/internal/biz/user.go](../application/sys-rpc/internal/biz/user.go)

按这个顺序看，会比直接扎进 `ent` 生成代码里更容易理解。

---

## 十、一句话总结

在当前工程里，你可以把 `ent` 理解成：

- `schema` 负责定义表结构和模型约束
- `make ent` 负责生成强类型数据库访问代码
- `Resources` 负责初始化 `ent.Client`
- `internal/data` 负责把 `ent` 访问能力包装成项目 repo
- `withTx(...)` 负责统一事务
- `biz` / `service` 只消费这些 repo 能力，不直接操作数据库

也就是说，`ent` 在我们工程里不是单独裸用，而是已经被纳入了 Kratos 的：

- `data`
- `biz`
- `service`

这套分层里。
