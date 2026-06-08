/**
 * UserListPage 组件测试：mock 分页 API，验证表格渲染用户数据。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { UserListPage } from '../pages/UserListPage';
import * as userApi from '../api/userApi';

vi.mock('../api/userApi', () => ({
  fetchUserPage: vi.fn(),
}));

describe('UserListPage', () => {
  it('renders user rows from mocked page API', async () => {
    vi.mocked(userApi.fetchUserPage).mockResolvedValue({
      list: [
        {
          id: 'u1',
          username: 'admin',
          nickname: '管理员',
          status: 1,
          createdAt: '2026-01-01T00:00:00Z',
          updatedAt: '2026-01-01T00:00:00Z',
        },
      ],
      total: 1,
      page: 1,
      pageSize: 20,
    });

    render(<UserListPage />);

    await waitFor(() => {
      expect(screen.getByText('admin')).toBeInTheDocument();
    });
    expect(screen.getByText('管理员')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: '用户管理' })).toBeInTheDocument();
  });
});