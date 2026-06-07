# GIN-MAM 设计规格

> MMS 蓝本重构 · 纯媒资管理（MAM）  
> 日期：2026-06-07  
> 状态：已批准（2026-06-07）

## 1. 项目概述

### 1.1 目标

将 `C:\Users\24291\IdeaProjects\mms`（Spring Boot + Vue3 + UniApp）重构为 **GIN-MAM**：

| 层级 | 目标技术栈 |
|------|-----------|
| 后端 | Gin + Viper + GORM + Zap + MySQL + Elasticsearch + Redis |
| Web | React + Ant Design + Vite |
| 移动端 | Flutter（替代 UniApp） |

- **不兼容**旧 MMS API 与数据，以 MMS 为**蓝本**反推领域模型
- 剥离选题、舆情、宣传指令等非核心功能
- AI 处理保留入口，MVP 显示「开发中」

### 1.2 首期里程碑（D）

Web 全栈闭环 + Flutter 最小审核端：

- 系统管理（用户/角色/组织/菜单/配置）
- 媒资核心（上传/编目/检索/预览/回收/分享/版本）
- 转码任务
- 媒资多级审核工作流
- Flutter：登录 + 待审列表 + 审核操作

### 1.3 Monorepo 布局

```
gin-mam/
├── backend/          # Go
├── web/              # React
├── mobile/           # Flutter
├── docs/
│   ├── CONSTRAINTS.md
│   └── superpowers/specs/
├── deploy/           # Dockerfile, nginx, config 模板
└── docker-compose.yml
```

---

## 2. 项目约束（CONSTRAINTS 摘要）

完整约束写入 `docs/CONSTRAINTS.md`，核心条目如下。

### 2.1 ID 约束

- 所有业务表主键：`varchar(32)`，值为 **UUID v7 去横杠**（32 位 hex）
- **禁止**：雪花 ID、自增 ID、带横杠 UUID、`contentSourceId` 等第二标识体系
- 外键命名：`{entity}_id`，类型一律 `varchar(32)`
- 对外标识与内部主键合一

### 2.2 表命名

- 前缀：`gm_`，**单数**，禁止复数（`gm_asset` ✅，`gm_assets` ❌）
- **禁止** GORM 全局 `TablePrefix`
- 每个实体 `TableName()` 返回完整表名

### 2.3 租户

- 所有业务表含 `tenant_id varchar(32)`
- MVP 固定 `"default"`，不做隔离逻辑
- 二期启用行级租户过滤

### 2.4 配置分层

| 层级 | 来源 | 内容 |
|------|------|------|
| L0 | `config.yaml`（Viper） | MySQL、Redis、ES 连接；HTTP 端口；日志；JWT Secret |
| L1 | `gm_sys_config` 表 | S3、预览链、转码 URL、分享前缀、功能开关等**一切业务参数** |
| L2 | `gm_catalog_config` 表 | 栏目级编目与转码组绑定 |

- 读取统一经 `ConfigService`，Redis 缓存 `config:{tenant_id}:{key}`，变更主动失效
- Key 命名：`MAM_{CATEGORY}_{NAME}` 全大写

### 2.5 依赖注入

- **禁止** wire / dig
- 统一 `internal/app/graph.go` 手动组装
- Repository / 外部客户端 / TxManager / Service 用 **interface**，mockgen 生成 mock

### 2.6 Mock

- 统一 **mockgen**（`go.uber.org/mock/mockgen`）
- `//go:generate mockgen -source=xxx.go -destination=mock/xxx_mock.go -package=mock`
- Makefile：`make generate-mocks`

### 2.7 测试

- 前后端均必须有单元测试，CI 强制跑
- 后端：`testing` + gomock；核心 domain 覆盖率 ≥ 70%
- 前端：Vitest + Testing Library
- Flutter：mocktail / mockito

### 2.8 事务

- 禁止 Service 直接 `db.Transaction()`
- 统一 `infra/tx.Manager.Run(ctx, propagation, fn)`
- 传播类型：`Required`（默认）/ `RequiresNew` / `Mandatory` / `NotSupported` / `Never`
- Repository 经 `TxFromContext(ctx)` 取当前 tx

### 2.9 HTTP 响应

**成功：无外层包装，REST 风格**

- HTTP 2xx + 业务 JSON，字段 **camelCase**
- 禁止 `ResultDTO` / `CommonResult` / `{ code, data, message }`

**失败：snake_case 错误字段名，与成功体隔离**

- HTTP 4xx/5xx + 错误 JSON
- 错误体字段名固定为 **`err_code`**、**`err_msg`**（蛇形命名）
- 成功体用 camelCase（如 `previewUrl`），错误体用 snake_case 键名，**一眼可区分**
- `err_code` 值为**数值型**业务错误码（如 `40001`），不是 snake_case 字符串
- Handler 统一 `httpx.Fail` / `BizError`，禁止手写错误 JSON

```json
// 成功 200
{ "id": "...", "previewUrl": "...", "createdAt": "..." }

// 失败 404
{ "err_code": 40401, "err_msg": "媒资不存在" }
```

---

## 3. 架构

### 3.1 推荐方案：Flat Domain Monolith

```
backend/
  cmd/server/main.go
  internal/
    app/graph.go
    pkg/id/              # UUID v7
    pkg/httpx/           # OK / Fail / BizError / ErrorVo
    domain/
      sys/ asset/ catalog/ transcode/ workflow/ search/
    infra/
      mysql/ redis/ elasticsearch/ storage/ config/ tx/
```

转码/ES 消费量大时，可将 consumer 抽到 `cmd/worker/`，无需推翻架构。

### 3.2 基础设施

| 组件 | 用途 |
|------|------|
| MySQL | 主数据 |
| Redis | 配置缓存、分片上传、工作流锁、JWT 黑名单 |
| Elasticsearch | 媒资检索 |
| S3 兼容存储 | 文件存储（配置表驱动） |
| Nginx Alpine | 与 Go 二进制同容器部署 |

### 3.3 部署

单容器：Alpine Nginx 静态资源 + 反代 `/api/` → Go `:8080`

```
/app/gin-mam-server
/usr/share/nginx/html/     # web/dist
/etc/nginx/conf.d/
```

Flutter 独立构建 APK/IPA，API 指向同一 `/api/` 入口。

---

## 4. 数据模型

### 4.1 核心表

```
gm_user / gm_role / gm_org / gm_menu / gm_sys_config
gm_asset / gm_asset_file / gm_asset_metadata
gm_catalog / gm_catalog_config
gm_transcode_task / gm_transcode_group
gm_workflow_def / gm_workflow_level_user
gm_workflow_instance / gm_workflow_operate
gm_workflow_instance_level_user
```

### 4.2 媒资简化（对比 MMS）

| MMS 问题 | GIN-MAM |
|----------|---------|
| 雪花 id + contentSourceId 双 ID | 仅 `gm_asset.id`（UUID v7 hex） |
| 预览 URL 散落 | DB 存 `storage_path`，响应时 ConfigService 拼链 |
| 多表职责模糊 | asset / asset_file / asset_metadata 分离 |

### 4.3 预览 URL 拼装

```
图片预览 = Get("MAM_IMAGE_ACCESS_DOMAIN") + asset_file.storage_path
转码播放 = Get("MAM_VOD_ACCESS_DOMAIN")    + transcode.output_path
Office   = Get("MAM_OFFICE_PREVIEW_URL")     + encode(source_url)
```

### 4.4 MVP 配置 Key（L1）

**STORAGE：** `MAM_STORAGE_S3_ENDPOINT` / `_ACCESS_KEY` / `_SECRET_KEY` / `_BUCKET` / `_REGION` / `_PATH_STYLE`

**PREVIEW：** `MAM_IMAGE_ACCESS_DOMAIN` / `MAM_VIDEO_ACCESS_DOMAIN` / `MAM_VOD_ACCESS_DOMAIN` / `MAM_OTHER_ACCESS_DOMAIN` / `MAM_IMAGE_CDN_DOMAIN` / `MAM_OFFICE_PREVIEW_PLAT` / `MAM_OFFICE_PREVIEW_URL` / `MAM_OFFICE_PREVIEW_MAX_SIZE`

**TRANSCODE：** `MAM_TRANSCODE_URL` / `MAM_TRANSCODE_CALLBACK_URL`

**SHARE：** `MAM_PC_SHARE_URL_PREFIX` / `MAM_H5_SHARE_URL_PREFIX` / `MAM_SHORT_URL_PREFIX`

**SYSTEM：** `MAM_DOMAIN_URL` / `MAM_UPLOAD_CHUNK_SIZE`

---

## 5. API 设计

### 5.1 通用约定

- 前缀：`/api/v1/`
- 鉴权：JWT Bearer
- 分页成功体：`{ "list": [], "total": 0, "page": 1, "pageSize": 20 }`（camelCase，无包装）

### 5.2 路由

```
/api/v1/
├── auth/
├── sys/user|role|org|menu|config/
├── asset/ + upload/ + recycle/ + share/
├── catalog/
├── transcode/ + callback
├── search/asset
├── asset-workflow/
└── ai/*  → 501 开发中
```

### 5.3 媒资工作流 API（对标 MMS `/biz/resources/workflow/*`）

| Method | Path | 说明 |
|--------|------|------|
| POST | `/asset-workflow/submit` | 送审 |
| POST | `/asset-workflow/audit` | 单独审核 |
| POST | `/asset-workflow/multi-audit` | 批量审核 |
| POST | `/asset-workflow/revoke` | 撤回 |
| POST | `/asset-workflow/assigned-to-me` | 待我审核（Flutter 核心） |
| POST | `/asset-workflow/created-by-me` | 我发起的 |
| POST | `/asset-workflow/audited-by-me` | 我审过的 |
| POST | `/asset-workflow/all` | 全部记录 |

### 5.4 错误码规划（数值 err_code）

| err_code | HTTP | 含义 |
|----------|------|------|
| 40100 | 401 | 未登录 |
| 40300 | 403 | 无权限 |
| 40401 | 404 | 媒资不存在 |
| 40901 | 409 | 工作流锁定冲突 |
| 42201 | 422 | 参数校验失败 |
| 50000 | 500 | 内部错误 |

域前缀段：40xxx 客户端，50xxx 服务端；具体码表实现阶段在 `docs/err_codes.md` 维护。

### 5.5 ErrorVo 结构

```go
type ErrorVo struct {
    ErrCode int    `json:"err_code"`
    ErrMsg  string `json:"err_msg"`
}
// 校验失败可选
type ValidationErrorVo struct {
    ErrCode int               `json:"err_code"`
    ErrMsg  string            `json:"err_msg"`
    Details []FieldError      `json:"details,omitempty"`
}
```

---

## 6. 工作流

### 6.1 引擎

对标 MMS `mms-plugin-workflow` 多级审核 + `BizResourcesWorkFLowController` 媒资封装。  
**不是** Snowy BPMN，**不是** `biz/modular/workflow` 图引擎。

### 6.2 审核状态

| status | 名称 |
|--------|------|
| 0 | 待审核 |
| 1 | 审核中 |
| 2 | 通过 |
| 3 | 打回 |

流转：送审 → 逐级审核 → 通过/打回；可撤回；打回后可重新送审。

### 6.3 并发

Redis 锁 `wf:lock:{instance_id}` + `locker_id` 字段。

---

## 7. 前端

### 7.1 Web（React + Ant Design + Vite）

```
web/src/features/
  auth/ dashboard/
  sys/{user,role,org,menu,config}/
  asset/{list,detail,upload,recycle,share}/
  catalog/ transcode/
  workflow/{def,audit}/
  ai/   # 开发中占位
```

axios：`2xx → response.data` 即业务体；`4xx/5xx → err_code / err_msg`

### 7.2 Flutter MVP

```
mobile/lib/features/
  auth/login
  audit/{pending_list, detail, action}
```

仅登录 + 待审 + 审核；上传/编目留 Web。

---

## 8. 事务模板

```go
type Manager interface {
    Run(ctx context.Context, prop Propagation, fn func(ctx context.Context) error) error
}
```

| 传播 | 行为 |
|------|------|
| Required | 有则加入，无则新建（默认） |
| RequiresNew | 新建独立事务 |
| Mandatory | 必须在事务内 |
| NotSupported | 挂起事务，非事务执行 |
| Never | 禁止在事务内 |

---

## 9. MVP 分期

| 阶段 | 内容 |
|------|------|
| M1 | backend 骨架 + sys + JWT + graph.go + tx + httpx |
| M2 | asset + catalog + upload + ES |
| M3 | workflow 送审/审核 |
| M4 | transcode |
| M5 | Web 全页面 |
| M6 | Flutter 审核端 |

---

## 10. 明确不做（MVP）

- 选题 / 舆情 / 宣传指令 / 知识库 / BPMN 通用流程
- 旧 MMS 数据迁移与 API 兼容
- AI 处理实功能（仅占位）
- 多租户运行时隔离（仅预留字段）
- wire 依赖注入
- ResultDTO 外层包装

---

## 11. 参考蓝本（MMS 路径）

| 能力 | MMS 参考 |
|------|----------|
| 媒资 | `mms-plugin-biz/modular/resources` |
| 编目 | `mms-plugin-biz/modular/catalog` |
| 转码配置 | `mms-plugin-biz/modular/resourcesconfig` |
| 工作流 | `mms-plugin-workflow` + `BizResourcesWorkFLowController` |
| 系统管理 | `mms-plugin-sys` |
| 配置表 | `mms-plugin-dev/modular/config`（DevConfig + Redis） |
| 存储 | `DevFileMinIoUtil` + ConfigService |
| 部署 | 根目录 `Dockerfile`（Nginx + JAR） |
