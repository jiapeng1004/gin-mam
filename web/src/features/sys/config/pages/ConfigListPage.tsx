/**
 * 系统配置列表占位。
 * 后续接口（规划）：GET /api/v1/sys/config/page；GET/PUT /api/v1/sys/config/:key
 */

import type { ColumnsType } from 'antd/es/table';
import { SysCrudShell } from '../../components/SysCrudShell';

const columns: ColumnsType<Record<string, unknown>> = [
  { title: '配置键', dataIndex: 'key', key: 'key' },
  { title: '配置值', dataIndex: 'value', key: 'value' },
  { title: '说明', dataIndex: 'remark', key: 'remark' },
];

/**
 * 系统配置占位页。
 * @returns 空表格与 API 说明
 */
export function ConfigListPage() {
  return (
    <SysCrudShell
      title="系统配置"
      description="规划分页 GET /api/v1/sys/config/page；按 key 读取 GET /api/v1/sys/config/:key；更新 PUT /api/v1/sys/config/:key（对接 ConfigService）。"
      columns={columns}
    />
  );
}