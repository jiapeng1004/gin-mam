// Package model 定义工作流领域 GORM 数据库模型。
package model

// TableNameWorkflowLevelUser 流程定义层级审核人表物理表名。
const TableNameWorkflowLevelUser = "gm_workflow_level_user"

// WorkflowLevelUser 流程定义中某层级的指定审核人，对应表 gm_workflow_level_user。
// Level 从 1 开始，与 WorkflowDef.AuditLevel 配合使用。
type WorkflowLevelUser struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// WorkflowDefID 所属流程定义 ID。
	WorkflowDefID string `gorm:"column:workflow_def_id;size:32;not null"`
	// Level 审核层级（1-based）。
	Level int `gorm:"column:level;not null"`
	// UserID 该层级的审核人用户 ID。
	UserID string `gorm:"column:user_id;size:32;not null"`
}

// TableName 返回 GORM 映射的完整表名 gm_workflow_level_user。
func (WorkflowLevelUser) TableName() string { return TableNameWorkflowLevelUser }
