package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameTranscodeProfile = "gm_transcode_profile"

type TranscodeProfile struct {
	ID             string         `gorm:"column:id;primaryKey;size:32"`
	TenantID       string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	Name           string         `gorm:"column:name;size:128;not null"`
	Alias          string         `gorm:"column:alias;size:128"`
	TranscodeType  int            `gorm:"column:transcode_type;not null;default:0"`
	DefinitionType int            `gorm:"column:definition_type;not null;default:0"`
	Param          string         `gorm:"column:param;type:text"`
	SpecialParam   string         `gorm:"column:special_param;type:text"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (TranscodeProfile) TableName() string { return TableNameTranscodeProfile }