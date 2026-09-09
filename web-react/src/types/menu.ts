import type { SnowflakeId } from './api';

export interface MenuRecord {
  id: SnowflakeId;
  menuName: string;
  parentId: SnowflakeId;
  sort?: number;
  path: string;
  component?: string;
  query?: string;
  isFrame?: number;
  isCache?: number;
  menuType: number;
  visible?: number;
  status?: number;
  perms?: string;
  apiPermissionIds?: SnowflakeId[];
  icon?: string;
  remark?: string;
  children?: MenuRecord[];
}

export interface MenuRouteRecord extends MenuRecord {
  fullPath: string;
}
