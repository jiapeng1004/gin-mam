// Package config 定义 L1 业务配置 GORM 数据库模型。
package config

import (
	"time"

	"gorm.io/gorm"
)

// TableNameSysConfig 系统业务配置表物理表名。
const TableNameSysConfig = "gm_sys_config"

// SysConfig 系统级业务配置项（L1），对应表 gm_sys_config。
// 存储 S3、预览域名、转码地址等运行时配置，经 Redis 缓存后供各 domain 读取。
type SysConfig struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// ConfigKey 配置键，如 MAM_IMAGE_ACCESS_DOMAIN。
	ConfigKey string `gorm:"column:config_key;size:128;not null"`
	// ConfigValue 配置值。
	ConfigValue string `gorm:"column:config_value;type:text"`
	// Category 配置分类，如 STORAGE、PREVIEW、TRANSCODE。
	Category string `gorm:"column:category;size:64"`
	// Remark 配置说明。
	Remark string `gorm:"column:remark;size:256"`
	// SortCode 同分类下展示排序。
	SortCode int `gorm:"column:sort_code;default:0"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_sys_config。
func (SysConfig) TableName() string { return TableNameSysConfig }
