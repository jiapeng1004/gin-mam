package model

import "time"

const TableNameWorkflowOperate = "gm_workflow_operate"

type WorkflowOperate struct {
	ID         string    `gorm:"column:id;primaryKey;size:32"`
	TenantID   string    `gorm:"column:tenant_id;size:32;not null;default:default"`
	InstanceID string    `gorm:"column:instance_id;size:32;not null"`
	OperatorID string    `gorm:"column:operator_id;size:32;not null"`
	Action     int8      `gorm:"column:action;not null"`
	Remark     string    `gorm:"column:remark;type:text"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
}

func (WorkflowOperate) TableName() string { return TableNameWorkflowOperate }
