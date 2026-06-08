package sys

import (
	"context"
	"time"

	"github.com/gin-mam/backend/internal/domain/sys/model"
)

const defaultTenant = "default"

const (
	ErrCodeAuthFailed = 40100
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	UserPage(ctx context.Context, page, pageSize int) (*UserPageResult, error)
	EnsureSeedAdmin(ctx context.Context) error
	CreateUser(ctx context.Context, in CreateUserInput) (*UserVO, error)
	UpdateUser(ctx context.Context, id string, in UpdateUserInput) (*UserVO, error)
	DeleteUser(ctx context.Context, id string) error
	GetUser(ctx context.Context, id string) (*UserVO, error)
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type UserVO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname,omitempty"`
	OrgID     *string   `json:"orgId,omitempty"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserPageResult struct {
	List     []UserVO `json:"list"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

type CreateUserInput struct {
	Username string
	Password string
	Nickname string
	OrgID    *string
	Status   int8
}

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
