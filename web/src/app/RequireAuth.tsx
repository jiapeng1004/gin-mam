/**
 * 路由守卫：未登录用户访问受保护路由时重定向到登录页。
 */

import { Navigate, useLocation } from 'react-router-dom';
import type { ReactNode } from 'react';
import { isAuthenticated } from './auth';

export interface RequireAuthProps {
  /** 需要登录才能渲染的子树（通常为 MainLayout 或业务页） */
  children: ReactNode;
}

/**
 * 包裹受保护路由；无会话时带上来源 location 跳转 /login。
 * @param props - 子组件
 * @returns 子组件或重定向
 */
export function RequireAuth({ children }: RequireAuthProps) {
  const location = useLocation();

  if (!isAuthenticated()) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
}