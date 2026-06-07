package sys

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/sys/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=repository.go -destination=mock/repository_mock.go -package=mock
type Repository interface {
	FindUserByUsername(ctx context.Context, tenantID, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, tenantID, id string) error
	GetUserByID(ctx context.Context, tenantID, id string) (*model.User, error)
	ListUsers(ctx context.Context, tenantID string, page, pageSize int) ([]model.User, int64, error)
}
