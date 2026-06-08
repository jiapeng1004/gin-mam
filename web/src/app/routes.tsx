/**
 * 路由表：集中声明 URL 与页面组件的映射及布局嵌套关系。
 */

import { Navigate } from 'react-router-dom';
import type { RouteObject } from 'react-router-dom';
import { RequireAuth } from './RequireAuth';
import { MainLayout } from './Layout/MainLayout';
import { LoginPage } from '../features/auth/LoginPage';
import { AssetListPage } from '../features/asset/AssetListPage';
import { CatalogPage } from '../features/catalog/CatalogPage';
import { ReviewPage } from '../features/review/ReviewPage';
import { TranscodePage } from '../features/transcode/TranscodePage';
import { SystemPage } from '../features/system/SystemPage';
import { AiPage } from '../features/ai/AiPage';

/**
 * 应用路由配置：/login 无 Layout；其余业务路由挂在 MainLayout 下并受 RequireAuth 保护。
 * @returns react-router RouteObject 数组
 */
export function buildRoutes(): RouteObject[] {
  return [
    {
      path: '/login',
      element: <LoginPage />,
    },
    {
      path: '/',
      element: (
        <RequireAuth>
          <MainLayout />
        </RequireAuth>
      ),
      children: [
        { index: true, element: <Navigate to="/asset/list" replace /> },
        { path: 'asset/list', element: <AssetListPage /> },
        { path: 'catalog', element: <CatalogPage /> },
        { path: 'review', element: <ReviewPage /> },
        { path: 'transcode', element: <TranscodePage /> },
        { path: 'system', element: <SystemPage /> },
        { path: 'ai', element: <AiPage /> },
      ],
    },
    { path: '*', element: <Navigate to="/asset/list" replace /> },
  ];
}