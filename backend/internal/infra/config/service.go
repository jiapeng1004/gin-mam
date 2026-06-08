// Package config L1 业务配置服务：gm_sys_config + Redis 缓存。
package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-mam/backend/internal/pkg/id"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const defaultTenant = "default"

// ErrNotFound 配置键不存在时返回的错误。
var ErrNotFound = errors.New("config: not found")

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 业务配置读写服务接口。
// 读取顺序：Redis 缓存 → MySQL gm_sys_config；写入时 upsert 并失效缓存。
type Service interface {
	// Get 按配置键读取值，缓存未命中时回源数据库并回填 Redis。
	Get(ctx context.Context, key string) (string, error)
	// Set 写入或更新配置值，并清除对应 Redis 缓存。
	Set(ctx context.Context, key, value string) error
	// Invalidate 仅清除指定键的 Redis 缓存。
	Invalidate(ctx context.Context, key string) error
}

type service struct {
	db  *gorm.DB
	rdb redis.Cmdable
}

func NewService(db *gorm.DB, rdb redis.Cmdable) Service {
	return &service{db: db, rdb: rdb}
}

func cacheKey(key string) string {
	return fmt.Sprintf("config:%s:%s", defaultTenant, key)
}

func (s *service) Get(ctx context.Context, key string) (string, error) {
	redisKey := cacheKey(key)
	val, err := s.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		return val, nil
	}
	if !errors.Is(err, redis.Nil) {
		return "", err
	}

	var row SysConfig
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND config_key = ?", defaultTenant, key).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	if setErr := s.rdb.Set(ctx, redisKey, row.ConfigValue, 0).Err(); setErr != nil {
		return row.ConfigValue, setErr
	}
	return row.ConfigValue, nil
}

func (s *service) Set(ctx context.Context, key, value string) error {
	now := s.db.NowFunc()
	var row SysConfig
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND config_key = ?", defaultTenant, key).
		First(&row).Error
	switch {
	case err == nil:
		row.ConfigValue = value
		row.UpdatedAt = now
		if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
			return err
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		row = SysConfig{
			ID:          id.New(),
			TenantID:    defaultTenant,
			ConfigKey:   key,
			ConfigValue: value,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
			return err
		}
	default:
		return err
	}
	return s.Invalidate(ctx, key)
}

func (s *service) Invalidate(ctx context.Context, key string) error {
	return s.rdb.Del(ctx, cacheKey(key)).Err()
}