/**
 * 认证相关常量与简易守卫辅助：脚手架占位，正式鉴权在后续任务实现。
 */

/** localStorage 中存放访问令牌的键名 */
export const AUTH_TOKEN_KEY = 'gin_mam_auth_token';

/**
 * 判断是否已登录（脚手架：仅检查本地是否存在 token）。
 * @returns 有 token 视为已登录
 */
export function isAuthenticated(): boolean {
  return Boolean(localStorage.getItem(AUTH_TOKEN_KEY));
}