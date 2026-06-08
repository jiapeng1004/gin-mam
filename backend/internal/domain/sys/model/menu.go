// Package model 定义系统管理领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameMenu 菜单/权限资源表物理表名。
const TableNameMenu = "gm_menu"

// Menu 前端菜单与权限资源，对应表 gm_menu。
// Type 区分目录、菜单、按钮等节点类型。
type Menu struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// ParentID 父菜单 ID，根菜单为空。
	ParentID *string `gorm:"column:parent_id;size:32"`
	// Title 菜单标题。
	Title string `gorm:"column:title;size:128;not null"`
	// Path 前端路由路径。
	Path string `gorm:"column:path;size:256"`
	// Component 前端组件路径或名称。
	Component string `gorm:"column:component;size:256"`
	// Permission 权限标识符，用于后端鉴权。
	Permission string `gorm:"column:permission;size:128"`
	// Type 节点类型：如 0 目录、1 菜单、2 按钮（具体枚举由业务约定）。
	Type int8 `gorm:"column:type;not null"`
	// SortCode 同级排序权重。
	SortCode int `gorm:"column:sort_code;default:0"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_menu。
func (Menu) TableName() string { return TableNameMenu }
