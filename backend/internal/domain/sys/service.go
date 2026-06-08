// Package sys 系统管理领域：认证、用户管理与种子数据。
package sys

import (
	"context"
	"time"

	"github.com/gin-mam/backend/internal/domain/sys/model"
)

const defaultTenant = "default"

const (
	// ErrCodeAuthFailed 登录失败（用户名或密码错误）的业务错误码。
	ErrCodeAuthFailed = 40100
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 系统管理应用服务接口。
type Service interface {
	// Login 校验用户名密码并签发 JWT，返回 token 与过期时间。
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	// UserPage 分页查询用户列表（不含密码）。
	UserPage(ctx context.Context, page, pageSize int) (*UserPageResult, error)
	// EnsureSeedAdmin 启动时确保存在默认管理员（admin/admin123），幂等。
	EnsureSeedAdmin(ctx context.Context) error
	// CreateUser 创建新用户，密码入库前哈希。
	CreateUser(ctx context.Context, in CreateUserInput) (*UserVO, error)
	// UpdateUser 更新用户资料，Password 非空时重置密码。
	UpdateUser(ctx context.Context, id string, in UpdateUserInput) (*UserVO, error)
	// DeleteUser 软删除用户。
	DeleteUser(ctx context.Context, id string) error
	// GetUser 按 ID 获取用户详情。
	GetUser(ctx context.Context, id string) (*UserVO, error)
}

// JWTConfig JWT 签发配置，由 L0 config.yaml 注入。
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// LoginResult 登录成功响应体。
type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// UserVO 用户视图对象，不含敏感字段。
type UserVO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname,omitempty"`
	OrgID     *string   `json:"orgId,omitempty"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserPageResult 用户分页结果。
type UserPageResult struct {
	List     []UserVO `json:"list"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

// CreateUserInput 创建用户入参。
type CreateUserInput struct {
	Username string
	Password string
	Nickname string
	OrgID    *string
	Status   int8
}

// UpdateUserInput 更新用户入参，指针字段 nil 表示不修改。
type UpdateUserInput struct {
	Nickname *string
	OrgID    *string
	Status   *int8
	Password *string
}

func userToVO(u *model.User) UserVO {
	return UserVO{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		OrgID:     u.OrgID,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
