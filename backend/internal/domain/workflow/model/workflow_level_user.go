package model

const TableNameWorkflowLevelUser = "gm_workflow_level_user"

type WorkflowLevelUser struct {
	ID            string `gorm:"column:id;primaryKey;size:32"`
	TenantID      string `gorm:"column:tenant_id;size:32;not null;default:default"`
	WorkflowDefID string `gorm:"column:workflow_def_id;size:32;not null"`
	Level         int    `gorm:"column:level;not null"`
	UserID        string `gorm:"column:user_id;size:32;not null"`
}

func (WorkflowLevelUser) TableName() string { return TableNameWorkflowLevelUser }
