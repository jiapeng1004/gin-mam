/**
 * TranscodeTaskPage 组件测试：mock 任务分页 API。
 */

import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { TranscodeTaskPage } from '../pages/TranscodeTaskPage';
import * as transcodeApi from '../api/transcodeApi';

vi.mock('../api/transcodeApi', () => ({
  fetchTranscodeTaskPage: vi.fn(),
  fetchTranscodeGroupPage: vi.fn(),
  createTranscodeTask: vi.fn(),
}));

describe('TranscodeTaskPage', () => {
  it('renders task rows from mocked API', async () => {
    vi.mocked(transcodeApi.fetchTranscodeTaskPage).mockResolvedValue({
      list: [
        {
          id: 't1',
          assetId: 'a1',
          transcodeGroupId: 'g1',
          status: 1,
          createdBy: 'u1',
          createdAt: '2026-01-01T00:00:00Z',
          updatedAt: '2026-01-01T00:00:00Z',
        },
      ],
      total: 1,
      page: 1,
      pageSize: 20,
    });

    render(<TranscodeTaskPage />);

    await waitFor(() => {
      expect(screen.getByText('a1')).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: '转码任务' })).toBeInTheDocument();
  });
});