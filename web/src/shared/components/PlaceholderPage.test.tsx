/**
 * PlaceholderPage 组件 smoke 测试。
 */

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { PlaceholderPage } from './PlaceholderPage';

describe('PlaceholderPage', () => {
  it('renders title and default description', () => {
    render(<PlaceholderPage title="测试模块" />);
    expect(screen.getByRole('heading', { name: '测试模块' })).toBeInTheDocument();
    expect(screen.getByText('功能开发中，敬请期待。')).toBeInTheDocument();
  });
});