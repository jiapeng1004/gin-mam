/**
 * LoginPage 组件测试：mock 认证 API，验证表单渲染。
 */

import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import { LoginPage } from '../pages/LoginPage';

vi.mock('../api/authApi', () => ({
  login: vi.fn(),
}));

describe('LoginPage', () => {
  it('renders login form fields', () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    );

    expect(screen.getByLabelText('用户名')).toBeInTheDocument();
    expect(screen.getByLabelText('密码')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /登\s*录/ })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: '媒资管理系统' })).toBeInTheDocument();
  });
});