/**
 * 转码组管理页：分页列表与创建、编辑、删除。
 */

import { Button, Form, Input, InputNumber, Modal, Popconfirm, Space, Table, Typography, message } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { useCallback, useEffect, useState } from 'react';
import {
  createTranscodeGroup,
  deleteTranscodeGroup,
  fetchTranscodeGroupPage,
  updateTranscodeGroup,
  type TranscodeGroupVO,
} from '../api/transcodeApi';

const DEFAULT_PAGE_SIZE = 20;

/** 转码组表单字段 */
interface GroupFormValues {
  name: string;
  message?: string;
  groupType?: number;
  param?: string;
  strategyType?: number;
  defaultFlag?: number;
  availableFlag?: number;
}

/**
 * 转码组 CRUD 页面。
 * @returns 转码组管理 UI
 */
export function TranscodeGroupPage() {
  const [form] = Form.useForm<GroupFormValues>();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<TranscodeGroupVO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<TranscodeGroupVO | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async (p: number, ps: number) => {
    setLoading(true);
    try {
      const result = await fetchTranscodeGroupPage(p, ps);
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

  const openCreate = () => {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({ groupType: 0, strategyType: 0, defaultFlag: 0, availableFlag: 1 });
    setModalOpen(true);
  };

  const openEdit = (row: TranscodeGroupVO) => {
    setEditing(row);
    form.setFieldsValue({
      name: row.name,
      message: row.message,
      groupType: row.groupType,
      param: row.param,
      strategyType: row.strategyType,
      defaultFlag: row.defaultFlag,
      availableFlag: row.availableFlag,
    });
    setModalOpen(true);
  };

  const onSubmit = async () => {
    const values = await form.validateFields();
    setSubmitting(true);
    try {
      if (editing) {
        await updateTranscodeGroup(editing.id, values);
        message.success('已更新');
      } else {
        await createTranscodeGroup(values);
        message.success('已创建');
      }
      setModalOpen(false);
      void load(page, pageSize);
    } finally {
      setSubmitting(false);
    }
  };

  const onDelete = async (id: string) => {
    await deleteTranscodeGroup(id);
    message.success('已删除');
    void load(page, pageSize);
  };

  const columns: ColumnsType<TranscodeGroupVO> = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '类型', dataIndex: 'groupType', key: 'groupType' },
    {
      title: '可用',
      dataIndex: 'availableFlag',
      key: 'availableFlag',
      render: (v: number) => (v === 1 ? '是' : '否'),
    },
    { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_, row) => (
        <Space>
          <Button type="link" onClick={() => openEdit(row)}>
            编辑
          </Button>
          <Popconfirm title="确认删除？" onConfirm={() => void onDelete(row.id)}>
            <Button type="link" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? DEFAULT_PAGE_SIZE);
  };

  return (
    <div>
      <Space style={{ marginBottom: 16, width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          转码组
        </Typography.Title>
        <Button type="primary" onClick={openCreate}>
          新建转码组
        </Button>
      </Space>
      <Table<TranscodeGroupVO>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={data}
        pagination={{ current: page, pageSize, total, showSizeChanger: true }}
        onChange={onTableChange}
      />
      <Modal
        title={editing ? '编辑转码组' : '新建转码组'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={() => void onSubmit()}
        confirmLoading={submitting}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="message" label="说明">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="groupType" label="组类型">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="strategyType" label="策略类型">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="param" label="参数 JSON">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="defaultFlag" label="默认组">
            <InputNumber min={0} max={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="availableFlag" label="是否可用">
            <InputNumber min={0} max={1} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}