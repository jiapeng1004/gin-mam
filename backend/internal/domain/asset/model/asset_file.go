package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameAssetFile = "gm_asset_file"

type AssetFile struct {
	ID          string         `gorm:"column:id;primaryKey;size:32"`
	TenantID    string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	AssetID     string         `gorm:"column:asset_id;size:32;not null"`
	Version     int            `gorm:"column:version;not null;default:1"`
	StoragePath string         `gorm:"column:storage_path;size:512;not null"`
	MimeType    string         `gorm:"column:mime_type;size:128"`
	FileSize    *int64         `gorm:"column:file_size"`
	Checksum    string         `gorm:"column:checksum;size:64"`
	IsMaster    int8           `gorm:"column:is_master;not null;default:1"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (AssetFile) TableName() string { return TableNameAssetFile }