// Package model 定义编目领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameCatalog 编目树节点表物理表名。
const TableNameCatalog = "gm_catalog"

// Catalog 编目树节点，对应表 gm_catalog。
// 通过 ParentID 构成多级树形结构，SortCode 控制同级排序。
type Catalog struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// ParentID 父节点 ID，根节点为空。
	ParentID *string `gorm:"column:parent_id;size:32"`
	// Name 节点显示名称。
	Name string `gorm:"column:name;size:128;not null"`
	// SortCode 同级排序权重，数值越小越靠前。
	SortCode int `gorm:"column:sort_code;default:0"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_catalog。
func (Catalog) TableName() string { return TableNameCatalog }
