/**
 * 主布局：登录后页面的 Ant Design Layout（侧栏 + 顶栏 + 内容区）。
 */

import {
  AuditOutlined,
  CloudUploadOutlined,
  DatabaseOutlined,
  RobotOutlined,
  SettingOutlined,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Layout, Menu, theme, Typography } from 'antd';
import type { MenuProps } from 'antd';
import { useMemo, type ReactNode } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

const { Header, Sider, Content } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

/** 侧栏菜单项与路由 path 的对应关系 */
const MENU_ITEMS: { key: string; label: string; icon: ReactNode }[] = [
  { key: '/asset/list', label: '媒资管理', icon: <VideoCameraOutlined /> },
  { key: '/catalog', label: '编目', icon: <DatabaseOutlined /> },
  { key: '/review', label: '审核中心', icon: <AuditOutlined /> },
  { key: '/transcode', label: '转码', icon: <CloudUploadOutlined /> },
  { key: '/system', label: '系统管理', icon: <SettingOutlined /> },
  { key: '/ai', label: 'AI（占位）', icon: <RobotOutlined /> },
];

/**
 * 应用主框架：左侧导航 + 顶栏标题 + 子路由出口。
 * @returns MainLayout 结构
 */
export function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const items: MenuItem[] = useMemo(
    () =>
      MENU_ITEMS.map((item) => ({
        key: item.key,
        icon: item.icon,
        label: item.label,
      })),
    [],
  );

  const selectedKey = useMemo(() => {
    const match = MENU_ITEMS.find((item) => location.pathname.startsWith(item.key));
    return match?.key ?? '/asset/list';
  }, [location.pathname]);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg" collapsedWidth={0}>
        <div
          style={{
            height: 48,
            margin: 16,
            color: '#fff',
            fontWeight: 600,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          Gin-MAM
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={items}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ padding: '0 24px', background: colorBgContainer }}>
          <Typography.Text strong>媒资管理系统</Typography.Text>
        </Header>
        <Content style={{ margin: '16px' }}>
          <div
            style={{
              padding: 24,
              minHeight: 360,
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
            }}
          >
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  );
}