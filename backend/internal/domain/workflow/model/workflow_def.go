// Package model 定义工作流领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameWorkflowDef 工作流定义表物理表名。
const TableNameWorkflowDef = "gm_workflow_def"

// WorkflowDef 审核流程定义模板，对应表 gm_workflow_def。
// AuditLevel 表示总审核层级数；各级审核人见 gm_workflow_level_user。
type WorkflowDef struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// Name 流程定义名称。
	Name string `gorm:"column:name;size:128;not null"`
	// AuditLevel 审核总层级数（从 1 开始计数）。
	AuditLevel int `gorm:"column:audit_level;not null"`
	// Status 定义状态：1 启用、0 停用。
	Status int8 `gorm:"column:status;not null;default:1"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_workflow_def。
func (WorkflowDef) TableName() string { return TableNameWorkflowDef }
