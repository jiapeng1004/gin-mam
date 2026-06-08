/**
 * AssetRecyclePage 测试：验证软删说明文案。
 */

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { AssetRecyclePage } from '../pages/AssetRecyclePage';

describe('AssetRecyclePage', () => {
  it('shows soft delete explanation', () => {
    render(<AssetRecyclePage />);
    expect(screen.getByRole('heading', { name: '媒资回收站' })).toBeInTheDocument();
    expect(screen.getByText('软删除说明')).toBeInTheDocument();
  });
});