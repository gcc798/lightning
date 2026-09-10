-- +goose Up
WITH seed (
  id, menu_name, parent_id, sort, path, component, menu_type, perms, icon, remark
) AS (
  VALUES
    (1880159541355580001::bigint, '系统管理', 0::bigint, 10, 'system', '', 0, '', 'system', '系统内置目录'),
    (1880159541355580002::bigint, '监控审计', 0::bigint, 20, 'monitor', '', 0, '', 'monitor', '系统内置目录'),

    (1880159541355580011::bigint, '用户管理', 1880159541355580001::bigint, 10, 'user', 'system/user/index', 1, 'user.read', 'user', '系统内置菜单'),
    (1880159541355580012::bigint, '角色管理', 1880159541355580001::bigint, 20, 'role', 'system/role/index', 1, 'role.read', 'team', '系统内置菜单'),
    (1880159541355580013::bigint, '菜单管理', 1880159541355580001::bigint, 30, 'menu', 'system/menu/index', 1, 'menu.read', 'bars', '系统内置菜单'),
    (1880159541355580014::bigint, '组织管理', 1880159541355580001::bigint, 40, 'org', 'system/org/index', 1, 'org.read', 'tree', '系统内置菜单'),
    (1880159541355580015::bigint, '字典管理', 1880159541355580001::bigint, 50, 'dict', 'system/dict/index', 1, 'dict.read', 'dict', '系统内置菜单'),
    (1880159541355580016::bigint, '参数配置', 1880159541355580001::bigint, 60, 'config', 'system/config/index', 1, 'config.read', 'control', '系统内置菜单'),
    (1880159541355580017::bigint, 'API 权限', 1880159541355580001::bigint, 70, 'api-permission', 'system/apiPermission/index', 1, 'api_permission.read', 'api', '系统内置菜单'),
    (1880159541355580018::bigint, '登录日志', 1880159541355580002::bigint, 10, 'login-log', 'monitor/loginLog/index', 1, 'login_log.read', 'logininfor', '系统内置菜单'),
    (1880159541355580019::bigint, '操作日志', 1880159541355580002::bigint, 20, 'oper-log', 'monitor/operLog/index', 1, 'oper_log.read', 'bug', '系统内置菜单'),

    (1880159541355580101::bigint, '新增用户', 1880159541355580011::bigint, 10, '', '', 2, 'user.create', '', '系统内置按钮'),
    (1880159541355580102::bigint, '编辑用户', 1880159541355580011::bigint, 20, '', '', 2, 'user.update', '', '系统内置按钮'),
    (1880159541355580103::bigint, '删除用户', 1880159541355580011::bigint, 30, '', '', 2, 'user.delete', '', '系统内置按钮'),
    (1880159541355580104::bigint, 'API 权限授权', 1880159541355580011::bigint, 40, '', '', 2, 'api_permission.assign', '', '系统内置按钮'),
    (1880159541355580105::bigint, '用户角色分配', 1880159541355580011::bigint, 50, '', '', 2, 'role.assign', '', '系统内置按钮'),

    (1880159541355580201::bigint, '新增角色', 1880159541355580012::bigint, 10, '', '', 2, 'role.create', '', '系统内置按钮'),
    (1880159541355580202::bigint, '编辑角色', 1880159541355580012::bigint, 20, '', '', 2, 'role.update', '', '系统内置按钮'),
    (1880159541355580203::bigint, '删除角色', 1880159541355580012::bigint, 30, '', '', 2, 'role.delete', '', '系统内置按钮'),
    (1880159541355580204::bigint, 'API 权限授权', 1880159541355580012::bigint, 40, '', '', 2, 'api_permission.assign', '', '系统内置按钮'),
    (1880159541355580205::bigint, '角色用户分配', 1880159541355580012::bigint, 50, '', '', 2, 'role.assign', '', '系统内置按钮'),

    (1880159541355580301::bigint, '新增菜单', 1880159541355580013::bigint, 10, '', '', 2, 'menu.create', '', '系统内置按钮'),
    (1880159541355580302::bigint, '编辑菜单', 1880159541355580013::bigint, 20, '', '', 2, 'menu.update', '', '系统内置按钮'),
    (1880159541355580303::bigint, '删除菜单', 1880159541355580013::bigint, 30, '', '', 2, 'menu.delete', '', '系统内置按钮'),

    (1880159541355580401::bigint, '新增组织', 1880159541355580014::bigint, 10, '', '', 2, 'org.create', '', '系统内置按钮'),
    (1880159541355580402::bigint, '编辑组织', 1880159541355580014::bigint, 20, '', '', 2, 'org.update', '', '系统内置按钮'),
    (1880159541355580403::bigint, '删除组织', 1880159541355580014::bigint, 30, '', '', 2, 'org.delete', '', '系统内置按钮'),

    (1880159541355580501::bigint, '新增字典', 1880159541355580015::bigint, 10, '', '', 2, 'dict.create', '', '系统内置按钮'),
    (1880159541355580502::bigint, '编辑字典', 1880159541355580015::bigint, 20, '', '', 2, 'dict.update', '', '系统内置按钮'),
    (1880159541355580503::bigint, '删除字典', 1880159541355580015::bigint, 30, '', '', 2, 'dict.delete', '', '系统内置按钮'),

    (1880159541355580601::bigint, '新增配置', 1880159541355580016::bigint, 10, '', '', 2, 'config.create', '', '系统内置按钮'),
    (1880159541355580602::bigint, '编辑配置', 1880159541355580016::bigint, 20, '', '', 2, 'config.update', '', '系统内置按钮'),
    (1880159541355580603::bigint, '删除配置', 1880159541355580016::bigint, 30, '', '', 2, 'config.delete', '', '系统内置按钮'),

    (1880159541355580701::bigint, '新增 API 权限', 1880159541355580017::bigint, 10, '', '', 2, 'api_permission.create', '', '系统内置按钮'),
    (1880159541355580702::bigint, '编辑 API 权限', 1880159541355580017::bigint, 20, '', '', 2, 'api_permission.update', '', '系统内置按钮'),
    (1880159541355580703::bigint, '删除 API 权限', 1880159541355580017::bigint, 30, '', '', 2, 'api_permission.delete', '', '系统内置按钮'),

    (1880159541355580801::bigint, '清理登录日志', 1880159541355580018::bigint, 10, '', '', 2, 'login_log.delete', '', '系统内置按钮'),
    (1880159541355580901::bigint, '清理操作日志', 1880159541355580019::bigint, 10, '', '', 2, 'oper_log.delete', '', '系统内置按钮')
)
INSERT INTO s_menu (
  id, menu_name, parent_id, sort, path, component, query, is_frame, is_cache,
  menu_type, visible, status, perms, icon, remark, create_by, update_by,
  created_time, updated_time
)
SELECT
  id, menu_name, parent_id, sort, path, component, '', 0, 0,
  menu_type, 0, 0, perms, icon, remark, 0, 0, NOW(), NOW()
FROM seed
ON CONFLICT (id) DO UPDATE SET
  menu_name = EXCLUDED.menu_name,
  parent_id = EXCLUDED.parent_id,
  sort = EXCLUDED.sort,
  path = EXCLUDED.path,
  component = EXCLUDED.component,
  menu_type = EXCLUDED.menu_type,
  visible = EXCLUDED.visible,
  status = EXCLUDED.status,
  perms = EXCLUDED.perms,
  icon = EXCLUDED.icon,
  remark = EXCLUDED.remark,
  updated_time = NOW();

INSERT INTO m_menu_api_permission (
  id, menu_id, permission_id, create_by, update_by, created_time, updated_time
)
SELECT
  1880159541355590000::bigint + ROW_NUMBER() OVER (ORDER BY menu.id),
  menu.id,
  permission.id,
  0,
  0,
  NOW(),
  NOW()
FROM s_menu menu
JOIN s_api_permission permission ON permission.code = menu.perms
WHERE menu.id BETWEEN 1880159541355580011 AND 1880159541355580901
  AND menu.perms <> ''
ON CONFLICT (menu_id, permission_id) DO NOTHING;

-- 普通用户默认可进入全部业务页面，但只派生页面本身的查询权限。
INSERT INTO m_role_menu (
  id, role_id, menu_id, create_by, update_by, created_time, updated_time
)
SELECT
  1880159541355600000::bigint + ROW_NUMBER() OVER (ORDER BY menu.id),
  role.id,
  menu.id,
  0,
  0,
  NOW(),
  NOW()
FROM s_role role
CROSS JOIN s_menu menu
WHERE role.role_key = 'user'
  AND menu.id IN (
    1880159541355580001, 1880159541355580002,
    1880159541355580011, 1880159541355580012, 1880159541355580013,
    1880159541355580014, 1880159541355580015, 1880159541355580016,
    1880159541355580017, 1880159541355580018, 1880159541355580019
  );

INSERT INTO m_role_api_permission (
  id, role_id, permission_id, source, create_by, update_by, created_time, updated_time
)
SELECT
  1880159541355610000::bigint + ROW_NUMBER() OVER (ORDER BY permission.id),
  role.id,
  permission.id,
  1,
  0,
  0,
  NOW(),
  NOW()
FROM s_role role
JOIN s_api_permission permission ON permission.code IN (
  'user.read', 'role.read', 'menu.read', 'org.read', 'dict.read',
  'config.read', 'api_permission.read', 'login_log.read', 'oper_log.read'
)
WHERE role.role_key = 'user'
ON CONFLICT (role_id, permission_id, source) DO NOTHING;

-- +goose Down
DELETE FROM m_role_api_permission
WHERE role_id = (SELECT id FROM s_role WHERE role_key = 'user')
  AND source = 1
  AND permission_id IN (
    SELECT id FROM s_api_permission WHERE code IN (
      'user.read', 'role.read', 'menu.read', 'org.read', 'dict.read',
      'config.read', 'api_permission.read', 'login_log.read', 'oper_log.read'
    )
  );
DELETE FROM m_role_menu
WHERE role_id = (SELECT id FROM s_role WHERE role_key = 'user')
  AND menu_id BETWEEN 1880159541355580001 AND 1880159541355580019;
DELETE FROM m_menu_api_permission
WHERE menu_id BETWEEN 1880159541355580001 AND 1880159541355580901;
DELETE FROM s_menu
WHERE id BETWEEN 1880159541355580001 AND 1880159541355580901;
