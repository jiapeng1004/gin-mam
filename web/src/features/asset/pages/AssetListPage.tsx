/**
 * 媒资列表页：分页表格、关键词搜索、跳转详情。
 */

import { Button, Input, Space, Table, Typography } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchAssetPage, type AssetVO } from '../api/assetApi';

const DEFAULT_PAGE_SIZE = 20;

function formatStatus(status: number): string {
  if (status === 1) return '已发布';
  if (status === 2) return '审核中';
  return '草稿';
}

/**
 * 媒资列表与检索。
 * @returns 媒资表格 UI
 */
export function AssetListPage() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<AssetVO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [keyword, setKeyword] = useState('');
  const [searchKeyword, setSearchKeyword] = useState('');

  const load = useCallback(async (p: number, ps: number, kw: string) => {
    setLoading(true);
    try {
      const result = await fetchAssetPage({
        page: p,
        pageSize: ps,
        keyword: kw.trim() || undefined,
      });
      setData(result.list);
      setTotal(result.total);
      setPage(result.page);
      setPageSize(result.pageSize);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(page, pageSize, searchKeyword);
  }, [load, page, pageSize, searchKeyword]);

  const columns: ColumnsType<AssetVO> = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (v: number) => formatStatus(v),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Button type="link" onClick={() => navigate(`/asset/${record.id}`)}>
          详情
        </Button>
      ),
    },
  ];

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? DEFAULT_PAGE_SIZE);
  };

  return (
    <div>
      <Typography.Title level={4}>媒资列表</Typography.Title>
      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search
          allowClear
          placeholder="搜索标题或关键词"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onSearch={(value) => {
            setSearchKeyword(value);
            setPage(1);
          }}
          style={{ width: 280 }}
        />
        <Button type="primary" onClick={() => navigate('/asset/upload')}>
          上传媒资
        </Button>
      </Space>
      <Table<AssetVO>
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