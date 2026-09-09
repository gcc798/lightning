import type { ItemType } from 'antd/es/menu/interface';
import type { MenuRecord, MenuRouteRecord } from '@/types/menu';
import { getMenuIconNode } from './icons';
import { isNumericValue, toNumberValue, toOptionalNumber } from './number';

function toStringValue(value: unknown) {
  return typeof value === 'string' ? value : '';
}

export function normalizeMenuRecord(raw: unknown): MenuRecord {
  const source = (raw ?? {}) as Record<string, unknown>;
  const children = normalizeMenuTree(source.children);
  const menu: MenuRecord = {
    id: (source.id ?? '') as MenuRecord['id'],
    menuName: toStringValue(source.menuName),
    parentId: (source.parentId ?? 0) as MenuRecord['parentId'],
    sort: toOptionalNumber(source.sort),
    path: toStringValue(source.path),
    component: toStringValue(source.component),
    query: toStringValue(source.query),
    isFrame: toOptionalNumber(source.isFrame),
    isCache: toOptionalNumber(source.isCache),
    menuType: toNumberValue(source.menuType),
    visible: toOptionalNumber(source.visible),
    status: toOptionalNumber(source.status),
    perms: toStringValue(source.perms),
    apiPermissionIds: Array.isArray(source.apiPermissionIds)
      ? source.apiPermissionIds as MenuRecord['apiPermissionIds']
      : [],
    icon: toStringValue(source.icon),
    remark: toStringValue(source.remark),
  };

  if (children.length > 0) {
    menu.children = children;
  }

  return menu;
}

export function normalizeMenuTree(raw: unknown): MenuRecord[] {
  if (!Array.isArray(raw)) {
    return [];
  }

  return raw.map(normalizeMenuRecord);
}

export function joinMenuPath(parentPath: string, currentPath: string) {
  if (!currentPath) {
    return parentPath || '/';
  }

  if (currentPath.startsWith('/')) {
    return currentPath;
  }

  const normalizedParent = parentPath === '/' ? '' : parentPath.replace(/\/$/, '');
  return `${normalizedParent}/${currentPath}`.replace(/\/+/g, '/');
}

export function isMenuHidden(menu: MenuRecord) {
  return isNumericValue(menu.visible, 1) || isNumericValue(menu.status, 1);
}

export function findFirstNavigablePath(
  menuTree: MenuRecord[],
  parentPath = '',
): string | null {
  for (const menu of menuTree) {
    if (isMenuHidden(menu) || isNumericValue(menu.menuType, 2)) {
      continue;
    }

    const fullPath = joinMenuPath(parentPath, menu.path);

    if (isNumericValue(menu.menuType, 1)) {
      return fullPath;
    }

    if (menu.children?.length) {
      const childPath = findFirstNavigablePath(menu.children, fullPath);
      if (childPath) {
        return childPath;
      }
    }
  }

  return null;
}

export function extractPermissions(menuTree: MenuRecord[]) {
  const permissions: string[] = [];

  const walk = (nodes: MenuRecord[]) => {
    nodes.forEach((node) => {
      if (node.perms) {
        permissions.push(node.perms);
      }

      if (node.children?.length) {
        walk(node.children);
      }
    });
  };

  walk(menuTree);
  return permissions;
}

export function flattenLeafMenus(
  menuTree: MenuRecord[],
  parentPath = '',
): MenuRouteRecord[] {
  const routes: MenuRouteRecord[] = [];

  menuTree.forEach((menu) => {
    if (isMenuHidden(menu)) {
      return;
    }

    const fullPath = joinMenuPath(parentPath, menu.path);

    if (isNumericValue(menu.menuType, 1)) {
      routes.push({ ...menu, fullPath });
    }

    if (menu.children?.length) {
      routes.push(...flattenLeafMenus(menu.children, fullPath));
    }
  });

  return routes;
}

export function buildSidebarMenus(
  menuTree: MenuRecord[],
  parentPath = '',
): ItemType[] {
  return menuTree
    .filter((menu) => !isMenuHidden(menu) && !isNumericValue(menu.menuType, 2))
    .map((menu) => {
      const fullPath = joinMenuPath(parentPath, menu.path);
      const icon = getMenuIconNode(menu.icon);

      if (menu.children?.length) {
        return {
          key: fullPath,
          icon,
          label: menu.menuName,
          children: buildSidebarMenus(menu.children, fullPath),
        } satisfies ItemType;
      }

      return {
        key: fullPath,
        icon,
        label: menu.menuName,
      } satisfies ItemType;
    });
}
