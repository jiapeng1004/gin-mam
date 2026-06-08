/**
 * AssetListPage 测试：mock 分页 API，验证表格渲染。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import { AssetListPage } from '../pages/AssetListPage';
import * as assetApi from '../api/assetApi';

vi.mock('../api/assetApi', () => ({
  fetchAssetPage: vi.fn(),
}));

describe('AssetListPage', () => {
  it('renders asset rows from mocked page API', async () => {
    vi.mocked(assetApi.fetchAssetPage).mockResolvedValue({
      list: [
        {
          id: 'a1',
          title: '测试视频',
          type: 'video',
          status: 0,
          createdBy: 'admin',
          createdAt: '2026-01-01T00:00:00Z',
          updatedAt: '2026-01-01T00:00:00Z',
        },
      ],
      total: 1,
      page: 1,
      pageSize: 20,
    });

    render(
      <MemoryRouter>
        <AssetListPage />
      </MemoryRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText('测试视频')).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: '媒资列表' })).toBeInTheDocument();
  });
});