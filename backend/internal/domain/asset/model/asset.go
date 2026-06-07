package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameAsset = "gm_asset"

type Asset struct {
	ID          string         `gorm:"column:id;primaryKey;size:32"`
	TenantID    string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	Title       string         `gorm:"column:title;size:256;not null"`
	Type        string         `gorm:"column:type;size:32;not null"`
	Status      int8           `gorm:"column:status;not null;default:0"`
	CatalogID   *string        `gorm:"column:catalog_id;size:32"`
	Description string         `gorm:"column:description;type:text"`
	CreatedBy   string         `gorm:"column:created_by;size:32;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Asset) TableName() string { return TableNameAsset }