package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCatalogConfig = "gm_catalog_config"

type CatalogConfig struct {
	ID          string         `gorm:"column:id;primaryKey;size:32"`
	TenantID    string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	CatalogID   string         `gorm:"column:catalog_id;size:32;not null"`
	ConfigKey   string         `gorm:"column:config_key;size:128;not null"`
	ConfigValue string         `gorm:"column:config_value;type:text"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (CatalogConfig) TableName() string { return TableNameCatalogConfig }