package config_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-mam/backend/internal/infra/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		NowFunc: func() time.Time { return time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&config.SysConfig{}))
	return db
}

func newTestRedis(t *testing.T) (*miniredis.Miniredis, redis.Cmdable) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	return mr, redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func TestGet_CacheMiss_DBHit_PopulatesRedis(t *testing.T) {
	db := openTestDB(t)
	mr, rdb := newTestRedis(t)
	svc := config.NewService(db, rdb)
	ctx := context.Background()

	require.NoError(t, db.Create(&config.SysConfig{
		ID:          "cfg1",
		TenantID:    "default",
		ConfigKey:   "MAM_TEST_KEY",
		ConfigValue: "from-db",
		CreatedAt:   db.NowFunc(),
		UpdatedAt:   db.NowFunc(),
	}).Error)

	val, err := svc.Get(ctx, "MAM_TEST_KEY")
	require.NoError(t, err)
	require.Equal(t, "from-db", val)
	cached, err := mr.Get("config:default:MAM_TEST_KEY")
	require.NoError(t, err)
	require.Equal(t, "from-db", cached)

	val2, err := svc.Get(ctx, "MAM_TEST_KEY")
	require.NoError(t, err)
	require.Equal(t, "from-db", val2)
}

func TestGet_NotFound(t *testing.T) {
	db := openTestDB(t)
	_, rdb := newTestRedis(t)
	svc := config.NewService(db, rdb)

	_, err := svc.Get(context.Background(), "MISSING")
	require.ErrorIs(t, err, config.ErrNotFound)
}

func TestInvalidate_NextGetReadsDBAgain(t *testing.T) {
	db := openTestDB(t)
	mr, rdb := newTestRedis(t)
	svc := config.NewService(db, rdb)
	ctx := context.Background()
	key := "MAM_INVALIDATE"

	require.NoError(t, db.Create(&config.SysConfig{
		ID:          "cfg2",
		TenantID:    "default",
		ConfigKey:   key,
		ConfigValue: "v1",
		CreatedAt:   db.NowFunc(),
		UpdatedAt:   db.NowFunc(),
	}).Error)

	_, err := svc.Get(ctx, key)
	require.NoError(t, err)
	cached, err := mr.Get("config:default:" + key)
	require.NoError(t, err)
	require.Equal(t, "v1", cached)

	require.NoError(t, db.Model(&config.SysConfig{}).
		Where("config_key = ?", key).
		Update("config_value", "v2").Error)

	valCached, err := svc.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, "v1", valCached)

	require.NoError(t, svc.Invalidate(ctx, key))
	require.False(t, mr.Exists("config:default:"+key))

	valFresh, err := svc.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, "v2", valFresh)
}

func TestSet_WritesDBAndInvalidatesCache(t *testing.T) {
	db := openTestDB(t)
	mr, rdb := newTestRedis(t)
	svc := config.NewService(db, rdb)
	ctx := context.Background()
	key := "MAM_SET"

	mr.Set("config:default:"+key, "stale")

	require.NoError(t, svc.Set(ctx, key, "new-value"))

	var row config.SysConfig
	require.NoError(t, db.Where("config_key = ?", key).First(&row).Error)
	require.Equal(t, "new-value", row.ConfigValue)
	require.False(t, mr.Exists("config:default:"+key))

	val, err := svc.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, "new-value", val)
}