/**
 * HTTP 客户端：基于 axios，对齐后端 httpx 无包装 JSON 与 ErrorVo 错误体。
 */

import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { clearToken, getToken } from './authStorage';
import { ApiError } from './errors';

/** 401 未授权时跳转登录页的路径 */
const LOGIN_PATH = '/login';

/**
 * 未授权时将浏览器导航至登录页（单独导出便于单测 mock）。
 */
export function redirectToLogin(): void {
  window.location.assign(LOGIN_PATH);
}

/**
 * 共享 axios 实例：开发环境经 Vite 代理访问 /api/v1。
 */
export const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
});

apiClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getToken();
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`);
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => {
    // 2xx：httpx.OK/Created 直接 JSON 序列化业务对象，无 code/data 外层
    return response.data;
  },
  (error: AxiosError<{ err_code?: number; err_msg?: string }>) => {
    const status = error.response?.status;
    const body = error.response?.data;

    if (status === 401) {
      clearToken();
      redirectToLogin();
    }

    // 4xx/5xx：httpx.Fail 写入 ErrorVo { err_code, err_msg }
    if (body && typeof body.err_code === 'number') {
      return Promise.reject(new ApiError(body.err_code, body.err_msg ?? ''));
    }

    const fallbackMsg = error.message || '请求失败';
    const fallbackCode = status ?? 0;
    return Promise.reject(new ApiError(fallbackCode, fallbackMsg));
  },
);

/**
 * GET 请求，成功时返回响应体业务数据。
 */
export async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  return (await apiClient.get<T>(url, { params })) as T;
}

/**
 * POST 请求，成功时返回响应体业务数据。
 */
export async function post<T>(url: string, body?: unknown): Promise<T> {
  return (await apiClient.post<T>(url, body)) as T;
}

/**
 * PUT 请求，成功时返回响应体业务数据。
 */
export async function put<T>(url: string, body?: unknown): Promise<T> {
  return (await apiClient.put<T>(url, body)) as T;
}

/**
 * DELETE 请求，成功时返回响应体业务数据。
 */
export async function del<T>(url: string): Promise<T> {
  return (await apiClient.delete<T>(url)) as T;
}
