/**
 * AuditCenterPage 组件测试：mock 待审列表 API，验证 Tab 与表格数据。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { AuditCenterPage } from '../pages/AuditCenterPage';
import * as workflowApi from '../api/workflowApi';

vi.mock('../api/workflowApi', () => ({
  fetchAssignedToMe: vi.fn(),
  fetchCreatedByMe: vi.fn(),
  fetchAuditedByMe: vi.fn(),
  fetchAllInstances: vi.fn(),
  auditWorkflow: vi.fn(),
}));

describe('AuditCenterPage', () => {
  it('renders assigned instances on default tab', async () => {
    vi.mocked(workflowApi.fetchAssignedToMe).mockResolvedValue({
      list: [
        {
          id: 'inst1',
          workflowDefId: 'def1',
          assetId: 'asset1',
          level: 0,
          auditLevel: 2,
          auditStatus: 1,
          createdBy: 'u1',
          createdAt: '2026-01-01T00:00:00Z',
          updatedAt: '2026-01-01T00:00:00Z',
        },
      ],
      total: 1,
      page: 1,
      pageSize: 20,
    });
    vi.mocked(workflowApi.fetchCreatedByMe).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      pageSize: 20,
    });
    vi.mocked(workflowApi.fetchAuditedByMe).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      pageSize: 20,
    });
    vi.mocked(workflowApi.fetchAllInstances).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      pageSize: 20,
    });

    render(<AuditCenterPage />);

    await waitFor(() => {
      expect(screen.getByText('asset1')).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: '审核中心' })).toBeInTheDocument();
  });
});