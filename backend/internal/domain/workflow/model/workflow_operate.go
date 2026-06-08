// Package model 定义工作流领域 GORM 数据库模型。
package model

import "time"

// TableNameWorkflowOperate 工作流操作审计表物理表名。
const TableNameWorkflowOperate = "gm_workflow_operate"

// WorkflowOperate 流程实例操作记录，对应表 gm_workflow_operate。
// 记录提交、通过、打回、撤回等操作，便于审计追溯。
type WorkflowOperate struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// InstanceID 关联的流程实例 ID。
	InstanceID string `gorm:"column:instance_id;size:32;not null"`
	// OperatorID 操作人用户 ID。
	OperatorID string `gorm:"column:operator_id;size:32;not null"`
	// Action 操作类型：1 提交、2 通过、3 打回、4 撤回。
	Action int8 `gorm:"column:action;not null"`
	// Remark 操作备注或审核意见。
	Remark string `gorm:"column:remark;type:text"`
	// CreatedAt 操作发生时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

// TableName 返回 GORM 映射的完整表名 gm_workflow_operate。
func (WorkflowOperate) TableName() string { return TableNameWorkflowOperate }
