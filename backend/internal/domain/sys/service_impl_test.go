package sys_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/domain/sys/model"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{
		NowFunc: func() time.Time { return time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))
	return db
}

func seedUser(t *testing.T, db *gorm.DB, username, password string) *model.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	now := db.NowFunc()
	user := &model.User{
		ID:        "user001",
		TenantID:  "default",
		Username:  username,
		Password:  string(hash),
		Nickname:  username,
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, db.Create(user).Error)
	return user
}

func TestLogin_WrongPassword_ReturnsBizError40100(t *testing.T) {
	db := openTestDB(t)
	seedUser(t, db, "admin", "admin123")
	repo := sys.NewMySQLRepository(db)
	svc := sys.NewService(repo, sys.JWTConfig{Secret: "test-secret", ExpireHours: 24})

	_, err := svc.Login(context.Background(), "admin", "wrong-password")
	require.Error(t, err)
	biz, ok := err.(*httpx.BizError)
	require.True(t, ok)
	require.Equal(t, 401, biz.HTTPStatus)
	require.Equal(t, sys.ErrCodeAuthFailed, biz.ErrCode)
	require.Equal(t, "用户名或密码错误", biz.ErrMsg)
}

func TestLogin_Success_ReturnsToken(t *testing.T) {
	db := openTestDB(t)
	seedUser(t, db, "admin", "admin123")
	repo := sys.NewMySQLRepository(db)
	svc := sys.NewService(repo, sys.JWTConfig{Secret: "test-secret", ExpireHours: 24})

	result, err := svc.Login(context.Background(), "admin", "admin123")
	require.NoError(t, err)
	require.NotEmpty(t, result.Token)
	require.False(t, result.ExpiresAt.IsZero())
}

func TestEnsureSeedAdmin_CreatesAdmin(t *testing.T) {
	db := openTestDB(t)
	repo := sys.NewMySQLRepository(db)
	svc := sys.NewService(repo, sys.JWTConfig{Secret: "test-secret", ExpireHours: 24})
	ctx := context.Background()

	require.NoError(t, svc.EnsureSeedAdmin(ctx))
	require.NoError(t, svc.EnsureSeedAdmin(ctx))

	_, err := svc.Login(ctx, "admin", "admin123")
	require.NoError(t, err)
}

