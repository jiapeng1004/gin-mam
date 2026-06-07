package catalog

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/catalog/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=repository.go -destination=mock/repository_mock.go -package=mock
type Repository interface {
	ListCatalogs(ctx context.Context, tenantID string) ([]model.Catalog, error)
	GetCatalogByID(ctx context.Context, tenantID, id string) (*model.Catalog, error)
	CreateCatalog(ctx context.Context, row *model.Catalog) error
	UpdateCatalog(ctx context.Context, row *model.Catalog) error
	DeleteCatalog(ctx context.Context, tenantID, id string) error

	ListConfigs(ctx context.Context, tenantID, catalogID string) ([]model.CatalogConfig, error)
	UpsertConfig(ctx context.Context, row *model.CatalogConfig) error
	DeleteConfigsNotInKeys(ctx context.Context, tenantID, catalogID string, keys []string) error
}