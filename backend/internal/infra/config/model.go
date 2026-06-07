package config

import (
	"time"

	"gorm.io/gorm"
)

const TableNameSysConfig = "gm_sys_config"

// SysConfig maps gm_sys_config (L1 business configuration).
type SysConfig struct {
	ID          string         `gorm:"column:id;primaryKey;size:32"`
	TenantID    string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	ConfigKey   string         `gorm:"column:config_key;size:128;not null"`
	ConfigValue string         `gorm:"column:config_value;type:text"`
	Category    string         `gorm:"column:category;size:64"`
	Remark      string         `gorm:"column:remark;size:256"`
	SortCode    int            `gorm:"column:sort_code;default:0"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SysConfig) TableName() string { return TableNameSysConfig }