/**
 * 应用路由入口：挂载 BrowserRouter 与路由表。
 */

import { BrowserRouter, useRoutes } from 'react-router-dom';
import { buildRoutes } from './routes';

/**
 * 根据 buildRoutes 渲染匹配的路由节点。
 * @returns 当前路由对应的页面树
 */
function AppRoutes() {
  return useRoutes(buildRoutes());
}

/**
 * 根组件：提供前端路由上下文。
 * @returns 应用根 UI
 */
export default function App() {
  return (
    <BrowserRouter>
      <AppRoutes />
    </BrowserRouter>
  );
}