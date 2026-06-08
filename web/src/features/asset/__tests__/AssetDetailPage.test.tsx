/**
 * AssetDetailPage 测试：mock 详情 API，验证 previewUrl 展示。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import { AssetDetailPage } from '../pages/AssetDetailPage';
import * as assetApi from '../api/assetApi';

vi.mock('../api/assetApi', () => ({
  fetchAssetById: vi.fn(),
}));

describe('AssetDetailPage', () => {
  it('shows asset detail including previewUrl', async () => {
    vi.mocked(assetApi.fetchAssetById).mockResolvedValue({
      id: 'a1',
      title: '封面图',
      type: 'image',
      status: 1,
      previewUrl: 'https://cdn.example/preview.jpg',
      createdBy: 'u1',
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
    });

    render(
      <MemoryRouter initialEntries={['/asset/a1']}>
        <Routes>
          <Route path="/asset/:id" element={<AssetDetailPage />} />
        </Routes>
      </MemoryRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText('https://cdn.example/preview.jpg')).toBeInTheDocument();
    });
    expect(screen.getAllByText('封面图').length).toBeGreaterThan(0);
  });
});
