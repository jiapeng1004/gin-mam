/**
 * 转码任务页：任务分页列表与创建任务表单。
 */

import { Button, Form, Input, Modal, Select, Space, Table, Typography, message } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { useCallback, useEffect, useState } from 'react';
import {
  createTranscodeTask,
  fetchTranscodeGroupPage,
  fetchTranscodeTaskPage,
  type TranscodeGroupVO,
  type TranscodeTaskVO,
} from '../api/transcodeApi';

const DEFAULT_PAGE_SIZE = 20;

const TASK_STATUS_LABEL: Record<number, string> = {
  0: '待处理',
  1: '转码中',
  2: '成功',
  3: '失败',
};

/** 创建任务表单 */
interface CreateTaskFormValues {
  assetId: string;
  assetFileId?: string;
  transcodeGroupId: string;
}

/**
 * 转码任务列表与创建 UI。
 * @returns 转码任务页面
 */
export function TranscodeTaskPage() {
  const [form] = Form.useForm<CreateTaskFormValues>();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<TranscodeTaskVO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [groups, setGroups] = useState<TranscodeGroupVO[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async (p: number, ps: number) => {
    setLoading(true);
    try {
      const result = await fetchTranscodeTaskPage(p, ps);
      setData(result.list);
      setTotal(result.total);
      setPage(result.page);
      setPageSize(result.pageSize);
    } finally {
      setLoading(false);
    }
  }, []);

  const loadGroups = useCallback(async () => {
    const result = await fetchTranscodeGroupPage(1, 100);
    setGroups(result.list);
  }, []);

  useEffect(() => {
    void load(page, pageSize);
  }, [load, page, pageSize]);

  const openCreate = () => {
    form.resetFields();
    void loadGroups();
    setModalOpen(true);
  };

  const onSubmit = async () => {
    const values = await form.validateFields();
    setSubmitting(true);
    try {
      await createTranscodeTask({
        assetId: values.assetId,
        assetFileId: values.assetFileId || undefined,
        transcodeGroupId: values.transcodeGroupId,
      });
      message.success('转码任务已创建');
      setModalOpen(false);
      void load(1, pageSize);
      setPage(1);
    } finally {
      setSubmitting(false);
    }
  };

  const columns: ColumnsType<TranscodeTaskVO> = [
    { title: '任务 ID', dataIndex: 'id', key: 'id', ellipsis: true },
    { title: '媒资 ID', dataIndex: 'assetId', key: 'assetId', ellipsis: true },
    { title: '转码组', dataIndex: 'transcodeGroupId', key: 'transcodeGroupId', ellipsis: true },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (v: number) => TASK_STATUS_LABEL[v] ?? String(v),
    },
    { title: '输出路径', dataIndex: 'outputPath', key: 'outputPath', ellipsis: true },
    { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt' },
  ];

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? DEFAULT_PAGE_SIZE);
  };

  return (
    <div>
      <Space style={{ marginBottom: 16, width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          转码任务
        </Typography.Title>
        <Button type="primary" onClick={openCreate}>
          创建任务
        </Button>
      </Space>
      <Table<TranscodeTaskVO>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={data}
        pagination={{ current: page, pageSize, total, showSizeChanger: true }}
        onChange={onTableChange}
      />
      <Modal
        title="创建转码任务"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={() => void onSubmit()}
        confirmLoading={submitting}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="assetId" label="媒资 ID" rules={[{ required: true }]}>
            <Input placeholder="asset id" />
          </Form.Item>
          <Form.Item name="assetFileId" label="媒资文件 ID（可选）">
            <Input />
          </Form.Item>
          <Form.Item name="transcodeGroupId" label="转码组" rules={[{ required: true }]}>
            <Select
              placeholder="选择转码组"
              options={groups.map((g) => ({ value: g.id, label: g.name }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}