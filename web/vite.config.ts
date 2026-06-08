/**
 * Vite 构建配置：React 插件、开发服务器端口与后端 API 代理。
 */

import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      // 浏览器请求 /api/* 时转发到 Go 后端，避免开发期 CORS
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
});