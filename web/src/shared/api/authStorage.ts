/**
 * 访问令牌持久化：读写 localStorage，供 axios 拦截器注入 Authorization。
 */

import { AUTH_TOKEN_KEY } from '../../app/auth';

/**
 * 读取当前保存的访问令牌。
 * @returns 令牌字符串；未登录时返回 null
 */
export function getToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

/**
 * 保存访问令牌。
 * @param token - 登录接口返回的 JWT 或访问令牌
 */
export function setToken(token: string): void {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

/**
 * 清除本地访问令牌（登出或 401 时使用）。
 */
export function clearToken(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}
