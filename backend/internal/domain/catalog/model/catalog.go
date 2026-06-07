package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCatalog = "gm_catalog"

type Catalog struct {
	ID        string         `gorm:"column:id;primaryKey;size:32"`
	TenantID  string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	ParentID  *string        `gorm:"column:parent_id;size:32"`
	Name      string         `gorm:"column:name;size:128;not null"`
	SortCode  int            `gorm:"column:sort_code;default:0"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Catalog) TableName() string { return TableNameCatalog }