/**
 * API 与全局常量：集中维护前端请求基址，便于环境与代理切换。
 */

/** 开发环境通过 Vite proxy 转发；生产环境由网关或同源 /api 提供 */
export const API_BASE = "/api";
