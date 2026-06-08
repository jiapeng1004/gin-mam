/**
 * AI 能力占位页：展示「开发中」状态，预留后续智能能力入口。
 */

import { Result } from 'antd';

/**
 * AI 模块占位页，使用 Ant Design Result 提示开发中。
 * @returns 占位 UI
 */
export function AiPlaceholderPage() {
  return (
    <Result
      status="info"
      title="开发中"
      subTitle="智能编目、语义检索等 AI 能力正在规划中，敬请期待。"
    />
  );
}