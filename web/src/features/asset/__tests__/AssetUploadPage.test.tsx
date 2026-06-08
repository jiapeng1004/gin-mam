/**
 * AssetUploadPage 测试：验证上传页标题与开始上传按钮。
 */

import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import { AssetUploadPage } from '../pages/AssetUploadPage';

vi.mock('../api/assetApi', () => ({
  initAssetUpload: vi.fn(),
  uploadAssetChunk: vi.fn(),
  completeAssetUpload: vi.fn(),
}));

describe('AssetUploadPage', () => {
  it('renders upload form controls', () => {
    render(
      <MemoryRouter>
        <AssetUploadPage />
      </MemoryRouter>,
    );
    expect(screen.getByRole('heading', { name: '上传媒资' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '开始上传' })).toBeInTheDocument();
  });
});