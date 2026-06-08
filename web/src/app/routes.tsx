/**
 * 路由表：集中声明 URL 与页面组件的映射及布局嵌套关系。
 */

import { Navigate } from 'react-router-dom';
import type { RouteObject } from 'react-router-dom';
import { RequireAuth } from './RequireAuth';
import { MainLayout } from './Layout/MainLayout';
import { LoginPage } from '../features/auth/pages/LoginPage';
import { AssetListPage } from '../features/asset/pages/AssetListPage';
import { CatalogTreePage } from '../features/catalog/pages/CatalogTreePage';
import { AuditCenterPage } from '../features/workflow/pages/AuditCenterPage';
import { WorkflowDefPage } from '../features/workflow/pages/WorkflowDefPage';
import { TranscodeGroupPage } from '../features/transcode/pages/TranscodeGroupPage';
import { TranscodeTaskPage } from '../features/transcode/pages/TranscodeTaskPage';
import { AiPlaceholderPage } from '../features/ai/AiPlaceholderPage';
import { UserListPage } from '../features/sys/user/pages/UserListPage';
import { RoleListPage } from '../features/sys/role/pages/RoleListPage';
import { OrgListPage } from '../features/sys/org/pages/OrgListPage';
import { MenuListPage } from '../features/sys/menu/pages/MenuListPage';
import { ConfigListPage } from '../features/sys/config/pages/ConfigListPage';

/**
 * 应用路由配置：/login 无 Layout；其余业务路由挂在 MainLayout 下并由 RequireAuth 保护。
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
        { path: 'catalog', element: <CatalogTreePage /> },
        { path: 'review', element: <AuditCenterPage /> },
        { path: 'workflow/def', element: <WorkflowDefPage /> },
        { path: 'transcode', element: <Navigate to="/transcode/task" replace /> },
        { path: 'transcode/group', element: <TranscodeGroupPage /> },
        { path: 'transcode/task', element: <TranscodeTaskPage /> },
        { path: 'sys/user', element: <UserListPage /> },
        { path: 'sys/role', element: <RoleListPage /> },
        { path: 'sys/org', element: <OrgListPage /> },
        { path: 'sys/menu', element: <MenuListPage /> },
        { path: 'sys/config', element: <ConfigListPage /> },
        { path: 'system', element: <Navigate to="/sys/user" replace /> },
        { path: 'ai', element: <AiPlaceholderPage /> },
      ],
    },
    { path: '*', element: <Navigate to="/asset/list" replace /> },
  ];
}