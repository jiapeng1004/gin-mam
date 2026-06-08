// Package model 定义转码领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameTranscodeGroup 转码组表物理表名。
const TableNameTranscodeGroup = "gm_transcode_group"

// TranscodeGroup 转码策略组，对应表 gm_transcode_group。
// 定义一组转码参数模板，创建转码任务时引用。
type TranscodeGroup struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// Name 转码组名称。
	Name string `gorm:"column:name;size:128;not null"`
	// Message 转码组说明或备注。
	Message string `gorm:"column:message;type:text"`
	// GroupType 转码组类型，由业务枚举约定。
	GroupType int `gorm:"column:group_type;not null;default:0"`
	// Param 转码参数 JSON 或模板字符串。
	Param string `gorm:"column:param;type:text"`
	// StrategyType 调度策略类型。
	StrategyType int `gorm:"column:strategy_type;not null;default:0"`
	// DefaultFlag 是否默认组：1 是、0 否。
	DefaultFlag int8 `gorm:"column:default_flag;not null;default:0"`
	// AvailableFlag 是否可用：1 可用、0 停用。
	AvailableFlag int8 `gorm:"column:available_flag;not null;default:1"`
	// CreatedBy 创建人用户 ID。
	CreatedBy string `gorm:"column:created_by;size:32;not null"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_transcode_group。
func (TranscodeGroup) TableName() string { return TableNameTranscodeGroup }
