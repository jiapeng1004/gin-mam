// Package model 定义转码领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameTranscodeProfile 转码模板表物理表名。
const TableNameTranscodeProfile = "gm_transcode_profile"

// TranscodeProfile 转码输出规格模板，对应表 gm_transcode_profile。
// 描述清晰度、编码格式等转码输出参数，可被转码组引用。
type TranscodeProfile struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// Name 模板名称。
	Name string `gorm:"column:name;size:128;not null"`
	// Alias 模板别名或对外显示名。
	Alias string `gorm:"column:alias;size:128"`
	// TranscodeType 转码类型枚举。
	TranscodeType int `gorm:"column:transcode_type;not null;default:0"`
	// DefinitionType 清晰度/规格类型枚举。
	DefinitionType int `gorm:"column:definition_type;not null;default:0"`
	// Param 标准转码参数 JSON。
	Param string `gorm:"column:param;type:text"`
	// SpecialParam 特殊场景附加参数 JSON。
	SpecialParam string `gorm:"column:special_param;type:text"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_transcode_profile。
func (TranscodeProfile) TableName() string { return TableNameTranscodeProfile }
