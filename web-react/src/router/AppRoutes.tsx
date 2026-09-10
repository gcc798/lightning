import { lazy, Suspense, useEffect, useMemo } from 'react';
import type { LazyExoticComponent, ReactElement } from 'react';
import { Navigate, useLocation, useRoutes } from 'react-router-dom';
import type { RouteObject } from 'react-router-dom';
import { message } from '@/components/ui';
import type { MenuRouteRecord } from '@/types/menu';
import type { MenuRecord } from '@/types/menu';
import { PageLoading } from '@/components/common/PageLoading';
import { AppLayout } from '@/components/layout/AppLayout';
import { useAuthStore } from '@/store/auth';
import { usePermissionStore } from '@/store/permission';
import { findFirstNavigablePath, flattenLeafMenus, isMenuHidden, joinMenuPath } from '@/utils/menu';
import { isNumericValue } from '@/utils/number';

const LoginPage = lazy(() => import('@/pages/auth/login'));
const DashboardPage = lazy(() => import('@/pages/dashboard'));
const NotFoundPage = lazy(() => import('@/pages/error/404'));
const ForbiddenPage = lazy(() => import('@/pages/error/403'));
const ServerErrorPage = lazy(() => import('@/pages/error/500'));

const pageModules = import.meta.glob('../pages/**/*.tsx');

function normalizeComponentPath(component?: string) {
  if (!component) {
    return '';
  }

  return component
    .replace(/^\/?src\/views\//, '')
    .replace(/^\//, '');
}

function resolvePageComponent(menu: MenuRouteRecord): LazyExoticComponent<() => ReactElement> {
  const normalized = normalizeComponentPath(menu.component);
  const candidate = `../pages/${normalized}.tsx`;
  const moduleLoader = pageModules[candidate] ?? pageModules['../pages/error/404.tsx'];

  // 动态菜单只负责告诉前端“应该打开哪个页面”，真正的 React 页面文件
  // 仍然由本地 pages 目录承载。找不到时退回 404，避免白屏。
  return lazy(async () => {
    const module = (await moduleLoader()) as { default: () => ReactElement };
    return module;
  });
}

function RequireAuth() {
  const location = useLocation();
  const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
  const loadMenuTree = usePermissionStore((state) => state.loadMenuTree);
  const isMenuLoaded = usePermissionStore((state) => state.isMenuLoaded);
  const isMenuLoading = usePermissionStore((state) => state.isMenuLoading);
  const menuLoadError = usePermissionStore((state) => state.menuLoadError);

  useEffect(() => {
    if (isLoggedIn && !isMenuLoaded && !isMenuLoading) {
      // 首次进入后台时先拉菜单树。
      // React 版的动态路由和侧边栏都基于这份菜单数据生成。
      void loadMenuTree();
    }
  }, [isLoggedIn, isMenuLoaded, isMenuLoading, loadMenuTree]);

  useEffect(() => {
    if (menuLoadError) {
      message.error(menuLoadError);
    }
  }, [menuLoadError]);

  if (!isLoggedIn) {
    const redirect = encodeURIComponent(`${location.pathname}${location.search}`);
    return <Navigate replace to={`/login?redirect=${redirect}`} />;
  }

  if (!isMenuLoaded || isMenuLoading) {
    // 这里必须在 Outlet 渲染前拦住，否则用户直达动态路由时会先命中 404。
    return <PageLoading tip="正在加载菜单和权限..." />;
  }

  return <AppLayout />;
}

function buildDynamicRoutes(menuTree: MenuRouteRecord[]): RouteObject[] {
  return menuTree.map((menu) => {
    const PageComponent = resolvePageComponent(menu);
    return {
      // useRoutes 下挂在根布局的子路由要使用相对 path，因此这里去掉前导 /。
      path: menu.fullPath.replace(/^\//, ''),
      element: <PageComponent />,
    };
  });
}

function buildDirectoryRedirectRoutes(
  menuTree: MenuRecord[],
  parentPath = '',
): RouteObject[] {
  const routes: RouteObject[] = [];

  menuTree.forEach((menu) => {
    if (isMenuHidden(menu)) {
      return;
    }

    const fullPath = joinMenuPath(parentPath, menu.path);

    if (isNumericValue(menu.menuType, 0) && menu.children?.length) {
      const targetPath = findFirstNavigablePath(menu.children, fullPath);

      if (targetPath) {
        routes.push({
          path: fullPath.replace(/^\//, ''),
          element: <Navigate replace to={targetPath} />,
        });
      }
    }

    if (menu.children?.length) {
      routes.push(...buildDirectoryRedirectRoutes(menu.children, fullPath));
    }
  });

  return routes;
}

function TitleSync() {
  const location = useLocation();
  const menuTree = usePermissionStore((state) => state.menuTree);

  useEffect(() => {
    const titleMap = new Map<string, string>([
      ['/login', '登录'],
      ['/dashboard', '仪表盘'],
      ['/403', '403'],
      ['/404', '404'],
      ['/500', '500'],
    ]);

    // 页面标题优先从后端菜单树里取，保持和旧后台菜单配置一致。
    flattenLeafMenus(menuTree).forEach((menu) => {
      titleMap.set(menu.fullPath, menu.menuName);
    });

    const currentTitle = titleMap.get(location.pathname) ?? import.meta.env.VITE_APP_TITLE;
    document.title = `${currentTitle} - ${import.meta.env.VITE_APP_TITLE}`;
  }, [location.pathname, menuTree]);

  return null;
}

export function AppRoutes() {
  const menuTree = usePermissionStore((state) => state.menuTree);
  const isMenuLoaded = usePermissionStore((state) => state.isMenuLoaded);

  const routes = useMemo<RouteObject[]>(() => {
    const dynamicRoutes = isMenuLoaded ? buildDynamicRoutes(flattenLeafMenus(menuTree)) : [];
    const directoryRedirectRoutes = isMenuLoaded
      ? buildDirectoryRedirectRoutes(menuTree)
      : [];

    return [
      {
        path: '/login',
        element: (
          <Suspense fallback={<PageLoading />}>
            <LoginPage />
          </Suspense>
        ),
      },
      {
        path: '/403',
        element: (
          <Suspense fallback={<PageLoading />}>
            <ForbiddenPage />
          </Suspense>
        ),
      },
      {
        path: '/404',
        element: (
          <Suspense fallback={<PageLoading />}>
            <NotFoundPage />
          </Suspense>
        ),
      },
      {
        path: '/500',
        element: (
          <Suspense fallback={<PageLoading />}>
            <ServerErrorPage />
          </Suspense>
        ),
      },
      {
        path: '/',
        element: <RequireAuth />,
        children: [
          {
            index: true,
            element: <Navigate replace to="/dashboard" />,
          },
          {
            path: 'dashboard',
            element: <DashboardPage />,
          },
          ...directoryRedirectRoutes,
          ...dynamicRoutes,
          {
            path: '*',
            element: <NotFoundPage />,
          },
        ],
      },
    ];
  }, [isMenuLoaded, menuTree]);

  const element = useRoutes(routes);

  return (
    <>
      <TitleSync />
      {element}
    </>
  );
}
