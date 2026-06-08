package transcode

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/transcode/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=repository.go -destination=mock/repository_mock.go -package=mock
type Repository interface {
	CreateGroup(ctx context.Context, row *model.TranscodeGroup) error
	GetGroupByID(ctx context.Context, tenantID, id string) (*model.TranscodeGroup, error)
	UpdateGroup(ctx context.Context, row *model.TranscodeGroup) error
	DeleteGroup(ctx context.Context, tenantID, id string) error
	ListGroups(ctx context.Context, tenantID string, page, pageSize int) ([]model.TranscodeGroup, int64, error)

	CreateTask(ctx context.Context, row *model.TranscodeTask) error
	GetTaskByID(ctx context.Context, tenantID, id string) (*model.TranscodeTask, error)
	UpdateTask(ctx context.Context, row *model.TranscodeTask) error
	ListTasks(ctx context.Context, tenantID string, page, pageSize int) ([]model.TranscodeTask, int64, error)
}