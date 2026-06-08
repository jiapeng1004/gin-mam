/**
 * 审核中心页：四个 Tab 展示不同视角的流程实例，待审列表支持审核弹窗。
 */

import { Button, Form, Input, Modal, Radio, Space, Table, Tabs, Typography, message } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  auditWorkflow,
  fetchAllInstances,
  fetchAssignedToMe,
  fetchAuditedByMe,
  fetchCreatedByMe,
  type WorkflowInstanceVO,
} from '../api/workflowApi';

const DEFAULT_PAGE_SIZE = 20;

type AuditTabKey = 'created' | 'assigned' | 'audited' | 'all';

const AUDIT_STATUS_LABEL: Record<number, string> = {
  0: '待审',
  1: '审核中',
  2: '已通过',
  3: '已打回',
};

/** 各 Tab 对应的分页查询函数 */
const TAB_FETCHERS: Record<
  AuditTabKey,
  (page: number, pageSize: number) => ReturnType<typeof fetchCreatedByMe>
> = {
  created: fetchCreatedByMe,
  assigned: fetchAssignedToMe,
  audited: fetchAuditedByMe,
  all: fetchAllInstances,
};

/**
 * 审核中心：多 Tab 实例列表与审核操作。
 * @returns 审核中心 UI
 */
export function AuditCenterPage() {
  const [activeTab, setActiveTab] = useState<AuditTabKey>('assigned');
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<WorkflowInstanceVO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [auditOpen, setAuditOpen] = useState(false);
  const [auditTarget, setAuditTarget] = useState<WorkflowInstanceVO | null>(null);
  const [auditSubmitting, setAuditSubmitting] = useState(false);
  const [auditForm] = Form.useForm<{ pass: boolean; remark?: string }>();

  const load = useCallback(
    async (tab: AuditTabKey, p: number, ps: number) => {
      setLoading(true);
      try {
        const fetcher = TAB_FETCHERS[tab];
        const result = await fetcher(p, ps);
        setData(result.list);
        setTotal(result.total);
        setPage(result.page);
        setPageSize(result.pageSize);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  useEffect(() => {
    void load(activeTab, page, pageSize);
  }, [activeTab, load, page, pageSize]);

  const openAudit = (row: WorkflowInstanceVO) => {
    setAuditTarget(row);
    auditForm.setFieldsValue({ pass: true, remark: '' });
    setAuditOpen(true);
  };

  const submitAudit = async () => {
    if (!auditTarget) {
      return;
    }
    const values = await auditForm.validateFields();
    setAuditSubmitting(true);
    try {
      await auditWorkflow({
        instanceId: auditTarget.id,
        pass: values.pass,
        remark: values.remark,
      });
      message.success(values.pass ? '已通过' : '已打回');
      setAuditOpen(false);
      setAuditTarget(null);
      void load(activeTab, page, pageSize);
    } finally {
      setAuditSubmitting(false);
    }
  };

  const columns: ColumnsType<WorkflowInstanceVO> = useMemo(() => {
    const base: ColumnsType<WorkflowInstanceVO> = [
      { title: '实例 ID', dataIndex: 'id', key: 'id', ellipsis: true },
      { title: '媒资 ID', dataIndex: 'assetId', key: 'assetId', ellipsis: true },
      { title: '流程定义', dataIndex: 'workflowDefId', key: 'workflowDefId', ellipsis: true },
      {
        title: '进度',
        key: 'progress',
        render: (_, row) => `${row.level}/${row.auditLevel}`,
      },
      {
        title: '状态',
        dataIndex: 'auditStatus',
        key: 'auditStatus',
        render: (v: number) => AUDIT_STATUS_LABEL[v] ?? String(v),
      },
      { title: '发起人', dataIndex: 'createdBy', key: 'createdBy' },
      { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt' },
    ];
    if (activeTab === 'assigned') {
      base.push({
        title: '操作',
        key: 'actions',
        render: (_, row) => (
          <Button type="link" onClick={() => openAudit(row)}>
            审核
          </Button>
        ),
      });
    }
    return base;
  }, [activeTab]);

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? DEFAULT_PAGE_SIZE);
  };

  const onTabChange = (key: string) => {
    setActiveTab(key as AuditTabKey);
    setPage(1);
  };

  const table = (
    <Table<WorkflowInstanceVO>
      rowKey="id"
      loading={loading}
      columns={columns}
      dataSource={data}
      pagination={{ current: page, pageSize, total, showSizeChanger: true }}
      onChange={onTableChange}
    />
  );

  return (
    <div>
      <Typography.Title level={4}>审核中心</Typography.Title>
      <Tabs
        activeKey={activeTab}
        onChange={onTabChange}
        items={[
          { key: 'created', label: '我发起的', children: table },
          { key: 'assigned', label: '待我审', children: table },
          { key: 'audited', label: '我审过的', children: table },
          { key: 'all', label: '全部', children: table },
        ]}
      />
      <Modal
        title="审核"
        open={auditOpen}
        onCancel={() => setAuditOpen(false)}
        onOk={() => void submitAudit()}
        confirmLoading={auditSubmitting}
        destroyOnClose
      >
        <Space direction="vertical" style={{ width: '100%' }}>
          <Typography.Text type="secondary">
            实例：{auditTarget?.id} / 媒资：{auditTarget?.assetId}
          </Typography.Text>
          <Form form={auditForm} layout="vertical">
            <Form.Item name="pass" label="审核结果" rules={[{ required: true }]}>
              <Radio.Group>
                <Radio value={true}>通过</Radio>
                <Radio value={false}>打回</Radio>
              </Radio.Group>
            </Form.Item>
            <Form.Item name="remark" label="备注">
              <Input.TextArea rows={3} placeholder="可选" />
            </Form.Item>
          </Form>
        </Space>
      </Modal>
    </div>
  );
}