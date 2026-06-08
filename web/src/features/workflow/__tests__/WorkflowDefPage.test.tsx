/**
 * WorkflowDefPage 组件测试：mock 流程定义分页 API，验证列表渲染。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { WorkflowDefPage } from '../pages/WorkflowDefPage';
import * as workflowApi from '../api/workflowApi';

vi.mock('../api/workflowApi', () => ({
  fetchWorkflowDefPage: vi.fn(),
  createWorkflowDef: vi.fn(),
}));

describe('WorkflowDefPage', () => {
  it('renders workflow def rows from mocked API', async () => {
    vi.mocked(workflowApi.fetchWorkflowDefPage).mockResolvedValue({
      list: [
        {
          id: 'def1',
          name: '标准审核',
          auditLevel: 2,
          status: 1,
          createdAt: '2026-01-01T00:00:00Z',
          updatedAt: '2026-01-01T00:00:00Z',
        },
      ],
      total: 1,
      page: 1,
      pageSize: 20,
    });

    render(<WorkflowDefPage />);

    await waitFor(() => {
      expect(screen.getByText('标准审核')).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: '流程定义' })).toBeInTheDocument();
  });
});