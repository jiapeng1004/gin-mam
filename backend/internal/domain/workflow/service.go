package workflow

import (
	"context"
	"time"
)

const defaultTenant = "default"

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Submit(ctx context.Context, userID string, in SubmitRequest) (*InstanceVO, error)
	Audit(ctx context.Context, userID string, in AuditRequest) error
	MultiAudit(ctx context.Context, userID string, items []AuditRequest) error
	Revoke(ctx context.Context, userID string, instanceIDs []string) error
	ListAssignedToMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error)
	ListCreatedByMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error)
	ListAuditedByMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error)
	ListAll(ctx context.Context, page, pageSize int) (*PageResult, error)

	CreateDef(ctx context.Context, in CreateDefInput) (*DefVO, error)
	DefPage(ctx context.Context, page, pageSize int) (*DefPageResult, error)
}

type SubmitRequest struct {
	AssetID        string
	WorkflowDefID  string
}

type AuditRequest struct {
	InstanceID string
	Pass       bool
	Remark     string
}

type PageResult struct {
	List     []InstanceVO `json:"list"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

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

type DefPageResult struct {
	List     []DefVO `json:"list"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

type DefVO struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	AuditLevel int              `json:"auditLevel"`
	Status     int8             `json:"status"`
	LevelUsers []DefLevelUserVO `json:"levelUsers,omitempty"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
}

type DefLevelUserVO struct {
	Level  int    `json:"level"`
	UserID string `json:"userId"`
}

type CreateDefInput struct {
	Name       string
	AuditLevel int
	LevelUsers []DefLevelUserVO
}
