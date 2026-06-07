package model

import "time"

const TableNameAssetMetadata = "gm_asset_metadata"

type AssetMetadata struct {
	ID        string    `gorm:"column:id;primaryKey;size:32"`
	TenantID  string    `gorm:"column:tenant_id;size:32;not null;default:default"`
	AssetID   string    `gorm:"column:asset_id;size:32;not null"`
	MetaKey   string    `gorm:"column:meta_key;size:128;not null"`
	MetaValue string    `gorm:"column:meta_value;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (AssetMetadata) TableName() string { return TableNameAssetMetadata }