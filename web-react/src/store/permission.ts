import { create } from 'zustand';
import { menuApi } from '@/api/menu';
import type { MenuRecord } from '@/types/menu';
import { extractPermissions, normalizeMenuTree } from '@/utils/menu';
import { hasPermission as checkPermission } from '@/utils/permissions';

interface PermissionState {
  menuTree: MenuRecord[];
  permissions: string[];
  isMenuLoaded: boolean;
  isMenuLoading: boolean;
  menuLoadError: string;
  loadMenuTree: () => Promise<void>;
  reset: () => void;
  hasPermission: (permission?: string | string[]) => boolean;
}

export const usePermissionStore = create<PermissionState>()(
  (set, get) => ({
    menuTree: [],
    permissions: [],
    isMenuLoaded: false,
    isMenuLoading: false,
    menuLoadError: '',
    loadMenuTree: async () => {
      if (get().isMenuLoading) {
        return;
      }

      set({ isMenuLoading: true, menuLoadError: '' });

      try {
        const menuTree = normalizeMenuTree(await menuApi.getUserMenuTree());
        set({
          menuTree,
          permissions: extractPermissions(menuTree),
          isMenuLoaded: true,
          isMenuLoading: false,
          menuLoadError: '',
        });
      } catch (error) {
        set({
          menuTree: [],
          permissions: [],
          // 即使菜单加载失败，也要把 isMenuLoaded 置为 true，
          // 否则路由守卫会一直停留在“正在加载菜单”的状态里死循环。
          isMenuLoaded: true,
          isMenuLoading: false,
          menuLoadError:
            error instanceof Error ? error.message : '加载菜单失败',
        });
      }
    },
    reset: () =>
      set({
        menuTree: [],
        permissions: [],
        isMenuLoaded: false,
        isMenuLoading: false,
        menuLoadError: '',
      }),
    hasPermission: (permission) => checkPermission(get().permissions, permission),
  }),
);
