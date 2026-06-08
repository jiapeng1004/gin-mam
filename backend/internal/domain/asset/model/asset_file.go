// Package model 定义媒资领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameAssetFile 媒资文件表物理表名。
const TableNameAssetFile = "gm_asset_file"

// AssetFile 媒资物理文件记录，对应表 gm_asset_file。
// 同一媒资可有多版本文件，通过 Version 区分；IsMaster 标记当前主文件。
type AssetFile struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// AssetID 关联的媒资主表 ID。
	AssetID string `gorm:"column:asset_id;size:32;not null"`
	// Version 文件版本号，从 1 递增。
	Version int `gorm:"column:version;not null;default:1"`
	// StoragePath 对象存储中的相对路径，预览 URL 由 ConfigService 拼装。
	StoragePath string `gorm:"column:storage_path;size:512;not null"`
	// MimeType 文件 MIME 类型。
	MimeType string `gorm:"column:mime_type;size:128"`
	// FileSize 文件字节大小，可为空。
	FileSize *int64 `gorm:"column:file_size"`
	// Checksum 文件校验和（如 MD5/SHA256）。
	Checksum string `gorm:"column:checksum;size:64"`
	// IsMaster 是否主文件：1 是、0 否。
	IsMaster int8 `gorm:"column:is_master;not null;default:1"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_asset_file。
func (AssetFile) TableName() string { return TableNameAssetFile }
