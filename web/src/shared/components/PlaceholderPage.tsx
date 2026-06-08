/**
 * 通用占位页：M5 脚手架阶段用于菜单路由占位，后续任务替换为真实业务页。
 */

import { Typography } from 'antd';

export interface PlaceholderPageProps {
  /** 页面标题，通常与侧栏菜单文案一致 */
  title: string;
  /** 可选说明文案 */
  description?: string;
}

/**
 * 渲染功能模块占位页。
 * @param props - 标题与说明
 * @returns 占位 UI
 */
export function PlaceholderPage({ title, description }: PlaceholderPageProps) {
  return (
    <div style={{ padding: 24 }}>
      <Typography.Title level={3}>{title}</Typography.Title>
      <Typography.Paragraph type="secondary">
        {description ?? '功能开发中，敬请期待。'}
      </Typography.Paragraph>
    </div>
  );
}