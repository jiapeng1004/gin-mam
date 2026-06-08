/**
 * CatalogTreePage 测试：mock 树 API，验证节点名称渲染。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CatalogTreePage } from '../pages/CatalogTreePage';
import * as catalogApi from '../api/catalogApi';

vi.mock('../api/catalogApi', () => ({
  fetchCatalogTree: vi.fn(),
  fetchCatalogConfig: vi.fn(),
  createCatalog: vi.fn(),
  deleteCatalog: vi.fn(),
  putCatalogConfig: vi.fn(),
}));

describe('CatalogTreePage', () => {
  it('renders catalog tree nodes', async () => {
    vi.mocked(catalogApi.fetchCatalogTree).mockResolvedValue([
      { id: 'c1', name: '根编目', children: [{ id: 'c2', name: '子编目', children: [] }] },
    ]);

    render(<CatalogTreePage />);

    await waitFor(() => {
      expect(screen.getByText('根编目')).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: '编目管理' })).toBeInTheDocument();
  });
});