// Package model 定义系统管理领域 GORM 数据库模型。
package model

import (
	"time"

	"gorm.io/gorm"
)

// TableNameUser 系统用户表物理表名。
const TableNameUser = "gm_user"

// User 系统用户，对应表 gm_user。
// 密码字段存储哈希值，不存明文。
type User struct {
	// ID 主键，UUID v7 去横杠 32 位。
	ID string `gorm:"column:id;primaryKey;size:32"`
	// TenantID 租户标识。
	TenantID string `gorm:"column:tenant_id;size:32;not null;default:default"`
	// Username 登录用户名，租户内唯一。
	Username string `gorm:"column:username;size:64;not null"`
	// Password 密码哈希（bcrypt 等）。
	Password string `gorm:"column:password;size:128;not null"`
	// Nickname 用户昵称或显示名。
	Nickname string `gorm:"column:nickname;size:64"`
	// OrgID 所属组织 ID，可为空。
	OrgID *string `gorm:"column:org_id;size:32"`
	// Status 账号状态：1 启用、0 禁用。
	Status int8 `gorm:"column:status;not null;default:1"`
	// CreatedAt 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	// UpdatedAt 记录最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	// DeletedAt 软删除时间戳。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回 GORM 映射的完整表名 gm_user。
func (User) TableName() string { return TableNameUser }
