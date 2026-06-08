// Package workflow 工作流领域：多级审核流程定义与实例运行。
package workflow

import (
	"context"
	"time"
)

const defaultTenant = "default"

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 工作流应用服务接口。
// 实例审核使用 Redis 分布式锁（wf:lock:{id}）防止并发冲突。
type Service interface {
	// Submit 将媒资提交至指定流程定义，创建实例并同步媒资状态为审核中。
	Submit(ctx context.Context, userID string, in SubmitRequest) (*InstanceVO, error)
	// Audit 对单个实例执行通过或打回，校验当前用户是否为该层级审核人。
	Audit(ctx context.Context, userID string, in AuditRequest) error
	// MultiAudit 批量审核多个实例，逐条执行 Audit 逻辑。
	MultiAudit(ctx context.Context, userID string, items []AuditRequest) error
	// Revoke 撤回本人提交的、尚未完成的审核实例。
	Revoke(ctx context.Context, userID string, instanceIDs []string) error
	// ListAssignedToMe 分页查询待当前用户审核的实例（当前层级指派人）。
	ListAssignedToMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error)
	// ListCreatedByMe 分页查询当前用户发起的实例。
	ListCreatedByMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error)
	// ListAuditedByMe 分页查询当前用户已参与审核的实例。
	ListAuditedByMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error)
	// ListAll 分页查询全部实例（管理视角）。
	ListAll(ctx context.Context, page, pageSize int) (*PageResult, error)

	// CreateDef 创建流程定义及各级审核人配置。
	CreateDef(ctx context.Context, in CreateDefInput) (*DefVO, error)
	// DefPage 分页查询流程定义列表。
	DefPage(ctx context.Context, page, pageSize int) (*DefPageResult, error)
}

// SubmitRequest 提交审核入参。
type SubmitRequest struct {
	AssetID       string
	WorkflowDefID string
}

// AuditRequest 单条审核操作入参。
type AuditRequest struct {
	InstanceID string
	Pass       bool
	Remark     string
}

// PageResult 流程实例分页结果。
type PageResult struct {
	List     []InstanceVO `json:"list"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

// InstanceVO 流程实例视图对象。
type InstanceVO struct {
	ID            string    `json:"id"`
	WorkflowDefID string    `json:"workflowDefId"`
	AssetID       string    `json:"assetId"`
	Level         int       `json:"level"`
	AuditLevel    int       `json:"auditLevel"`
	AuditStatus   int8      `json:"auditStatus"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// DefPageResult 流程定义分页结果。
type DefPageResult struct {
	List     []DefVO `json:"list"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

// DefVO 流程定义视图对象。
type DefVO struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	AuditLevel int              `json:"auditLevel"`
	Status     int8             `json:"status"`
	LevelUsers []DefLevelUserVO `json:"levelUsers,omitempty"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
}

// DefLevelUserVO 流程定义某层审核人。
type DefLevelUserVO struct {
	Level  int    `json:"level"`
	UserID string `json:"userId"`
}

// CreateDefInput 创建流程定义入参。
type CreateDefInput struct {
	Name       string
	AuditLevel int
	LevelUsers []DefLevelUserVO
}
