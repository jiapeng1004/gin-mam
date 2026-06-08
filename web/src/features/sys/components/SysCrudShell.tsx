/**
 * 系统管理 CRUD 占位壳：空表格 + 开发中说明，供角色/组织等模块复用。
 */

import { Alert, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';

export interface SysCrudShellProps {
  /** 页面标题 */
  title: string;
  /** 功能说明与规划 API 的简短文案（展示在页面上） */
  description: string;
  /** 表格列定义（通常无数据，仅展示表头） */
  columns: ColumnsType<Record<string, unknown>>;
}

/**
 * 渲染「功能开发中」说明与空数据表格，作为 sys 子模块最小 CRUD 壳。
 * @param props - 标题、说明与列
 * @returns 占位列表 UI
 */
export function SysCrudShell({ title, description, columns }: SysCrudShellProps) {
  return (
    <div>
      <Typography.Title level={4}>{title}</Typography.Title>
      <Alert type="info" message="功能开发中" description={description} showIcon style={{ marginBottom: 16 }} />
      <Table<Record<string, unknown>>
        rowKey="id"
        columns={columns}
        dataSource={[]}
        pagination={false}
        locale={{ emptyText: '暂无数据，接口接入后将在此展示' }}
      />
    </div>
  );
}