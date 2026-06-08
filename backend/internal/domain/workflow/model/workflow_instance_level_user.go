// Package model 定义工作流领域 GORM 数据库模型。
package model

// TableNameWorkflowInstanceLevelUser 流程实例层级审核人表物理表名。
const TableNameWorkflowInstanceLevelUser = "gm_workflow_instance_level_user"

// WorkflowInstanceLevelUser 流程实例某层级的审核人及该层审核状态，对应表 gm_workflow_instance_level_user。
// 由流程定义复制而来，实例运行期间各层独立记录审核进度。
type WorkflowInstanceLevelUser struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// InstanceID 所属流程实例 ID。
	InstanceID string `gorm:"column:instance_id;size:32;not null"`
	// Level 审核层级（1-based）。
	Level int `gorm:"column:level;not null"`
	// UserID 该层级审核人用户 ID。
	UserID string `gorm:"column:user_id;size:32;not null"`
	// AuditStatus 该层级审核状态：0 待审、1 审核中、2 已通过、3 已打回。
	AuditStatus int8 `gorm:"column:audit_status;not null;default:0"`
}

// TableName 返回 GORM 映射的完整表名 gm_workflow_instance_level_user。
func (WorkflowInstanceLevelUser) TableName() string { return TableNameWorkflowInstanceLevelUser }
