// Package model 定义系统管理领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameOrg 组织架构表物理表名。
const TableNameOrg = "gm_org"

// Org 组织架构节点，对应表 gm_org。
// 通过 ParentID 构成多级组织树。
type Org struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// ParentID 上级组织 ID，根组织为空。
	ParentID *string `gorm:"column:parent_id;size:32"`
	// Name 组织名称。
	Name string `gorm:"column:name;size:128;not null"`
	// SortCode 同级排序权重。
	SortCode int `gorm:"column:sort_code;default:0"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_org。
func (Org) TableName() string { return TableNameOrg }
