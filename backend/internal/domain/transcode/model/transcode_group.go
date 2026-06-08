package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameTranscodeGroup = "gm_transcode_group"

type TranscodeGroup struct {
	ID            string         `gorm:"column:id;primaryKey;size:32"`
	TenantID      string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	Name          string         `gorm:"column:name;size:128;not null"`
	Message       string         `gorm:"column:message;type:text"`
	GroupType     int            `gorm:"column:group_type;not null;default:0"`
	Param         string         `gorm:"column:param;type:text"`
	StrategyType  int            `gorm:"column:strategy_type;not null;default:0"`
	DefaultFlag   int8           `gorm:"column:default_flag;not null;default:0"`
	AvailableFlag int8           `gorm:"column:available_flag;not null;default:1"`
	CreatedBy     string         `gorm:"column:created_by;size:32;not null"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (TranscodeGroup) TableName() string { return TableNameTranscodeGroup }