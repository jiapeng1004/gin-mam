/**
 * 认证 Hook：读取本地 token 状态并提供登出。
 */

import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { isAuthenticated as checkAuthenticated } from '../../../app/auth';
import { clearToken } from '../../../shared/api/authStorage';

/**
 * 提供当前是否已登录与登出方法（清除 token 并跳转登录页）。
 * @returns isAuthenticated、logout
 */
export function useAuth() {
  const navigate = useNavigate();

  const isAuthenticated = checkAuthenticated();

  const logout = useCallback(() => {
    clearToken();
    navigate('/login', { replace: true });
  }, [navigate]);

  return { isAuthenticated, logout };
}