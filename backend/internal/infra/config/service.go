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

var ErrNotFound = errors.New("config: not found")

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
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