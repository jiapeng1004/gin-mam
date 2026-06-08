// Package model 定义转码领域 GORM 数据库模型。
package model

import (
	"time"
)

// TableNameTranscodeTask 转码任务表物理表名。
const TableNameTranscodeTask = "gm_transcode_task"

// 转码任务状态枚举。
const (
	// TaskStatusPending 待提交或排队中。
	TaskStatusPending int8 = 0
	// TaskStatusRunning 转码进行中。
	TaskStatusRunning int8 = 1
	// TaskStatusSuccess 转码成功。
	TaskStatusSuccess int8 = 2
	// TaskStatusFailed 转码失败。
	TaskStatusFailed int8 = 3
)

// TranscodeTask 单次媒资转码任务，对应表 gm_transcode_task。
// 任务提交至外部转码服务，结果通过 HMAC 签名回调更新。
type TranscodeTask struct {
	// ID 主键，UUID v7 去横杠 32 位，同时作为回调 taskId。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// AssetID 关联媒资 ID。
	AssetID string `gorm:"column:asset_id;size:32;not null"`
	// AssetFileID 源文件 ID，可为空表示使用媒资主文件。
	AssetFileID *string `gorm:"column:asset_file_id;size:32"`
	// TranscodeGroupID 使用的转码组 ID。
	TranscodeGroupID string `gorm:"column:transcode_group_id;size:32;not null"`
	// ProfileID 使用的转码模板 ID，可为空。
	ProfileID *string `gorm:"column:profile_id;size:32"`
	// Status 任务状态，见 TaskStatus* 常量。
	Status int8 `gorm:"column:status;not null;default:0"`
	// ExternalJobID 外部转码服务返回的任务 ID。
	ExternalJobID string `gorm:"column:external_job_id;size:128"`
	// OutputPath 转码输出文件在对象存储中的路径。
	OutputPath string `gorm:"column:output_path;size:512"`
	// ErrorMsg 失败时的错误信息。
	ErrorMsg string `gorm:"column:error_msg;type:text"`
	// CallbackPayload 最近一次回调原始 JSON，便于排查。
	CallbackPayload *string `gorm:"column:callback_payload;type:text"`
	// CreatedBy 创建任务的用户 ID。
	CreatedBy string `gorm:"column:created_by;size:32;not null"`
	// CreatedAt 任务创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 任务最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

// TableName 返回 GORM 映射的完整表名 gm_transcode_task。
func (TranscodeTask) TableName() string { return TableNameTranscodeTask }
