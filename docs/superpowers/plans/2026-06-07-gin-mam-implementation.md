# GIN-MAM 实现计划

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 以 MMS 为蓝本，在 `gin-mam` monorepo 中实现纯 MAM 系统（Go 后端 + React Web + Flutter 审核端），首期里程碑 D。

**架构：** Flat Domain Monolith；`graph.go` 手动 DI；`gm_` 单数表 + UUID v7 主键；L1 配置走 `gm_sys_config` + Redis；REST 无包装响应；`infra/tx` 五传播事务。

**技术栈：** Go 1.23 / Gin / GORM / Viper / Zap / mockgen / MySQL / Redis / ES / S3 · React 18 / Ant Design 5 / Vite / Vitest · Flutter 3 / mocktail

**规格：** [`docs/superpowers/specs/2026-06-07-gin-mam-design.md`](../specs/2026-06-07-gin-mam-design.md)  
**约束：** [`docs/CONSTRAINTS.md`](../../CONSTRAINTS.md)

---

## 文件结构总览

### 后端（`backend/`）

| 路径 | 职责 |
|------|------|
| `cmd/server/main.go` | 入口：读 Viper → `app.Build()` → 监听 |
| `internal/app/graph.go` | 手动 DI 唯一入口 |
| `internal/app/router.go` | Gin 路由注册 |
| `internal/pkg/id/id.go` | UUID v7 生成（32 hex） |
| `internal/pkg/httpx/response.go` | OK/Created/Fail/ErrorVo/BizError |
| `internal/pkg/httpx/middleware.go` | ErrorMiddleware、JWT 中间件 |
| `internal/infra/tx/` | 五传播事务 Manager |
| `internal/infra/mysql/mysql.go` | GORM 连接 |
| `internal/infra/redis/redis.go` | go-redis 客户端 |
| `internal/infra/config/` | ConfigService（DB+Redis） |
| `internal/infra/storage/` | S3 Client interface + 实现 |
| `internal/infra/elasticsearch/` | ES Client |
| `internal/domain/sys/` | 用户/角色/组织/菜单/配置 |
| `internal/domain/asset/` | 媒资/上传/回收/分享 |
| `internal/domain/catalog/` | 编目 |
| `internal/domain/transcode/` | 转码 |
| `internal/domain/workflow/` | 多级审核 |
| `internal/domain/search/` | ES 检索 |
| `migrations/` | SQL 迁移 |
| `config.yaml` | L0 进程配置模板 |
| `Makefile` | test / generate-mocks / build |

### Web（`web/`）

| 路径 | 职责 |
|------|------|
| `src/app/` | 路由、Layout、Provider |
| `src/shared/api/client.ts` | axios + err_code/err_msg 拦截 |
| `src/features/*/` | 按域划分页面与测试 |

### Mobile（`mobile/`）

| 路径 | 职责 |
|------|------|
| `lib/features/auth/` | 登录 |
| `lib/features/audit/` | 待审/审核 |
| `lib/shared/api/` | Dio/http + ErrorVo 解析 |

### 部署（`deploy/` + 根目录）

| 路径 | 职责 |
|------|------|
| `docker-compose.yml` | MySQL + Redis + ES（+ MinIO 可选） |
| `deploy/Dockerfile` | Alpine Nginx + Go 二进制 |
| `deploy/nginx/default.conf` | SPA + `/api/` 反代 |

---

# 里程碑 M1：后端骨架 + 系统管理 + 基础设施

## 任务 1：初始化 Go 模块与目录

**文件：**
- 创建：`backend/go.mod`
- 创建：`backend/Makefile`
- 创建：`backend/config.yaml`
- 创建：`backend/.gitignore`

- [ ] **步骤 1：创建 go.mod**

```bash
cd backend
go mod init github.com/gin-mam/backend
```

- [ ] **步骤 2：安装核心依赖**

```bash
go get github.com/gin-gonic/gin@v1.10.0
go get gorm.io/gorm@v1.25.12
go get gorm.io/driver/mysql@v1.5.7
go get github.com/spf13/viper@v1.19.0
go get go.uber.org/zap@v1.27.0
go get github.com/redis/go-redis/v9@v9.7.0
go get github.com/golang-jwt/jwt/v5@v5.2.1
go get github.com/google/uuid@v1.6.0
go get go.uber.org/mock/mockgen@v0.5.0
go get github.com/stretchr/testify@v1.9.0
```

- [ ] **步骤 3：创建 config.yaml 模板**

```yaml
# backend/config.yaml
server:
  port: 8080
  mode: debug
mysql:
  dsn: "root:root@tcp(127.0.0.1:3306)/gin_mam?charset=utf8mb4&parseTime=True&loc=UTC"
redis:
  addr: "127.0.0.1:6379"
  db: 0
elasticsearch:
  addresses: ["http://127.0.0.1:9200"]
jwt:
  secret: "change-me-in-production"
  expire_hours: 24
log:
  level: debug
```

- [ ] **步骤 4：创建 Makefile**

```makefile
.PHONY: test generate-mocks build
generate-mocks:
	go generate ./...
test:
	go test ./... -race -count=1
build:
	go build -o bin/gin-mam-server ./cmd/server
```

- [ ] **步骤 5：Commit**

```bash
git add backend/
git commit -m "chore(backend): init go module and project skeleton"
```

---

## 任务 2：UUID v7 ID 生成器

**文件：**
- 创建：`backend/internal/pkg/id/id.go`
- 创建：`backend/internal/pkg/id/id_test.go`

- [ ] **步骤 1：编写失败测试**

```go
// backend/internal/pkg/id/id_test.go
package id_test

import (
	"testing"
	"github.com/gin-mam/backend/internal/pkg/id"
	"github.com/stretchr/testify/require"
)

func TestNew_Returns32HexChars(t *testing.T) {
	got := id.New()
	require.Len(t, got, 32)
	require.NotContains(t, got, "-")
	require.Regexp(t, `^[0-9a-f]{32}$`, got)
}

func TestNew_IsUnique(t *testing.T) {
	a, b := id.New(), id.New()
	require.NotEqual(t, a, b)
}
```

- [ ] **步骤 2：运行测试验证失败**

```bash
cd backend && go test ./internal/pkg/id/... -v
```

预期：FAIL，`undefined: id.New`

- [ ] **步骤 3：实现 id.go**

```go
// backend/internal/pkg/id/id.go
package id

import (
	"strings"
	"github.com/google/uuid"
)

func New() string {
	return strings.ReplaceAll(uuid.Must(uuid.NewV7()).String(), "-", "")
}
```

- [ ] **步骤 4：运行测试验证通过**

```bash
cd backend && go test ./internal/pkg/id/... -v
```

预期：PASS

- [ ] **步骤 5：Commit**

```bash
git add backend/internal/pkg/id/
git commit -m "feat(backend): add UUID v7 hex ID generator"
```

---

## 任务 3：HTTP 响应层（无包装 REST）

**文件：**
- 创建：`backend/internal/pkg/httpx/response.go`
- 创建：`backend/internal/pkg/httpx/error.go`
- 创建：`backend/internal/pkg/httpx/middleware.go`
- 创建：`backend/internal/pkg/httpx/response_test.go`

- [ ] **步骤 1：编写失败测试**

```go
// backend/internal/pkg/httpx/response_test.go
package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
)

func TestFail_ReturnsSnakeCaseErrorVo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	httpx.Fail(c, http.StatusNotFound, 40401, "媒资不存在")
	require.Equal(t, 404, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, float64(40401), body["err_code"])
	require.Equal(t, "媒资不存在", body["err_msg"])
	require.NotContains(t, body, "errCode")
}
```

- [ ] **步骤 2：运行测试验证失败**

```bash
cd backend && go test ./internal/pkg/httpx/... -v
```

- [ ] **步骤 3：实现 response.go 与 error.go**

```go
// backend/internal/pkg/httpx/error.go
package httpx

type ErrorVo struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type BizError struct {
	HTTPStatus int
	ErrCode    int
	ErrMsg     string
}

func (e *BizError) Error() string { return e.ErrMsg }

func NewBizError(httpStatus, errCode int, msg string) *BizError {
	return &BizError{HTTPStatus: httpStatus, ErrCode: errCode, ErrMsg: msg}
}
```

```go
// backend/internal/pkg/httpx/response.go
package httpx

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data any)       { c.JSON(http.StatusOK, data) }
func Created(c *gin.Context, data any)   { c.JSON(http.StatusCreated, data) }
func NoContent(c *gin.Context)           { c.Status(http.StatusNoContent) }

func Fail(c *gin.Context, httpStatus, errCode int, errMsg string) {
	c.JSON(httpStatus, ErrorVo{ErrCode: errCode, ErrMsg: errMsg})
}
```

```go
// backend/internal/pkg/httpx/middleware.go
package httpx

import "github.com/gin-gonic/gin"

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}
		err := c.Errors.Last().Err
		if biz, ok := err.(*BizError); ok {
			Fail(c, biz.HTTPStatus, biz.ErrCode, biz.ErrMsg)
			return
		}
		Fail(c, 500, 50000, "服务器内部错误")
	}
}
```

- [ ] **步骤 4：运行测试验证通过**

```bash
cd backend && go test ./internal/pkg/httpx/... -v
```

- [ ] **步骤 5：Commit**

```bash
git add backend/internal/pkg/httpx/
git commit -m "feat(backend): add REST response layer without ResultDTO wrapper"
```

---

## 任务 4：GORM 事务 Manager（五传播）

**文件：**
- 创建：`backend/internal/infra/tx/propagation.go`
- 创建：`backend/internal/infra/tx/context.go`
- 创建：`backend/internal/infra/tx/manager.go`
- 创建：`backend/internal/infra/tx/errors.go`
- 创建：`backend/internal/infra/tx/manager_test.go`

- [ ] **步骤 1：编写 Mandatory 失败测试**

```go
// backend/internal/infra/tx/manager_test.go 片段
func TestMandatory_WithoutTx_ReturnsError(t *testing.T) {
	db := openTestDB(t)
	m := tx.NewManager(db)
	err := m.Run(context.Background(), tx.Mandatory, func(ctx context.Context) error {
		return nil
	})
	require.ErrorIs(t, err, tx.ErrTxMandatory)
}
```

- [ ] **步骤 2：实现 propagation / context / errors / manager**

`Propagation` 枚举：`Required`, `RequiresNew`, `Mandatory`, `NotSupported`, `Never`

`Manager` interface + `NewManager(db *gorm.DB) Manager`

`TxFromContext` / `ContextWithTx` / `ContextWithoutTx`

- [ ] **步骤 3：编写 Required / RequiresNew 集成测试**

使用 SQLite 内存库或 testcontainers MySQL 验证：
- Required 嵌套共提交
- RequiresNew 内层回滚不影响外层

- [ ] **步骤 4：运行测试**

```bash
cd backend && go test ./internal/infra/tx/... -v
```

- [ ] **步骤 5：Commit**

```bash
git add backend/internal/infra/tx/
git commit -m "feat(backend): add GORM transaction manager with 5 propagations"
```

---

## 任务 5：MySQL 迁移 — 系统表

**文件：**
- 创建：`backend/migrations/001_sys.up.sql`
- 创建：`backend/migrations/001_sys.down.sql`

- [ ] **步骤 1：编写 001_sys.up.sql**

```sql
CREATE TABLE gm_user (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    username    VARCHAR(64) NOT NULL,
    password    VARCHAR(128) NOT NULL,
    nickname    VARCHAR(64),
    org_id      VARCHAR(32),
    status      TINYINT NOT NULL DEFAULT 1,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL,
    UNIQUE KEY uk_tenant_username (tenant_id, username)
);

CREATE TABLE gm_role (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    name        VARCHAR(64) NOT NULL,
    code        VARCHAR(64) NOT NULL,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL,
    UNIQUE KEY uk_tenant_code (tenant_id, code)
);

CREATE TABLE gm_org (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    parent_id   VARCHAR(32),
    name        VARCHAR(128) NOT NULL,
    sort_code   INT DEFAULT 0,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL
);

CREATE TABLE gm_menu (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    parent_id   VARCHAR(32),
    title       VARCHAR(128) NOT NULL,
    path        VARCHAR(256),
    component   VARCHAR(256),
    permission  VARCHAR(128),
    type        TINYINT NOT NULL,
    sort_code   INT DEFAULT 0,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL
);

CREATE TABLE gm_sys_config (
    id           VARCHAR(32) PRIMARY KEY,
    tenant_id    VARCHAR(32) NOT NULL DEFAULT 'default',
    config_key   VARCHAR(128) NOT NULL,
    config_value TEXT,
    category     VARCHAR(64),
    remark       VARCHAR(256),
    sort_code    INT DEFAULT 0,
    created_at   DATETIME(3) NOT NULL,
    updated_at   DATETIME(3) NOT NULL,
    deleted_at   DATETIME(3) NULL,
    UNIQUE KEY uk_tenant_key (tenant_id, config_key)
);

CREATE TABLE gm_user_role (
    user_id VARCHAR(32) NOT NULL,
    role_id VARCHAR(32) NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE gm_role_menu (
    role_id VARCHAR(32) NOT NULL,
    menu_id VARCHAR(32) NOT NULL,
    PRIMARY KEY (role_id, menu_id)
);
```

- [ ] **步骤 2：编写 down 迁移并本地验证**

```bash
docker compose up -d mysql
# 使用 golang-migrate 或手动 mysql 执行
mysql -h127.0.0.1 -uroot -proot gin_mam < backend/migrations/001_sys.up.sql
```

- [ ] **步骤 3：Commit**

```bash
git add backend/migrations/
git commit -m "feat(backend): add sys table migrations"
```

---

## 任务 6：ConfigService（DB + Redis 缓存）

**文件：**
- 创建：`backend/internal/infra/config/service.go`
- 创建：`backend/internal/infra/config/service_test.go`
- 创建：`backend/internal/infra/config/mock/service_mock.go`（go:generate）

- [ ] **步骤 1：定义 Service interface + mockgen**

```go
//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Invalidate(ctx context.Context, key string) error
}
```

- [ ] **步骤 2：编写测试（mock Redis）**

验证：DB miss → 空；DB hit → 回填 Redis；Invalidate 后重新读 DB

- [ ] **步骤 3：实现 Get/Set/Invalidate**

Redis key：`config:default:{key}`

- [ ] **步骤 4：运行测试 + generate-mocks**

```bash
cd backend && go generate ./internal/infra/config/... && go test ./internal/infra/config/... -v
```

- [ ] **步骤 5：Commit**

```bash
git commit -m "feat(backend): add ConfigService with Redis cache"
```

---

## 任务 7：sys 域 — 登录 JWT + 用户 CRUD

**文件：**
- 创建：`backend/internal/domain/sys/model/*.go`
- 创建：`backend/internal/domain/sys/repository.go` + `repository_mysql.go`
- 创建：`backend/internal/domain/sys/service.go` + `service_impl.go` + `service_impl_test.go`
- 创建：`backend/internal/domain/sys/handler.go` + `handler_test.go`
- 创建：`backend/internal/domain/sys/mock/*.go`

- [ ] **步骤 1：User model + TableName**

```go
func (User) TableName() string { return "gm_user" }
```

- [ ] **步骤 2：编写 Login 失败测试（用户不存在 → BizError 40100）**

- [ ] **步骤 3：实现 bcrypt 密码校验 + JWT 签发**

- [ ] **步骤 4：实现 Handler**

```
POST /api/v1/auth/login     → 200 { token, expiresAt }
GET  /api/v1/sys/user/page  → 200 { list, total, page, pageSize }
```

- [ ] **步骤 5：运行测试**

```bash
cd backend && go test ./internal/domain/sys/... -v
```

- [ ] **步骤 6：Commit**

```bash
git commit -m "feat(backend): add sys auth and user management"
```

---

## 任务 8：graph.go + main.go + router

**文件：**
- 创建：`backend/internal/app/graph.go`
- 创建：`backend/internal/app/router.go`
- 创建：`backend/cmd/server/main.go`
- 创建：`docker-compose.yml`（根目录）

- [ ] **步骤 1：实现 graph.Build()**

手动注入：mysql → redis → txMgr → configSvc → sysSvc → sysHandler → router

- [ ] **步骤 2：router 注册 /api/v1 + ErrorMiddleware + JWT 中间件**

- [ ] **步骤 3：main.go 读 Viper 启动**

- [ ] **步骤 4：docker-compose.yml**

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: gin_mam
    ports: ["3306:3306"]
  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
  elasticsearch:
    image: elasticsearch:8.15.0
    environment:
      discovery.type: single-node
      xpack.security.enabled: "false"
    ports: ["9200:9200"]
```

- [ ] **步骤 5：冒烟测试**

```bash
docker compose up -d
cd backend && go run ./cmd/server
curl -X POST http://localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin123"}'
```

- [ ] **步骤 6：Commit**

```bash
git commit -m "feat(backend): wire graph.go and start HTTP server (M1 complete)"
```

**M1 验收：** 登录成功返回 camelCase token；错误返回 `{ err_code, err_msg }`；sys 用户分页可用。

---

# 里程碑 M2：媒资 + 编目 + 上传 + ES

## 任务 9：媒资表迁移

**文件：**
- 创建：`backend/migrations/002_asset.up.sql`

表：`gm_asset`, `gm_asset_file`, `gm_asset_metadata`, `gm_catalog`, `gm_catalog_config`

- [ ] **步骤 1：编写迁移 SQL（所有 id varchar(32)，tenant_id，软删）**
- [ ] **步骤 2：执行迁移**
- [ ] **步骤 3：Commit**

---

## 任务 10：S3 Storage（ConfigService 驱动）

**文件：**
- 创建：`backend/internal/infra/storage/storage.go`（interface）
- 创建：`backend/internal/infra/storage/s3.go`
- 创建：`backend/internal/infra/storage/s3_test.go`

- [ ] **步骤 1：Storage interface** — `Put`, `Get`, `Delete`, `PresignPut`
- [ ] **步骤 2：S3 实现从 ConfigService 读 MAM_STORAGE_S3_* keys**
- [ ] **步骤 3：mockgen + 单元测试**
- [ ] **步骤 4：Commit**

---

## 任务 11：asset 域 — CRUD + 预览 URL 拼装

**文件：**
- 创建：`backend/internal/domain/asset/*.go`

- [ ] **步骤 1：GetByID 测试 — mock ConfigService 返回 domain，断言 previewUrl**
- [ ] **步骤 2：实现 Service enrichPreviewURL（tx.NotSupported 读配置）**
- [ ] **步骤 3：Handler CRUD + 分页**
- [ ] **步骤 4：注册路由 `/api/v1/asset/*`**
- [ ] **步骤 5：Commit**

---

## 任务 12：分片上传

**文件：**
- 创建：`backend/internal/domain/asset/upload.go`
- Redis key：`upload:{uploadId}:chunks`

- [ ] **步骤 1：init/chunk/complete 三接口测试**
- [ ] **步骤 2：实现（chunk 大小读 MAM_UPLOAD_CHUNK_SIZE）**
- [ ] **步骤 3：complete 后写 gm_asset_file + S3**
- [ ] **步骤 4：Commit**

---

## 任务 13：catalog 域

- [ ] **步骤 1：树形 catalog CRUD + catalog_config**
- [ ] **步骤 2：测试 + 路由 `/api/v1/catalog/*`**
- [ ] **步骤 3：Commit**

---

## 任务 14：ES 检索

**文件：**
- 创建：`backend/internal/domain/search/*.go`
- 创建：`backend/internal/infra/elasticsearch/client.go`

- [ ] **步骤 1：asset 创建/更新时异步索引（同进程 goroutine MVP）**
- [ ] **步骤 2：POST `/api/v1/search/asset`**
- [ ] **步骤 3：测试 mock ES client**
- [ ] **步骤 4：Commit**

**M2 验收：** 上传文件 → 创建媒资 → ES 可搜 → 详情含 previewUrl。

---

# 里程碑 M3：媒资工作流

## 任务 15：工作流表迁移

**文件：**
- 创建：`backend/migrations/003_workflow.up.sql`

表：`gm_workflow_def`, `gm_workflow_level_user`, `gm_workflow_instance`, `gm_workflow_instance_level_user`, `gm_workflow_operate`

- [ ] **步骤 1：编写迁移**
- [ ] **步骤 2：Commit**

---

## 任务 16：workflow 域 — 状态机 + Redis 锁

**文件：**
- 创建：`backend/internal/domain/workflow/*.go`
- 创建：`backend/internal/infra/lock/redis_lock.go`

- [ ] **步骤 1：状态机纯函数测试（Pending→Auditing→Passed/Rejected）**
- [ ] **步骤 2：Submit 测试（tx.Required，mock assetRepo + workflowRepo）**
- [ ] **步骤 3：Audit + locker_id + Redis wf:lock:{id}**
- [ ] **步骤 4：RequiresNew 写 operate 日志**
- [ ] **步骤 5：Handler 注册 `/api/v1/asset-workflow/*`**
- [ ] **步骤 6：Commit**

**M3 验收：** 送审 → 待审列表 → 通过/打回 → 媒资 status 同步。

---

# 里程碑 M4：转码

## 任务 17：转码表 + 域

**文件：**
- 创建：`backend/migrations/004_transcode.up.sql`
- 创建：`backend/internal/domain/transcode/*.go`

- [ ] **步骤 1：gm_transcode_group + gm_transcode_task 迁移**
- [ ] **步骤 2：创建任务 — 读 MAM_TRANSCODE_URL 调外部服务**
- [ ] **步骤 3：POST `/api/v1/transcode/callback` 签名校验更新状态**
- [ ] **步骤 4：测试 + Commit**

---

# 里程碑 M5：React Web

## 任务 18：Vite + React + Ant Design 脚手架

**文件：**
- 创建：`web/package.json`, `web/vite.config.ts`, `web/vitest.config.ts`

```bash
cd web
npm create vite@latest . -- --template react-ts
npm install antd axios react-router-dom
npm install -D vitest @testing-library/react @testing-library/jest-dom jsdom
```

- [ ] **步骤 1：搭建 Layout + 路由**
- [ ] **步骤 2：Commit**

---

## 任务 19：axios client（无包装响应）

**文件：**
- 创建：`web/src/shared/api/client.ts`
- 创建：`web/src/shared/api/client.test.ts`

```typescript
// 2xx → 直接返回 data
// 4xx/5xx → throw { errCode: data.err_code, errMsg: data.err_msg }
```

- [ ] **步骤 1：编写 client.test.ts 验证拦截器**
- [ ] **步骤 2：Commit**

---

## 任务 20：feature 页面（按优先级）

| 优先级 | Feature | 页面 |
|--------|---------|------|
| P0 | auth | Login |
| P0 | sys | user, role, org, menu, config |
| P0 | asset | list, detail, upload, recycle |
| P0 | workflow | def, audit（四 Tab） |
| P1 | catalog | tree + config |
| P1 | transcode | group, task |
| P2 | ai | 开发中占位 |

- [ ] **每个 feature：组件 + 至少 1 个 Vitest 测试**
- [ ] **Commit 按 feature 拆分**

---

# 里程碑 M6：Flutter 审核端

## 任务 21：Flutter 项目初始化（Ant Design Flutter 生态）

```bash
cd mobile && flutter create . --org com.ginmam
flutter pub add antd_flutter_mobile dio mocktail flutter_riverpod
```

- UI 组件库：**`antd_flutter_mobile`**（Ant Design Mobile，与 Web Ant Design 设计语言一致）
- 不用 `ant_design_flutter`（偏 Web/桌面，官方不推荐移动端）
- [ ] **步骤 1：api_client 解析 ErrorVo（err_code 数值）**
- [ ] **步骤 2：AntApp/主题壳 + 基础路由**
- [ ] **步骤 3：Commit**

---

## 任务 22：auth + audit 三页面

| 页面 | API |
|------|-----|
| LoginPage | POST `/api/v1/auth/login` |
| PendingListPage | POST `/api/v1/asset-workflow/assigned-to-me` |
| AuditDetailPage | GET `/api/v1/asset/:id` + POST audit |

- [ ] **步骤 1：每页 widget test / repository test**
- [ ] **步骤 2：Commit**

**M6 验收：** Flutter 登录 → 看待审 → 通过/打回。

---

# 里程碑 M7：部署

## 任务 23：Dockerfile + Nginx

**文件：**
- 创建：`deploy/Dockerfile`
- 创建：`deploy/nginx/default.conf`

- [ ] **步骤 1：多阶段构建 Go + web/dist + nginx:alpine**
- [ ] **步骤 2：验证容器内 /api/ 反代 + SPA**
- [ ] **步骤 3：Commit**

---

# 规格覆盖自检

| 规格章节 | 对应任务 |
|----------|----------|
| ID varchar(32) UUID v7 | 任务 2 |
| gm_ 单数表 TableName | 任务 5,9,15,17 |
| 配置 L0/L1 分层 | 任务 1,6 |
| graph.go 手动 DI | 任务 8 |
| mockgen | 任务 6,7,10,16 |
| 五传播事务 | 任务 4 |
| REST 无包装 + err_code/err_msg | 任务 3,19 |
| S3 ConfigService | 任务 10 |
| 预览 URL 拼装 | 任务 11 |
| 工作流多级审核 | 任务 16 |
| 转码 | 任务 17 |
| Web 全页面 | 任务 18-20 |
| Flutter 审核 | 任务 21-22 |
| Nginx 单容器 | 任务 23 |
| AI 占位 | 任务 20 router `/ai` → 501 |

---

# 错误码表（实现时创建）

**文件：** `docs/err_codes.md`（任务 3 完成后创建）

| err_code | HTTP | 含义 |
|----------|------|------|
| 40100 | 401 | 未登录 |
| 40300 | 403 | 无权限 |
| 40401 | 404 | 媒资不存在 |
| 40901 | 409 | 工作流锁定 |
| 42201 | 422 | 参数校验失败 |
| 50100 | 501 | 功能开发中 |
| 50000 | 500 | 内部错误 |
