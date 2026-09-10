import { Suspense, useEffect, useMemo, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import {
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
} from '@/utils/icons';
import { Avatar, Button, Content, Dropdown, Header, Layout, Menu, Sider, Space, type MenuItem } from '@/components/ui';
import type { MenuRecord } from '@/types/menu';
import { PageLoading } from '@/components/common/PageLoading';
import { ThemeSwitcher } from '@/components/common/ThemeSwitcher';
import { usePermissionStore } from '@/store/permission';
import { useAppStore } from '@/store/app';
import { useAuthStore } from '@/store/auth';
import { getMenuIconNode } from '@/utils/icons';
import { findFirstNavigablePath, isMenuHidden, joinMenuPath } from '@/utils/menu';
import { isNumericValue } from '@/utils/number';

function findOpenKeys(pathname: string) {
  const parts = pathname.split('/').filter(Boolean);
  if (parts.length <= 1) {
    return [];
  }

  return parts.slice(0, -1).map((_, index) => joinMenuPath('', parts.slice(0, index + 1).join('/')));
}

function findMenuTrail(
  menuTree: MenuRecord[],
  pathname: string,
  parentPath = '',
  trail: MenuRecord[] = [],
): MenuRecord[] {
  for (const menu of menuTree) {
    if (isMenuHidden(menu)) {
      continue;
    }

    const fullPath = joinMenuPath(parentPath, menu.path);
    const nextTrail = [...trail, menu];

    if (fullPath === pathname && !isNumericValue(menu.menuType, 2)) {
      return nextTrail;
    }

    const childTrail = findMenuTrail(menu.children ?? [], pathname, fullPath, nextTrail);
    if (childTrail.length) {
      return childTrail;
    }
  }

  return [];
}

function buildMenuItems(
  menuTree: MenuRecord[],
  navigate: (path: string) => void,
  parentPath = '',
): MenuItem[] {
  return menuTree
    .filter((menu) => !isMenuHidden(menu) && !isNumericValue(menu.menuType, 2))
    .map((menu) => {
      const fullPath = joinMenuPath(parentPath, menu.path);
      const icon = getMenuIconNode(menu.icon);
      const targetPath =
        isNumericValue(menu.menuType, 1)
          ? fullPath
          : findFirstNavigablePath(menu.children ?? [], fullPath) ?? fullPath;
      const childItems = buildMenuItems(menu.children ?? [], navigate, fullPath);

      // 这里不能直接看原始 children.length。
      // 因为很多“页面菜单”下面挂的是按钮权限节点（menuType=2），它们不应该把页面菜单渲染成可展开子菜单。
      if (childItems.length > 0) {
        return {
          key: fullPath,
          icon,
          label: menu.menuName,
          // 目录标题点击时直接进入第一个真实页面，避免只展开不跳转造成“没反应”的错觉。
          onTitleClick: () => {
            if (targetPath && targetPath !== fullPath) {
              navigate(targetPath);
            }
          },
          children: childItems,
        } satisfies MenuItem;
      }

      return {
        key: fullPath,
        icon,
        label: menu.menuName,
        onClick: () => navigate(targetPath),
      } satisfies MenuItem;
    });
}

export function AppLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const menuTree = usePermissionStore((state) => state.menuTree);
  const collapsed = useAppStore((state) => state.collapsed);
  const setCollapsed = useAppStore((state) => state.setCollapsed);
  const toggleCollapsed = useAppStore((state) => state.toggleCollapsed);
  const authStore = useAuthStore();

  const menuItems = useMemo(
    () => buildMenuItems(menuTree, (path) => navigate(path)),
    [menuTree, navigate],
  );
  const currentOpenKeys = useMemo(() => findOpenKeys(location.pathname), [location.pathname]);
  const currentTrail = useMemo(
    () => findMenuTrail(menuTree, location.pathname),
    [location.pathname, menuTree],
  );
  const [openKeys, setOpenKeys] = useState<string[]>(currentOpenKeys);

  useEffect(() => {
    setOpenKeys(currentOpenKeys);
  }, [currentOpenKeys]);

  return (
    <div className="app-shell-wrap">
      <Layout className="app-shell">
        <Sider
          breakpoint="lg"
          className="app-sider"
          collapsed={collapsed}
          collapsible
          onCollapse={setCollapsed}
          theme="light"
          trigger={null}
          width={224}
        >
          <button
            className="app-logo"
            type="button"
            onClick={() => navigate('/dashboard')}
          >
            {!collapsed ? (
              <span className="app-logo-text">
                <strong>Admin</strong>
              </span>
            ) : null}
          </button>

          <Menu
            className="app-menu"
            inlineIndent={14}
            mode="inline"
            openKeys={openKeys}
            selectedKeys={[location.pathname]}
            items={menuItems}
            onOpenChange={(keys) => setOpenKeys(keys as string[])}
            theme="light"
          />
        </Sider>

        <Layout>
          <Header className="app-header">
            <Space align="center" className="app-header-group" size={12}>
              <Button
                className="header-action-btn"
                icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                type="text"
                onClick={toggleCollapsed}
              />
              <div className="header-page-context">
                <nav aria-label="当前页面" className="header-breadcrumb">
                  {(currentTrail.length ? currentTrail.map((menu) => menu.menuName) : ['工作台'])
                    .map((title, index, items) => (
                      <span
                        className={`header-breadcrumb-item${index === items.length - 1 ? ' is-current' : ''}`}
                        key={`${title}-${index}`}
                      >
                        <span>{title}</span>
                        {index < items.length - 1 ? (
                          <span className="header-breadcrumb-separator">/</span>
                        ) : null}
                      </span>
                    ))}
                </nav>
              </div>
            </Space>

            <Space align="center" className="app-header-group" size={12}>
              <ThemeSwitcher />
              <Dropdown
                menu={{
                  items: [
                    {
                      key: 'logout',
                      icon: <LogoutOutlined />,
                      label: '退出登录',
                      onClick: async () => {
                        await authStore.logout();
                        navigate('/login', { replace: true });
                      },
                    },
                  ],
                }}
              >
                <Button
                  aria-label="用户菜单"
                  className="header-user-btn"
                  icon={<Avatar icon={<UserOutlined />} shape="square" size={28} />}
                  type="text"
                />
              </Dropdown>
            </Space>
          </Header>

          <Content className="app-content">
            <Suspense fallback={<PageLoading />}>
              <div className="app-view" key={location.pathname}>
                <Outlet />
              </div>
            </Suspense>
          </Content>
        </Layout>
      </Layout>
    </div>
  );
}
