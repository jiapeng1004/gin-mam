/**
 * 媒资回收站说明页：后端暂无独立 recycle 列表 API，说明软删机制。
 */

import { Alert, Empty, Typography } from 'antd';

/**
 * 回收站占位与说明。
 * @returns 说明与空列表 UI
 */
export function AssetRecyclePage() {
  return (
    <div>
      <Typography.Title level={4}>媒资回收站</Typography.Title>
      <Alert
        type="info"
        showIcon
        message="软删除说明"
        description="DELETE /api/v1/asset/:id 会在数据库设置 deleted_at（GORM 软删）。当前后端未提供回收站列表或恢复接口，此处为占位；已删除媒资不会出现在正常分页列表中。"
        style={{ marginBottom: 24 }}
      />
      <Empty description="暂无回收站列表数据（待后端提供 recycle API）" />
    </div>
  );
}