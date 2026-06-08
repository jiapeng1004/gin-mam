/**
 * 媒资详情页：Descriptions 展示 previewUrl 等字段。
 */

import { Button, Descriptions, Image, Spin, Typography, message } from 'antd';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { fetchAssetById, type AssetVO } from '../api/assetApi';

/**
 * 单条媒资详情。
 * @returns 详情 Descriptions UI
 */
export function AssetDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [asset, setAsset] = useState<AssetVO | null>(null);

  const load = useCallback(async (assetId: string) => {
    setLoading(true);
    try {
      const vo = await fetchAssetById(assetId);
      setAsset(vo);
    } catch {
      message.error('加载媒资失败');
      setAsset(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (id) {
      void load(id);
    }
  }, [id, load]);

  if (!id) {
    return <Typography.Text type="danger">缺少媒资 ID</Typography.Text>;
  }

  if (loading) {
    return <Spin />;
  }

  if (!asset) {
    return (
      <div>
        <Typography.Text type="secondary">未找到媒资</Typography.Text>
        <Button style={{ marginLeft: 8 }} onClick={() => navigate('/asset/list')}>
          返回列表
        </Button>
      </div>
    );
  }

  return (
    <div>
      <SpaceLikeHeader navigate={navigate} title={asset.title} />
      {asset.previewUrl ? (
        <div style={{ marginBottom: 16 }}>
          <Image src={asset.previewUrl} alt={asset.title} style={{ maxHeight: 240 }} />
        </div>
      ) : null}
      <Descriptions bordered column={1} size="small">
        <Descriptions.Item label="ID">{asset.id}</Descriptions.Item>
        <Descriptions.Item label="标题">{asset.title}</Descriptions.Item>
        <Descriptions.Item label="类型">{asset.type}</Descriptions.Item>
        <Descriptions.Item label="状态">{asset.status}</Descriptions.Item>
        <Descriptions.Item label="编目 ID">{asset.catalogId ?? '—'}</Descriptions.Item>
        <Descriptions.Item label="预览地址">{asset.previewUrl ?? '—'}</Descriptions.Item>
        <Descriptions.Item label="存储路径">{asset.storagePath ?? '—'}</Descriptions.Item>
        <Descriptions.Item label="描述">{asset.description || '—'}</Descriptions.Item>
        <Descriptions.Item label="创建人">{asset.createdBy}</Descriptions.Item>
        <Descriptions.Item label="创建时间">{asset.createdAt}</Descriptions.Item>
        <Descriptions.Item label="更新时间">{asset.updatedAt}</Descriptions.Item>
      </Descriptions>
    </div>
  );
}

function SpaceLikeHeader({ title, navigate }: { title: string; navigate: (path: string) => void }) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
      <Button onClick={() => navigate('/asset/list')}>返回列表</Button>
      <Typography.Title level={4} style={{ margin: 0 }}>
        {title}
      </Typography.Title>
    </div>
  );
}