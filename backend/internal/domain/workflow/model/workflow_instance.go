// Package model 定义工作流领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameWorkflowInstance 工作流实例表物理表名。
const TableNameWorkflowInstance = "gm_workflow_instance"

// WorkflowInstance 单次媒资审核流程实例，对应表 gm_workflow_instance。
// Level 表示当前已完成层级；AuditStatus 表示实例整体审核结论。
type WorkflowInstance struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// WorkflowDefID 引用的流程定义 ID。
	WorkflowDefID string `gorm:"column:workflow_def_id;size:32;not null"`
	// AssetID 被审核的媒资 ID。
	AssetID string `gorm:"column:asset_id;size:32;not null"`
	// Level 当前审核进度层级（0 表示尚未完成任何层级）。
	Level int `gorm:"column:level;not null;default:0"`
	// AuditLevel 流程定义的总审核层级数（冗余快照）。
	AuditLevel int `gorm:"column:audit_level;not null"`
	// AuditStatus 实例状态：0 待审、1 审核中、2 已通过、3 已打回。
	AuditStatus int8 `gorm:"column:audit_status;not null;default:0"`
	// LockerID 当前持锁用户 ID，用于 Redis 分布式锁关联，可为空。
	LockerID *string `gorm:"column:locker_id;size:32"`
	// CreatedBy 提交审核的用户 ID。
	CreatedBy string `gorm:"column:created_by;size:32;not null"`
	// CreatedAt 实例创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 实例最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_workflow_instance。
func (WorkflowInstance) TableName() string { return TableNameWorkflowInstance }
