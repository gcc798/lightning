# native 权限体系重构结果

## 1. 文档状态

本文记录 IAM 权限体系已经落地的设计。权限业务、model、迁移和 HTTP 路由全部属于 `native/application/iam/internal/`，不允许移入根 `internal` 或共享 HTTP 应用。

## 2. 最终模型

API 权限与菜单保持分离：

| 体系 | 表 | 用途 |
| --- | --- | --- |
| API 权限 | `s_api_permission`、`m_role_api_permission`、`m_user_api_permission` | 后端安全边界和真实鉴权依据 |
| 菜单 | `s_menu`、`m_role_menu` | 前端路由与按钮显隐 |
| 菜单到 API 权限 | `m_menu_api_permission` | 声明一个菜单需要派生哪些后端 API 权限 |

`s_menu.perms` 只用于前端按钮显隐，不参与后端权限派生。菜单与 API 权限通过 ID 多对多关联，不再依赖权限 code 字符串相等。

鉴权入口是 `application/iam/internal/domain/permission.go` 的 `CheckPermission`。路由在注册中显式传入 `resource` 和 `action`，不再根据权限名后缀猜测 action。

## 3. 授权来源

角色和用户 API 权限关联表包含 `source`：

- `0`：管理员手工授权。
- `1`：菜单派生授权。

唯一索引包含主体 ID、权限 ID 和 `source`。同一权限允许同时存在手工与派生两条来源记录，因此撤销任一来源不会误删另一来源。

写入规则：

- `AssignMenusToRole` 只替换该角色的派生记录。
- `AssignRolePermissions` 只替换该角色的手工记录。
- `AssignUserPermissions` 只替换该用户的手工记录。
- 鉴权查询对来源不做区分，任一有效来源均可放行。

这修复了“保存高级 API 权限后清空菜单派生权限”和“取消菜单时误删手工权限”两个数据正确性问题。

## 4. 菜单派生流程

```text
菜单关联 API permission IDs
        -> m_menu_api_permission
角色保存菜单树
        -> m_role_menu
        -> 重建 source=derived 的 m_role_api_permission
```

一个菜单可以关联多个 API 权限。角色菜单界面显示本次将派生的权限数量；高级 API 权限界面将派生项设为只读，只提交手工项。

## 5. 已完成阶段

### P1：授权来源隔离

- 为角色和用户 API 权限关联增加 `source`。
- 手工写入与菜单派生写入互不覆盖。
- 增加集成测试验证保存、撤销和重复迁移。

### P2：菜单与 API 权限多对多

- 建立 `m_menu_api_permission` 和唯一索引。
- Menu DTO 增加 `apiPermissionIds`。
- 菜单保存和查询均使用 permission ID，不再用 code 做后端关联。

### P3：前端入口收敛

- 角色菜单配置是日常授权主入口。
- API 权限入口标记为“高级”，用于纯 API 场景。
- 派生权限在高级面板只读展示。

## 6. 验收

数据库集成测试位于：

```text
native/application/iam/internal/migrations/permission_sources_test.go
```

它验证：

1. 菜单派生后保存手工权限，派生授权不丢失。
2. 清空手工权限，只删除手工来源。
3. 取消菜单，只删除派生来源，手工授权保留。

运行方式：

```bash
LIGHTNING_TEST_POSTGRES_DSN='...' go test ./application/iam/internal/migrations
```

## 7. 保留项与后续触发条件

- 保留 API 权限与菜单两套模型，不合并为单树。
- 保留 `m_user_api_permission`，用于用户直接授权。
- 不引入 Casbin，当前 SQL 查询足够直接。
- 暂不实现 client_credentials 权限表；出现不绑定用户的真实机器客户端时再增加。
- `super_admin` 目前仍由固定 `role_key` 判断；只有需要可配置超级角色时再改 Schema。
