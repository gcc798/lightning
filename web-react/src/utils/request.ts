import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios';
import JSONbig from 'json-bigint';
import { message } from '@/components/ui';
import { DEFAULT_CLIENT_ID } from '@/constants/auth';
import type { CommonResp } from '@/types/api';
import { useAuthStore } from '@/store/auth';

const jsonBig = JSONbig({ storeAsString: true });
const authWhiteList = ['/login', '/auth/refresh', '/captcha'];

let refreshPromise: Promise<void> | null = null;

function shouldSkipRefresh(url?: string) {
  return authWhiteList.some((item) => url?.includes(item));
}

function redirectToLogin() {
  const target = `${window.location.pathname}${window.location.search}`;
  const redirect = encodeURIComponent(target);
  window.location.href = `/login?redirect=${redirect}`;
}

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
  transformResponse: [
    (data) => {
      if (typeof data !== 'string') {
        return data;
      }

      try {
        return jsonBig.parse(data);
      } catch {
        return data;
      }
    },
  ],
});

service.interceptors.request.use((config) => {
  const { accessToken } = useAuthStore.getState();

  config.headers.clientId = DEFAULT_CLIENT_ID;

  // 所有业务请求统一在这里补 Authorization，避免页面层重复处理登录态。
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }

  return config;
});

service.interceptors.response.use(
  async (response) => {
    if (
      response.config.responseType === 'blob' ||
      response.config.responseType === 'arraybuffer'
    ) {
      return response.data;
    }

    const payload = response.data as CommonResp<unknown>;

    if (!payload || typeof payload.code !== 'number') {
      return response.data;
    }

    // sys-api 的绝大多数接口都返回 CommonResp，这里统一拆掉外层 data，
    // 让页面和 API 模块直接面对真正的业务数据结构。
    if (payload.code === 200) {
      return payload.data;
    }

    // 401 刷新逻辑集中在请求层，页面只关心“当前调用是否成功”。
    // 这样可以避免每个页面都重复写 token 续期与跳转逻辑。
    if (payload.code === 401 && !shouldSkipRefresh(response.config.url)) {
      const authStore = useAuthStore.getState();

      if (!authStore.refreshToken) {
        authStore.clearAuthState();
        message.error(payload.msg || '登录已过期，请重新登录');
        redirectToLogin();
        return Promise.reject(new Error(payload.msg || '登录已过期'));
      }

      // 并发请求只允许共享一个 refresh Promise，避免多个 401 同时触发多次刷新。
      if (!refreshPromise) {
        refreshPromise = authStore
          .refreshAccessToken()
          .finally(() => {
            refreshPromise = null;
          });
      }

      try {
        await refreshPromise;

        // 刷新成功后使用新的 accessToken 重放原始请求。
        const retryConfig = {
          ...response.config,
          headers: {
            ...response.config.headers,
            Authorization: `Bearer ${useAuthStore.getState().accessToken}`,
          },
        };

        return service.request(retryConfig);
      } catch (error) {
        authStore.clearAuthState();
        message.error(payload.msg || '登录已过期，请重新登录');
        redirectToLogin();
        return Promise.reject(error);
      }
    }

    message.error(payload.msg || '请求失败');

    const error = new Error(payload.msg || '请求失败') as Error & { code?: number };
    error.code = payload.code;
    return Promise.reject(error);
  },
  (error) => {
    // 这里处理真正的网络错误或 HTTP 层异常，和上面的业务 code 分支互补。
    if (error.response?.status === 401 && !shouldSkipRefresh(error.config?.url)) {
      useAuthStore.getState().clearAuthState();
      redirectToLogin();
    }

    if (error.response?.data?.msg) {
      message.error(error.response.data.msg);
    } else if (error.message) {
      message.error(error.message);
    } else {
      message.error('请求失败');
    }

    return Promise.reject(error);
  },
);

export const request = {
  get<T = unknown>(url: string, config?: AxiosRequestConfig) {
    return service.get<T, T>(url, config);
  },

  post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    return service.post<T, T>(url, data, config);
  },

  put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    return service.put<T, T>(url, data, config);
  },

  delete<T = unknown>(url: string, config?: AxiosRequestConfig) {
    return service.delete<T, T>(url, config);
  },
};
