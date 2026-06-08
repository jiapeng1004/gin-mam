package catalog

import (
	"context"
)

const defaultTenant = "default"

const ErrCodeNotFound = 40402

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Tree(ctx context.Context) ([]TreeNode, error)
	Create(ctx context.Context, in CreateInput) (*CatalogVO, error)
	Update(ctx context.Context, id string, in UpdateInput) (*CatalogVO, error)
	Delete(ctx context.Context, id string) error
	ListConfig(ctx context.Context, catalogID string) ([]ConfigVO, error)
	BatchUpdateConfig(ctx context.Context, catalogID string, items []ConfigItemInput) ([]ConfigVO, error)
}

type TreeNode struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	ParentID *string    `json:"parentId"`
	Children []TreeNode `json:"children"`
}

type CatalogVO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parentId,omitempty"`
	SortCode int     `json:"sortCode"`
}

type CreateInput struct {
	Name     string
	ParentID *string
	SortCode int
}

type UpdateInput struct {
	Name     *string
	ParentID *string
	SortCode *int
}

type ConfigVO struct {
	ID          string `json:"id"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
}

type ConfigItemInput struct {
	ConfigKey   string
	ConfigValue string
}