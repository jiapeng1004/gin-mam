package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameRole = "gm_role"

type Role struct {
	ID        string         `gorm:"column:id;primaryKey;size:32"`
	TenantID  string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	Name      string         `gorm:"column:name;size:64;not null"`
	Code      string         `gorm:"column:code;size:64;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Role) TableName() string { return TableNameRole }
