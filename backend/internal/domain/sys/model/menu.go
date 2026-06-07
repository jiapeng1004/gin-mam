package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameMenu = "gm_menu"

type Menu struct {
	ID         string         `gorm:"column:id;primaryKey;size:32"`
	TenantID   string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	ParentID   *string        `gorm:"column:parent_id;size:32"`
	Title      string         `gorm:"column:title;size:128;not null"`
	Path       string         `gorm:"column:path;size:256"`
	Component  string         `gorm:"column:component;size:256"`
	Permission string         `gorm:"column:permission;size:128"`
	Type       int8           `gorm:"column:type;not null"`
	SortCode   int            `gorm:"column:sort_code;default:0"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Menu) TableName() string { return TableNameMenu }
