package sys

import (
	"context"
	"errors"
	"time"

	"github.com/gin-mam/backend/internal/domain/sys/model"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/gin-mam/backend/internal/pkg/id"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	repo Repository
	jwt  JWTConfig
	now  func() time.Time
}

func NewService(repo Repository, jwtCfg JWTConfig) Service {
	if jwtCfg.ExpireHours <= 0 {
		jwtCfg.ExpireHours = 24
	}
	return &service{
		repo: repo,
		jwt:  jwtCfg,
		now:  time.Now,
	}
}

func (s *service) authFailed() *httpx.BizError {
	return httpx.NewBizError(401, ErrCodeAuthFailed, "用户名或密码错误")
}

type tokenClaims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.repo.FindUserByUsername(ctx, defaultTenant, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, s.authFailed()
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, s.authFailed()
	}
	expiresAt := s.now().Add(time.Duration(s.jwt.ExpireHours) * time.Hour)
	claims := tokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(s.now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwt.Secret))
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, ExpiresAt: expiresAt}, nil
}

func (s *service) UserPage(ctx context.Context, page, pageSize int) (*UserPageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	users, total, err := s.repo.ListUsers(ctx, defaultTenant, page, pageSize)
	if err != nil {
		return nil, err
	}
	list := make([]UserVO, 0, len(users))
	for i := range users {
		list = append(list, userToVO(&users[i]))
	}
	return &UserPageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *service) EnsureSeedAdmin(ctx context.Context) error {
	_, err := s.repo.FindUserByUsername(ctx, defaultTenant, "admin")
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := s.now()
	return s.repo.CreateUser(ctx, &model.User{
		ID:        id.New(),
		TenantID:  defaultTenant,
		Username:  "admin",
		Password:  string(hash),
		Nickname:  "管理员",
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (s *service) CreateUser(ctx context.Context, in CreateUserInput) (*UserVO, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := s.now()
	status := in.Status
	if status == 0 {
		status = 1
	}
	user := &model.User{
		ID:        id.New(),
		TenantID:  defaultTenant,
		Username:  in.Username,
		Password:  string(hash),
		Nickname:  in.Nickname,
		OrgID:     in.OrgID,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	vo := userToVO(user)
	return &vo, nil
}

func (s *service) UpdateUser(ctx context.Context, id string, in UpdateUserInput) (*UserVO, error) {
	user, err := s.repo.GetUserByID(ctx, defaultTenant, id)
	if err != nil {
		return nil, err
	}
	if in.Nickname != nil {
		user.Nickname = *in.Nickname
	}
	if in.OrgID != nil {
		user.OrgID = in.OrgID
	}
	if in.Status != nil {
		user.Status = *in.Status
	}
	if in.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hash)
	}
	user.UpdatedAt = s.now()
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	vo := userToVO(user)
	return &vo, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, defaultTenant, id)
}

func (s *service) GetUser(ctx context.Context, id string) (*UserVO, error) {
	user, err := s.repo.GetUserByID(ctx, defaultTenant, id)
	if err != nil {
		return nil, err
	}
	vo := userToVO(user)
	return &vo, nil
}
