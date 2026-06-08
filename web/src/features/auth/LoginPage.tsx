/**
 * 登录页：未挂载 MainLayout，后续任务接入真实认证表单与接口。
 */

import { Button, Card, Form, Input, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import { AUTH_TOKEN_KEY } from '../../app/auth';

/**
 * 登录占位页；提交后写入本地 token 并跳转媒资列表。
 * @returns 登录表单 UI
 */
export function LoginPage() {
  const navigate = useNavigate();

  const onFinish = () => {
    // 脚手架阶段：模拟登录成功，任务 20 将替换为真实 JWT/会话
    localStorage.setItem(AUTH_TOKEN_KEY, 'dev-placeholder-token');
    navigate('/asset/list', { replace: true });
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
          <Form.Item label="用户名" name="username" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input placeholder="admin" />
          </Form.Item>
          <Form.Item label="密码" name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password placeholder="******" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            登录
          </Button>
        </Form>
      </Card>
    </div>
  );
}