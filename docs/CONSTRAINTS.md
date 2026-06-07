# GIN-MAM 项目约束

本文档为 GIN-MAM 开发的硬性约定，所有代码必须遵守。

## ID

1. 所有业务表主键类型：`varchar(32)`
2. 生成规则：UUID v7，存储时去除横杠（32 位 hex）
3. 禁止：自增 ID、雪花 ID、带横杠 UUID、同一实体存在第二标识（如 contentSourceId）
4. 外键命名：`{entity}_id`，类型一律 `varchar(32)`

## 表命名

1. 前缀：`gm_`
2. 表名使用**单数**，禁止复数
3. 禁止 GORM 全局 `NamingStrategy.TablePrefix`
4. 每个实体必须实现 `TableName()` 返回完整表名，例如 `gm_asset`

## 租户

1. 所有业务表含 `tenant_id varchar(32)`
2. MVP 默认值 `"default"`，不做运行时隔离
3. 二期启用中间件自动注入 tenant 过滤

## 配置

1. **L0（YAML）**：仅 MySQL、Redis、ES、HTTP 端口、日志级别、JWT Secret
2. **L1（gm_sys_config）**：S3、预览链、转码、分享、功能开关等一切业务参数
3. **L2（gm_catalog_config）**：栏目级编目配置
4. 禁止将 S3 endpoint、预览域名等业务参数写入 YAML
5. 读取统一经 `ConfigService`，禁止业务代码直查配置表
6. Redis 缓存 key：`config:{tenant_id}:{config_key}`，变更时主动失效
7. 配置 key 命名：`MAM_{CATEGORY}_{NAME}` 全大写

## 依赖注入

1. 禁止 wire、dig 等 DI 框架
2. 统一在 `internal/app/graph.go` 手动组装依赖
3. 可 mock 的层必须定义 interface

## Mock

1. 统一使用 mockgen（`go.uber.org/mock/mockgen`）
2. 禁止 mockery、testify/mock 混用
3. 每个 interface 源文件顶部声明 `//go:generate mockgen ...`
4. mock 输出目录：`{package}/mock/`

## 测试

1. 前后端均必须有单元测试，CI 强制执行
2. 后端核心 domain 覆盖率目标 ≥ 70%
3. 前端：Vitest + Testing Library
4. Flutter：mocktail 或 mockito

## 事务

1. 禁止 Service 层直接调用 `db.Transaction()`
2. 统一使用 `infra/tx.Manager.Run(ctx, propagation, fn)`
3. 传播类型：Required（默认）、RequiresNew、Mandatory、NotSupported、Never
4. Repository 通过 `TxFromContext(ctx)` 获取当前事务
5. 各传播行为必须有 `manager_test.go` 单测

## HTTP 响应

1. **成功**：HTTP 2xx，直接返回业务 JSON，字段 **camelCase**，无外层包装
2. **失败**：HTTP 4xx/5xx，返回 `{ "err_code": <number>, "err_msg": "<string>" }`
3. 错误 JSON 的键名使用 **snake_case**（`err_code`、`err_msg`），与成功体 camelCase 字段天然隔离
4. `err_code` 值为**数值型**业务错误码，不是字符串枚举
5. 禁止 ResultDTO、ResultVo、CommonResult、`{ code, data, message }` 包装
6. Handler 禁止手动拼装错误 JSON，统一使用 `httpx.Fail` 或 `BizError`

## 软删除与时间

1. 统一 GORM 软删除 `deleted_at`
2. 时间字段：`created_at` / `updated_at`，UTC 存储
