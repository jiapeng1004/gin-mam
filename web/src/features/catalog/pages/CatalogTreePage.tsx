/**
 * 编目树页：左侧树、右侧编目配置表单。
 */

import {
  Button,
  Form,
  Input,
  InputNumber,
  Space,
  Spin,
  Tree,
  Typography,
  message,
} from 'antd';
import type { DataNode } from 'antd/es/tree';
import { useCallback, useEffect, useMemo, useState, type Key } from 'react';
import {
  createCatalog,
  deleteCatalog,
  fetchCatalogConfig,
  fetchCatalogTree,
  putCatalogConfig,
  type CatalogConfigInput,
  type CatalogTreeNode,
} from '../api/catalogApi';

function toTreeData(nodes: CatalogTreeNode[]): DataNode[] {
  return nodes.map((node) => ({
    key: node.id,
    title: node.name,
    children: node.children?.length ? toTreeData(node.children) : undefined,
  }));
}

/**
 * 编目树与元数据配置维护。
 * @returns 树 + 配置表单 UI
 */
export function CatalogTreePage() {
  const [tree, setTree] = useState<CatalogTreeNode[]>([]);
  const [loadingTree, setLoadingTree] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [configLoading, setConfigLoading] = useState(false);
  const [configForm] = Form.useForm<{ items: CatalogConfigInput[] }>();
  const [createForm] = Form.useForm<{ name: string; sortCode?: number }>();

  const loadTree = useCallback(async () => {
    setLoadingTree(true);
    try {
      const data = await fetchCatalogTree();
      setTree(data);
    } catch {
      message.error('加载编目树失败');
    } finally {
      setLoadingTree(false);
    }
  }, []);

  useEffect(() => {
    void loadTree();
  }, [loadTree]);

  const treeData = useMemo(() => toTreeData(tree), [tree]);

  const loadConfig = async (catalogId: string) => {
    setConfigLoading(true);
    try {
      const items = await fetchCatalogConfig(catalogId);
      configForm.setFieldsValue({
        items: items.map((item) => ({
          configKey: item.configKey,
          configValue: item.configValue,
        })),
      });
    } catch {
      message.error('加载编目配置失败');
    } finally {
      setConfigLoading(false);
    }
  };

  const onSelect = (keys: Key[]) => {
    const id = keys[0] ? String(keys[0]) : null;
    setSelectedId(id);
    if (id) {
      void loadConfig(id);
    } else {
      configForm.resetFields();
    }
  };

  const saveConfig = async () => {
    if (!selectedId) {
      message.warning('请先选择编目节点');
      return;
    }
    const values = await configForm.validateFields();
    try {
      await putCatalogConfig(selectedId, values.items ?? []);
      message.success('配置已保存');
      void loadConfig(selectedId);
    } catch {
      message.error('保存配置失败');
    }
  };

  const onCreateRoot = async () => {
    const values = await createForm.validateFields();
    try {
      await createCatalog({
        name: values.name,
        parentId: selectedId,
        sortCode: values.sortCode ?? 0,
      });
      message.success('创建成功');
      createForm.resetFields();
      void loadTree();
    } catch {
      message.error('创建失败');
    }
  };

  const onDeleteSelected = async () => {
    if (!selectedId) {
      message.warning('请先选择编目节点');
      return;
    }
    try {
      await deleteCatalog(selectedId);
      message.success('已删除');
      setSelectedId(null);
      configForm.resetFields();
      void loadTree();
    } catch {
      message.error('删除失败');
    }
  };

  return (
    <div>
      <Typography.Title level={4}>编目管理</Typography.Title>
      <div style={{ display: 'flex', gap: 24, minHeight: 420 }}>
        <div style={{ width: 320, borderRight: '1px solid #f0f0f0', paddingRight: 16 }}>
          <Space direction="vertical" style={{ width: '100%', marginBottom: 12 }}>
            <Form form={createForm} layout="inline" onFinish={() => void onCreateRoot()}>
              <Form.Item name="name" rules={[{ required: true, message: '名称' }]}>
                <Input placeholder="新编目名称" />
              </Form.Item>
              <Form.Item name="sortCode">
                <InputNumber placeholder="排序" min={0} />
              </Form.Item>
              <Button type="primary" htmlType="submit">
                新建
              </Button>
            </Form>
            <Button danger disabled={!selectedId} onClick={() => void onDeleteSelected()}>
              删除选中
            </Button>
          </Space>
          <Spin spinning={loadingTree}>
            <Tree
              treeData={treeData}
              selectedKeys={selectedId ? [selectedId] : []}
              onSelect={onSelect}
            />
          </Spin>
        </div>
        <div style={{ flex: 1 }}>
          <Typography.Title level={5}>编目配置</Typography.Title>
          {!selectedId ? (
            <Typography.Text type="secondary">请在左侧选择编目节点</Typography.Text>
          ) : (
            <Form form={configForm} layout="vertical" disabled={configLoading}>
              <Form.List name="items">
                {(fields, { add, remove }) => (
                  <>
                    {fields.map((field) => (
                      <Space key={field.key} align="baseline" style={{ display: 'flex', marginBottom: 8 }}>
                        <Form.Item
                          {...field}
                          name={[field.name, 'configKey']}
                          rules={[{ required: true, message: '键' }]}
                        >
                          <Input placeholder="configKey" />
                        </Form.Item>
                        <Form.Item {...field} name={[field.name, 'configValue']}>
                          <Input placeholder="configValue" style={{ width: 240 }} />
                        </Form.Item>
                        <Button type="link" onClick={() => remove(field.name)}>
                          删除
                        </Button>
                      </Space>
                    ))}
                    <Button type="dashed" onClick={() => add({ configKey: '', configValue: '' })} block>
                      添加配置项
                    </Button>
                  </>
                )}
              </Form.List>
              <Button type="primary" onClick={() => void saveConfig()} style={{ marginTop: 16 }}>
                保存配置
              </Button>
            </Form>
          )}
        </div>
      </div>
    </div>
  );
}
