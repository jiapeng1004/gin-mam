#!/bin/sh
# 单容器入口：后台启动 Go API，前台运行 Nginx（:80）
set -e

/app/gin-mam-server -config /app/config.yaml &
GO_PID=$!

cleanup() {
  kill "$GO_PID" 2>/dev/null || true
}
trap cleanup TERM INT

exec nginx -g 'daemon off;'
