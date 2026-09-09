import { request } from '@/utils/request';
import type { SnowflakeId } from '@/types/api';
import type { ApiPermissionRecord } from '@/types/system';

export interface RolePermissionGrants {
  manualIds: SnowflakeId[];
  derivedIds: SnowflakeId[];
}

export const apiPermissionApi = {
  tree: () => request.get<ApiPermissionRecord[]>('/api/v1/api-permission/tree'),
  list: () => request.get<ApiPermissionRecord[]>('/api/v1/api-permission'),
  create: (data: Partial<ApiPermissionRecord>) => request.post<ApiPermissionRecord>('/api/v1/api-permission', data),
  update: (id: SnowflakeId, data: Partial<ApiPermissionRecord>) =>
    request.put<string>(`/api/v1/api-permission/${id}`, data),
  delete: (id: SnowflakeId) => request.delete<string>(`/api/v1/api-permission/${id}`),
  getRolePermissions: (roleId: SnowflakeId) =>
    request.get<RolePermissionGrants>(`/api/v1/role/${roleId}/api-permissions`),
  assignRolePermissions: (roleId: SnowflakeId, permissionIds: SnowflakeId[]) =>
    request.post<string>(`/api/v1/role/${roleId}/api-permissions`, { permissionIds }),
  getUserPermissions: (userId: SnowflakeId) =>
    request.get<SnowflakeId[]>(`/api/v1/user/${userId}/api-permissions`),
  assignUserPermissions: (userId: SnowflakeId, permissionIds: SnowflakeId[]) =>
    request.post<string>(`/api/v1/user/${userId}/api-permissions`, { permissionIds }),
};
