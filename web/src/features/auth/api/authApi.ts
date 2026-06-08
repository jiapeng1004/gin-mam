/**
 * 认证 API：对接后端 POST /api/v1/auth/login。
 */

import { post } from '../../../shared/api/client';

/** 登录请求体，对应后端 loginRequest */
export interface LoginRequest {
  username: string;
  password: string;
}

/** 登录成功响应，camelCase：token、expiresAt（ISO 8601 字符串） */
export interface LoginResponse {
  token: string;
  expiresAt: string;
}

/**
 * 用户名密码登录，获取 JWT。
 * @param body - username、password
 * @returns token 与过期时间
 */
export async function login(body: LoginRequest): Promise<LoginResponse> {
  return post<LoginResponse>('/auth/login', body);
}