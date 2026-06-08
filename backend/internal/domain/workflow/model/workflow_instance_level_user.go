package model

const TableNameWorkflowInstanceLevelUser = "gm_workflow_instance_level_user"

type WorkflowInstanceLevelUser struct {
	ID          string `gorm:"column:id;primaryKey;size:32"`
	TenantID    string `gorm:"column:tenant_id;size:32;not null;default:default"`
	InstanceID  string `gorm:"column:instance_id;size:32;not null"`
	Level       int    `gorm:"column:level;not null"`
	UserID      string `gorm:"column:user_id;size:32;not null"`
	AuditStatus int8   `gorm:"column:audit_status;not null;default:0"`
}

func (WorkflowInstanceLevelUser) TableName() string { return TableNameWorkflowInstanceLevelUser }
