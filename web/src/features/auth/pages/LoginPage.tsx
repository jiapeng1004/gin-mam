/**
 * 登录页：未挂载 MainLayout，提交后调用认证 API 并写入 token。
 */

import { Button, Card, Form, Input, Typography, message } from 'antd';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api/authApi';
import { setToken } from '../../../shared/api/authStorage';
import { ApiError } from '../../../shared/api/errors';

interface LoginFormValues {
  username: string;
  password: string;
}

/**
 * 登录表单页；成功后将 JWT 存入 localStorage 并跳转媒资列表。
 * @returns 登录 UI
 */
export function LoginPage() {
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);

  const onFinish = async (values: LoginFormValues) => {
    setSubmitting(true);
    try {
      const result = await login(values);
      setToken(result.token);
      navigate('/asset/list', { replace: true });
    } catch (err) {
      const msg = err instanceof ApiError ? err.errMsg : '登录失败，请稍后重试';
      message.error(msg || '登录失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: '#f0f2f5',
      }}
    >
      <Card style={{ width: 400 }}>
        <Typography.Title level={4} style={{ textAlign: 'center', marginBottom: 24 }}>
          媒资管理系统
        </Typography.Title>
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item
            label="用户名"
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="admin" autoComplete="username" />
          </Form.Item>
          <Form.Item
            label="密码"
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="******" autoComplete="current-password" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={submitting}>
            登录
          </Button>
        </Form>
      </Card>
    </div>
  );
}