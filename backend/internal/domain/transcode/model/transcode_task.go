package model

import (
	"time"
)

const TableNameTranscodeTask = "gm_transcode_task"

const (
	TaskStatusPending int8 = 0
	TaskStatusRunning int8 = 1
	TaskStatusSuccess int8 = 2
	TaskStatusFailed  int8 = 3
)

type TranscodeTask struct {
	ID               string    `gorm:"column:id;primaryKey;size:32"`
	TenantID         string    `gorm:"column:tenant_id;size:32;not null;default:default"`
	AssetID          string    `gorm:"column:asset_id;size:32;not null"`
	AssetFileID      *string   `gorm:"column:asset_file_id;size:32"`
	TranscodeGroupID string    `gorm:"column:transcode_group_id;size:32;not null"`
	ProfileID        *string   `gorm:"column:profile_id;size:32"`
	Status           int8      `gorm:"column:status;not null;default:0"`
	ExternalJobID    string    `gorm:"column:external_job_id;size:128"`
	OutputPath       string    `gorm:"column:output_path;size:512"`
	ErrorMsg         string    `gorm:"column:error_msg;type:text"`
	CallbackPayload  *string   `gorm:"column:callback_payload;type:text"`
	CreatedBy        string    `gorm:"column:created_by;size:32;not null"`
	CreatedAt        time.Time `gorm:"column:created_at;not null"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null"`
}

func (TranscodeTask) TableName() string { return TableNameTranscodeTask }