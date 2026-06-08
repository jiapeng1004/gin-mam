/**
 * 媒资分片上传页：选择文件、进度条、调用 upload API。
 */

import { Button, Form, Input, Progress, Select, Typography, message } from 'antd';
import { useState, type ChangeEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  completeAssetUpload,
  initAssetUpload,
  uploadAssetChunk,
} from '../api/assetApi';

const ASSET_TYPES = [
  { value: 'image', label: '图片' },
  { value: 'video', label: '视频' },
  { value: 'audio', label: '音频' },
  { value: 'document', label: '文档' },
  { value: 'other', label: '其他' },
];

/**
 * 分片上传向导。
 * @returns 上传表单 UI
 */
export function AssetUploadPage() {
  const navigate = useNavigate();
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [percent, setPercent] = useState(0);
  const [form] = Form.useForm<{ title: string; type: string; catalogId?: string }>();

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const picked = e.target.files?.[0] ?? null;
    setFile(picked);
    if (picked && !form.getFieldValue('title')) {
      form.setFieldsValue({ title: picked.name });
    }
  };

  const handleUpload = async () => {
    if (!file) {
      message.warning('请先选择文件');
      return;
    }
    const values = await form.validateFields();
    setUploading(true);
    setPercent(0);
    try {
      const init = await initAssetUpload({
        fileName: file.name,
        fileSize: file.size,
        mimeType: file.type || undefined,
      });
      const chunkSize = init.chunkSize;
      const totalChunks = Math.ceil(file.size / chunkSize);
      for (let index = 0; index < totalChunks; index += 1) {
        const start = index * chunkSize;
        const end = Math.min(start + chunkSize, file.size);
        const blob = file.slice(start, end);
        await uploadAssetChunk(init.uploadId, index, blob, file.name);
        setPercent(Math.round(((index + 1) / totalChunks) * 90));
      }
      const asset = await completeAssetUpload({
        uploadId: init.uploadId,
        assetTitle: values.title,
        type: values.type,
        catalogId: values.catalogId?.trim() || undefined,
      });
      setPercent(100);
      message.success('上传成功');
      navigate(`/asset/${asset.id}`);
    } catch {
      message.error('上传失败');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <Typography.Title level={4}>上传媒资</Typography.Title>
      <Form form={form} layout="vertical" initialValues={{ type: 'video' }} style={{ maxWidth: 480 }}>
        <Form.Item label="选择文件">
          <input type="file" onChange={onFileChange} disabled={uploading} />
        </Form.Item>
        <Form.Item name="title" label="媒资标题" rules={[{ required: true, message: '\u8bf7\u8f93\u5165\u6807\u9898' }]}>
          <Input placeholder="媒资标题" />
        </Form.Item>
        <Form.Item name="type" label="媒资类型" rules={[{ required: true }]}>
          <Select options={ASSET_TYPES} />
        </Form.Item>
        <Form.Item name="catalogId" label="编目 ID（可选）">
          <Input placeholder="catalogId" />
        </Form.Item>
        {uploading || percent > 0 ? (
          <Progress percent={percent} status={uploading ? 'active' : undefined} />
        ) : null}
        <Button type="primary" onClick={() => void handleUpload()} loading={uploading} disabled={!file}>
          开始上传
        </Button>
      </Form>
    </div>
  );
}
