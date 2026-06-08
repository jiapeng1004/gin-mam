// Package model 定义媒资领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameAsset 媒资主表物理表名。
const TableNameAsset = "gm_asset"

// Asset 媒资主实体，对应表 gm_asset。
// 存储媒资元信息与业务状态；文件路径存于 gm_asset_file，扩展元数据存于 gm_asset_metadata。
type Asset struct {
	// ID 主键，UUID v7 去横杠后的 32 位十六进制字符串。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识，MVP 固定为 default。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// Title 媒资标题，用于列表展示与检索。
	Title string `gorm:"column:title;size:256;not null"`
	// Type 媒资类型，如 video、image、audio、document 等。
	Type string `gorm:"column:type;size:32;not null"`
	// Status 业务状态：0 草稿、1 审核中、2 已通过、3 已打回。
	Status int8 `gorm:"column:status;not null;default:0"`
	// CatalogID 所属编目节点 ID，可为空表示未归类。
	CatalogID *string `gorm:"column:catalog_id;size:32"`
	// Description 媒资描述或备注。
	Description string `gorm:"column:description;type:text"`
	// CreatedBy 创建人用户 ID。
	CreatedBy string `gorm:"column:created_by;size:32;not null"`
	// CreatedAt 记录创建时间（UTC）。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间（UTC）。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳，非空表示已逻辑删除。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_asset。
func (Asset) TableName() string { return TableNameAsset }
