/**
 * 流程定义页：分页展示 gm_workflow_def，并提供创建表单。
 */

import { Button, Form, Input, InputNumber, Space, Table, Typography, message } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { useCallback, useEffect, useState } from 'react';
import {
  createWorkflowDef,
  fetchWorkflowDefPage,
  type WorkflowDefVO,
} from '../api/workflowApi';

const DEFAULT_PAGE_SIZE = 20;

/** 创建表单字段结构 */
interface CreateDefFormValues {
  name: string;
  auditLevel: number;
  levelUsers: { level: number; userId: string }[];
}

/**
 * 流程定义列表与创建 UI。
 * @returns 流程定义管理页面
 */
export function WorkflowDefPage() {
  const [form] = Form.useForm<CreateDefFormValues>();
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [data, setData] = useState<WorkflowDefVO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);

  const load = useCallback(async (p: number, ps: number) => {
    setLoading(true);
    try {
      const result = await fetchWorkflowDefPage(p, ps);
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

  const columns: ColumnsType<WorkflowDefVO> = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '审核层级', dataIndex: 'auditLevel', key: 'auditLevel' },
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

  const onFinish = async (values: CreateDefFormValues) => {
    setSubmitting(true);
    try {
      await createWorkflowDef({
        name: values.name,
        auditLevel: values.auditLevel,
        levelUsers: values.levelUsers ?? [],
      });
      message.success('流程定义已创建');
      form.resetFields();
      void load(1, pageSize);
      setPage(1);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Typography.Title level={4}>流程定义</Typography.Title>
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        initialValues={{ auditLevel: 1, levelUsers: [{ level: 1, userId: '' }] }}
        style={{ marginBottom: 24, maxWidth: 640 }}
      >
        <Form.Item name="name" label="流程名称" rules={[{ required: true, message: '请输入名称' }]}>
          <Input placeholder="例如：标准两级审核" />
        </Form.Item>
        <Form.Item name="auditLevel" label="审核层级" rules={[{ required: true }]}>
          <InputNumber min={1} max={10} style={{ width: '100%' }} />
        </Form.Item>
        <Typography.Text type="secondary">各级审核人（level 从 1 起）</Typography.Text>
        <Form.List name="levelUsers">
          {(fields, { add, remove }) => (
            <>
              {fields.map(({ key, name, ...rest }) => (
                <Space key={key} align="baseline" style={{ display: 'flex', marginTop: 8 }}>
                  <Form.Item
                    {...rest}
                    name={[name, 'level']}
                    rules={[{ required: true, message: '层级' }]}
                  >
                    <InputNumber min={1} placeholder="层级" />
                  </Form.Item>
                  <Form.Item
                    {...rest}
                    name={[name, 'userId']}
                    rules={[{ required: true, message: '用户 ID' }]}
                  >
                    <Input placeholder="审核人用户 ID" style={{ width: 220 }} />
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} />
                </Space>
              ))}
              <Form.Item>
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  添加审核人
                </Button>
              </Form.Item>
            </>
          )}
        </Form.List>
        <Button type="primary" htmlType="submit" loading={submitting}>
          创建流程定义
        </Button>
      </Form>
      <Table<WorkflowDefVO>
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