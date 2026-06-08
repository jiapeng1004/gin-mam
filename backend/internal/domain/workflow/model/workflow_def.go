package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameWorkflowDef = "gm_workflow_def"

type WorkflowDef struct {
	ID         string         `gorm:"column:id;primaryKey;size:32"`
	TenantID   string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	Name       string         `gorm:"column:name;size:128;not null"`
	AuditLevel int            `gorm:"column:audit_level;not null"`
	Status     int8           `gorm:"column:status;not null;default:1"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (WorkflowDef) TableName() string { return TableNameWorkflowDef }
