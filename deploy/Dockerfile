# Gin MAM 单容器镜像：Alpine Nginx 托管 Web 静态资源 + 反代 Go API
#
# 构建（仓库根目录；Docker 内拉模块不稳定时需先 vendor）：
#   cd backend && make vendor && cd ..
#   docker build -f deploy/Dockerfile -t gin-mam:latest .
#
# 全栈启动：
#   docker compose up -d
#   浏览器访问 http://localhost:8080

# --- Stage 1: Go 后端 ---
FROM golang:1.24-alpine AS go-builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
COPY backend/vendor ./vendor
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-s -w" -o /out/gin-mam-server ./cmd/server

# --- Stage 2: React Web ---
FROM node:20-alpine AS web-builder

WORKDIR /src/web

COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# --- Stage 3: 运行镜像 ---
FROM nginx:1.27-alpine

RUN apk add --no-cache ca-certificates tzdata \
  && rm -f /etc/nginx/conf.d/default.conf

COPY --from=go-builder /out/gin-mam-server /app/gin-mam-server
COPY --from=web-builder /src/web/dist /usr/share/nginx/html
COPY deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY deploy/config.docker.yaml /app/config.yaml
COPY deploy/entrypoint.sh /entrypoint.sh

RUN sed -i 's/\r$//' /entrypoint.sh \
  && chmod +x /entrypoint.sh /app/gin-mam-server

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
