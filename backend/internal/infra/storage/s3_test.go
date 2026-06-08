package storage_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/gin-mam/backend/internal/infra/config"
	configmock "github.com/gin-mam/backend/internal/infra/config/mock"
	"github.com/gin-mam/backend/internal/infra/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewS3Storage_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)

	s := storage.NewS3Storage(cfg)
	require.NotNil(t, s)
	var _ storage.Storage = s
}

func TestS3Storage_MissingEndpointOnUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("", config.ErrNotFound)

	s := storage.NewS3Storage(cfg)
	err := s.Put(context.Background(), "k", strings.NewReader("x"), 1, "text/plain")
	require.ErrorIs(t, err, storage.ErrS3EndpointMissing)
}

func TestS3Storage_MissingAccessKeyOnUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("", config.ErrNotFound)

	s := storage.NewS3Storage(cfg)
	err := s.Put(context.Background(), "k", strings.NewReader("x"), 1, "text/plain")
	require.ErrorIs(t, err, storage.ErrS3AccessKeyMissing)
}

func TestS3Storage_MissingSecretKeyOnUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("ak", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("", config.ErrNotFound)

	s := storage.NewS3Storage(cfg)
	err := s.Put(context.Background(), "k", strings.NewReader("x"), 1, "text/plain")
	require.ErrorIs(t, err, storage.ErrS3SecretKeyMissing)
}

func TestS3Storage_MissingBucketOnUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("ak", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("sk", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Bucket).Return("", config.ErrNotFound)

	s := storage.NewS3Storage(cfg)
	err := s.Put(context.Background(), "k", strings.NewReader("x"), 1, "text/plain")
	require.ErrorIs(t, err, storage.ErrS3BucketMissing)
}

func TestS3Storage_InvalidPathStyleOnUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("ak", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("sk", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Bucket).Return("b", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Region).Return("", config.ErrNotFound)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3PathStyle).Return("not-a-bool", nil)

	s := storage.NewS3Storage(cfg)
	err := s.Put(context.Background(), "k", strings.NewReader("x"), 1, "text/plain")
	require.Error(t, err)
	require.False(t, errors.Is(err, storage.ErrS3BucketMissing))
}
