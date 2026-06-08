package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameWorkflowInstance = "gm_workflow_instance"

type WorkflowInstance struct {
	ID            string         `gorm:"column:id;primaryKey;size:32"`
	TenantID      string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	WorkflowDefID string         `gorm:"column:workflow_def_id;size:32;not null"`
	AssetID       string         `gorm:"column:asset_id;size:32;not null"`
	Level         int            `gorm:"column:level;not null;default:0"`
	AuditLevel    int            `gorm:"column:audit_level;not null"`
	AuditStatus   int8           `gorm:"column:audit_status;not null;default:0"`
	LockerID      *string        `gorm:"column:locker_id;size:32"`
	CreatedBy     string         `gorm:"column:created_by;size:32;not null"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (WorkflowInstance) TableName() string { return TableNameWorkflowInstance }
