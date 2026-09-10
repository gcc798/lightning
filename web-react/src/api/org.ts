import { request } from '@/utils/request';
import type { PageData, PageQuery, SnowflakeId } from '@/types/api';
import type { OrgRecord } from '@/types/system';

export const orgApi = {
  tree: () => request.get<OrgRecord[] | null>('/api/v1/org/tree').then((data) => data ?? []),
  page: (data: PageQuery & { orgName?: string; status?: number }) =>
    request.post<PageData<OrgRecord>>('/api/v1/org/page', data),
  detail: (id: SnowflakeId) => request.get<OrgRecord>(`/api/v1/org/${id}`),
  create: (data: Partial<OrgRecord>) => request.post<string>('/api/v1/org', data),
  update: (id: SnowflakeId, data: Partial<OrgRecord>) =>
    request.put<string>(`/api/v1/org/${id}`, data),
  delete: (id: SnowflakeId) => request.delete<string>(`/api/v1/org/${id}`),
  batchDelete: (ids: SnowflakeId[]) =>
    request.delete<string>('/api/v1/org/batch', { data: { ids } }),
};
