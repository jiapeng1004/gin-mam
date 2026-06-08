// Package model 定义媒资领域 GORM 数据库模型。
package model

import "time"

// TableNameAssetMetadata 媒资扩展元数据表物理表名。
const TableNameAssetMetadata = "gm_asset_metadata"

// AssetMetadata 媒资键值对扩展元数据，对应表 gm_asset_metadata。
// 用于存储非结构化或业务自定义字段，同一 AssetID 下 MetaKey 应唯一。
type AssetMetadata struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// AssetID 关联媒资 ID。
	AssetID string `gorm:"column:asset_id;size:32;not null"`
	// MetaKey 元数据键名。
	MetaKey string `gorm:"column:meta_key;size:128;not null"`
	// MetaValue 元数据值，文本存储。
	MetaValue string `gorm:"column:meta_value;type:text"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

// TableName 返回 GORM 映射的完整表名 gm_asset_metadata。
func (AssetMetadata) TableName() string { return TableNameAssetMetadata }
