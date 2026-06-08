package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "gm_user"

type User struct {
	ID        string         `gorm:"column:id;primaryKey;size:32"`
	TenantID  string         `gorm:"column:tenant_id;size:32;not null;default:default"`
	Username  string         `gorm:"column:username;size:64;not null"`
	Password  string         `gorm:"column:password;size:128;not null"`
	Nickname  string         `gorm:"column:nickname;size:64"`
	OrgID     *string        `gorm:"column:org_id;size:32"`
	Status    int8           `gorm:"column:status;not null;default:1"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string { return TableNameUser }
