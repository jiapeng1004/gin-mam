/**
 * 系统用户列表页：分页展示 gm_user，对接用户分页 API。
 */

import { Table, Typography } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { useCallback, useEffect, useState } from 'react';
import { fetchUserPage, type UserVO } from '../api/userApi';

const DEFAULT_PAGE_SIZE = 20;

/**
 * 用户管理列表，Table 分页加载 /sys/user/page 数据。
 * @returns 用户列表 UI
 */
export function UserListPage() {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<UserVO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);

  const load = useCallback(async (p: number, ps: number) => {
    setLoading(true);
    try {
      const result = await fetchUserPage(p, ps);
      setData(result.list);
      setTotal(result.total);
      setPage(result.page);
      setPageSize(result.pageSize);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(page, pageSize);
  }, [load, page, pageSize]);

  const columns: ColumnsType<UserVO> = [
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname', render: (v) => v || '—' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => (status === 1 ? '启用' : '停用'),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  ];

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? DEFAULT_PAGE_SIZE);
  };

  return (
    <div>
      <Typography.Title level={4}>用户管理</Typography.Title>
      <Table<UserVO>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={data}
        pagination={{ current: page, pageSize, total, showSizeChanger: true }}
        onChange={onTableChange}
      />
    </div>
  );
}