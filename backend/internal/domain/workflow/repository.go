package workflow

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/workflow/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=repository.go -destination=mock/repository_mock.go -package=mock
type Repository interface {
	CreateDef(ctx context.Context, def *model.WorkflowDef) error
	GetDefByID(ctx context.Context, tenantID, id string) (*model.WorkflowDef, error)
	ListDefs(ctx context.Context, tenantID string, page, pageSize int) ([]model.WorkflowDef, int64, error)
	CreateLevelUsers(ctx context.Context, rows []model.WorkflowLevelUser) error
	ListLevelUsersByDefID(ctx context.Context, tenantID, defID string) ([]model.WorkflowLevelUser, error)

	CreateInstance(ctx context.Context, inst *model.WorkflowInstance) error
	GetInstanceByID(ctx context.Context, tenantID, id string) (*model.WorkflowInstance, error)
	UpdateInstance(ctx context.Context, inst *model.WorkflowInstance) error
	DeleteInstance(ctx context.Context, tenantID, id string) error
	ListInstancesAssignedToUser(ctx context.Context, tenantID, userID string, page, pageSize int) ([]model.WorkflowInstance, int64, error)
	ListInstancesCreatedBy(ctx context.Context, tenantID, userID string, page, pageSize int) ([]model.WorkflowInstance, int64, error)
	ListInstancesAuditedBy(ctx context.Context, tenantID, userID string, page, pageSize int) ([]model.WorkflowInstance, int64, error)
	ListAllInstances(ctx context.Context, tenantID string, page, pageSize int) ([]model.WorkflowInstance, int64, error)

	CreateInstanceLevelUsers(ctx context.Context, rows []model.WorkflowInstanceLevelUser) error
	IsUserAssigneeAtLevel(ctx context.Context, tenantID, instanceID string, level int, userID string) (bool, error)
	UpdateInstanceLevelUserStatus(ctx context.Context, tenantID, instanceID string, level int, userID string, status int8) error

	CreateOperate(ctx context.Context, row *model.WorkflowOperate) error
}
