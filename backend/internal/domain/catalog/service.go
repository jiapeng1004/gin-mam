// Package catalog 编目领域：树形栏目 CRUD 与节点级配置管理。
package catalog

import (
	"context"
)

const defaultTenant = "default"

// ErrCodeNotFound 编目节点不存在时的业务错误码。
const ErrCodeNotFound = 40402

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 编目应用服务接口。
type Service interface {
	// Tree 返回完整编目树，根节点 ParentID 为空。
	Tree(ctx context.Context) ([]TreeNode, error)
	// Create 创建编目节点。
	Create(ctx context.Context, in CreateInput) (*CatalogVO, error)
	// Update 更新编目节点名称、父级或排序。
	Update(ctx context.Context, id string, in UpdateInput) (*CatalogVO, error)
	// Delete 软删除编目节点（需无子节点或按业务规则处理）。
	Delete(ctx context.Context, id string) error
	// ListConfig 列出指定编目节点的配置项。
	ListConfig(ctx context.Context, catalogID string) ([]ConfigVO, error)
	// BatchUpdateConfig 批量 upsert 编目配置项。
	BatchUpdateConfig(ctx context.Context, catalogID string, items []ConfigItemInput) ([]ConfigVO, error)
}

// TreeNode 编目树节点，含递归 Children。
type TreeNode struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	ParentID *string    `json:"parentId"`
	Children []TreeNode `json:"children"`
}

// CatalogVO 编目节点视图对象。
type CatalogVO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parentId,omitempty"`
	SortCode int     `json:"sortCode"`
}

// CreateInput 创建编目入参。
type CreateInput struct {
	Name     string
	ParentID *string
	SortCode int
}

// UpdateInput 更新编目入参，指针字段 nil 表示不修改。
type UpdateInput struct {
	Name     *string
	ParentID *string
	SortCode *int
}

// ConfigVO 编目配置项视图。
type ConfigVO struct {
	ID          string `json:"id"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
}

// ConfigItemInput 批量更新时的单条配置。
type ConfigItemInput struct {
	ConfigKey   string
	ConfigValue string
}
