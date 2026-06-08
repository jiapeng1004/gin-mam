// Package model 定义系统管理领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameRole 角色表物理表名。
const TableNameRole = "gm_role"

// Role 系统角色，对应表 gm_role。
// 用于 RBAC 权限模型，MVP 阶段表结构已预留。
type Role struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// Name 角色显示名称。
	Name string `gorm:"column:name;size:64;not null"`
	// Code 角色编码，租户内唯一，用于程序判断。
	Code string `gorm:"column:code;size:64;not null"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_role。
func (Role) TableName() string { return TableNameRole }
