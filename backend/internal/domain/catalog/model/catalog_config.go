// Package model 定义编目领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameCatalogConfig 编目节点配置表物理表名。
const TableNameCatalogConfig = "gm_catalog_config"

// CatalogConfig 编目节点级业务配置项，对应表 gm_catalog_config。
// 与全局 gm_sys_config 不同，配置仅作用于指定编目及其子树业务逻辑。
type CatalogConfig struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// CatalogID 所属编目节点 ID。
	CatalogID string `gorm:"column:catalog_id;size:32;not null"`
	// ConfigKey 配置键名。
	ConfigKey string `gorm:"column:config_key;size:128;not null"`
	// ConfigValue 配置值。
	ConfigValue string `gorm:"column:config_value;type:text"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_catalog_config。
func (CatalogConfig) TableName() string { return TableNameCatalogConfig }
