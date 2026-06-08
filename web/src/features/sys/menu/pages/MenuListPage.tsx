/**
 * 菜单管理列表占位。
 * 后续接口（规划）：GET /api/v1/sys/menu/tree；POST/PUT/DELETE /api/v1/sys/menu/:id
 */

import type { ColumnsType } from 'antd/es/table';
import { SysCrudShell } from '../../components/SysCrudShell';

const columns: ColumnsType<Record<string, unknown>> = [
  { title: '菜单名称', dataIndex: 'title', key: 'title' },
  { title: '路由路径', dataIndex: 'path', key: 'path' },
  { title: '类型', dataIndex: 'type', key: 'type' },
];

/**
 * 菜单管理占位页。
 * @returns 空表格与 API 说明
 */
export function MenuListPage() {
  return (
    <SysCrudShell
      title="菜单管理"
      description="规划菜单树 GET /api/v1/sys/menu/tree；增删改 POST/PUT/DELETE /api/v1/sys/menu/:id，用于动态侧栏与权限绑定。"
      columns={columns}
    />
  );
}