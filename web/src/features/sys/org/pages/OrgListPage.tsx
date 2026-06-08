/**
 * 组织管理列表占位。
 * 后续接口（规划）：GET /api/v1/sys/org/tree 或 /page；POST/PUT/DELETE /api/v1/sys/org/:id
 */

import type { ColumnsType } from 'antd/es/table';
import { SysCrudShell } from '../../components/SysCrudShell';

const columns: ColumnsType<Record<string, unknown>> = [
  { title: '组织名称', dataIndex: 'name', key: 'name' },
  { title: '上级组织', dataIndex: 'parentName', key: 'parentName' },
  { title: '排序', dataIndex: 'sort', key: 'sort' },
];

/**
 * 组织管理占位页。
 * @returns 空表格与 API 说明
 */
export function OrgListPage() {
  return (
    <SysCrudShell
      title="组织管理"
      description="规划树形 GET /api/v1/sys/org/tree（或分页 /page）；维护组织 POST/PUT/DELETE /api/v1/sys/org/:id。"
      columns={columns}
    />
  );
}