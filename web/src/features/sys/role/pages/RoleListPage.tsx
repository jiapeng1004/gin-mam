/**
 * 角色管理列表占位。
 * 后续接口（规划）：GET /api/v1/sys/role/page；POST/PUT/DELETE /api/v1/sys/role/:id
 */

import type { ColumnsType } from 'antd/es/table';
import { SysCrudShell } from '../../components/SysCrudShell';

const columns: ColumnsType<Record<string, unknown>> = [
  { title: '角色编码', dataIndex: 'code', key: 'code' },
  { title: '角色名称', dataIndex: 'name', key: 'name' },
  { title: '状态', dataIndex: 'status', key: 'status' },
];

/**
 * 角色管理占位页。
 * @returns 空表格与 API 说明
 */
export function RoleListPage() {
  return (
    <SysCrudShell
      title="角色管理"
      description="规划分页 GET /api/v1/sys/role/page；创建 POST /api/v1/sys/role；更新 PUT /api/v1/sys/role/:id；删除 DELETE /api/v1/sys/role/:id。"
      columns={columns}
    />
  );
}