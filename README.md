# GIN-MAM

纯媒资管理系统（MAM），由 MMS 蓝本反推领域模型，**不兼容**旧 MMS API/数据。

## 技术栈

| 层 | 技术 |
|----|------|
| 后端 | Go / Gin / GORM / MySQL / Redis / Elasticsearch / S3 |
| Web | React 18 + Ant Design 5 + Vite |
| 移动端 | Flutter + antd_flutter_mobile（审核端 MVP） |
| 部署 | Alpine Nginx + Go 单容器 |

约束详见 [docs/CONSTRAINTS.md](docs/CONSTRAINTS.md)。

## 目录

```
backend/    Go API
web/        React 管理端
mobile/     Flutter 审核端
deploy/     Dockerfile + Nginx
docs/       设计规格、实现计划、错误码
```

## 本地开发

### 依赖服务

```bash
docker compose up -d mysql redis elasticsearch
```

### 后端

```bash
cd backend
go run ./cmd/server          # :8080，配置见 config.yaml
make test
make swagger               # 重新生成 docs/swagger.*
```

- **Swagger UI**：http://localhost:8080/swagger/index.html
- 修改 Handler 注解后执行 `make swagger` 更新文档

默认账号：`admin` / `admin123`

### Web

```bash
cd web
npm install
npm run dev                  # :5173，/api 代理到 :8080
npm test
```

### 移动端

```bash
cd mobile
flutter pub get
flutter run                  # 模拟器默认 API: http://10.0.2.2:8080/api/v1
```

## 生产部署

### 拉取镜像（阿里云 ACR）

```bash
docker pull registry.cn-hangzhou.aliyuncs.com/jp_aoa/gin-mam:latest
docker run -d -p 8080:80 registry.cn-hangzhou.aliyuncs.com/jp_aoa/gin-mam:latest
```

`master` 分支 push 后 GitHub Actions 自动构建并推送 `:latest`（需配置仓库 Secret `ALIYUN_ACR_PASSWORD`）。

### 本地构建

```bash
cd backend && make vendor && cd ..
docker build -f deploy/Dockerfile -t gin-mam:latest .
docker compose up -d         # 含 app 服务，访问 http://localhost:8080
```

## API 约定

- 前缀：`/api/v1/`
- 鉴权：JWT Bearer（除 `POST /auth/login`、`POST /transcode/callback`）
- 错误码：[docs/err_codes.md](docs/err_codes.md)

## 文档

- [设计规格](docs/superpowers/specs/2026-06-07-gin-mam-design.md)
- [实现计划](docs/superpowers/plans/2026-06-07-gin-mam-implementation.md)
