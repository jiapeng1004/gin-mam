package storage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-mam/backend/internal/infra/config"
	configmock "github.com/gin-mam/backend/internal/infra/config/mock"
	"github.com/gin-mam/backend/internal/infra/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func fakeS3Config(ctrl *gomock.Controller) *configmock.MockService {
	m := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	m.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	m.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("test-access", nil)
	m.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("test-secret", nil)
	m.EXPECT().Get(ctx, storage.ConfigKeyS3Bucket).Return("test-bucket", nil)
	m.EXPECT().Get(ctx, storage.ConfigKeyS3Region).Return("us-east-1", nil)
	m.EXPECT().Get(ctx, storage.ConfigKeyS3PathStyle).Return("true", nil)
	return m
}

func TestNewS3Storage_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := fakeS3Config(ctrl)

	s, err := storage.NewS3Storage(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, s)
	var _ storage.Storage = s
}

func TestNewS3Storage_MissingEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("", config.ErrNotFound)

	_, err := storage.NewS3Storage(context.Background(), cfg)
	require.ErrorIs(t, err, storage.ErrS3EndpointMissing)
}

func TestNewS3Storage_MissingAccessKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("", config.ErrNotFound)

	_, err := storage.NewS3Storage(context.Background(), cfg)
	require.ErrorIs(t, err, storage.ErrS3AccessKeyMissing)
}

func TestNewS3Storage_MissingSecretKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("ak", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("", config.ErrNotFound)

	_, err := storage.NewS3Storage(context.Background(), cfg)
	require.ErrorIs(t, err, storage.ErrS3SecretKeyMissing)
}

func TestNewS3Storage_MissingBucket(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("ak", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("sk", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Bucket).Return("", config.ErrNotFound)

	_, err := storage.NewS3Storage(context.Background(), cfg)
	require.ErrorIs(t, err, storage.ErrS3BucketMissing)
}

func TestNewS3Storage_InvalidPathStyle(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	ctx := gomock.Any()
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Endpoint).Return("http://127.0.0.1:9000", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3AccessKey).Return("ak", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3SecretKey).Return("sk", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Bucket).Return("b", nil)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3Region).Return("", config.ErrNotFound)
	cfg.EXPECT().Get(ctx, storage.ConfigKeyS3PathStyle).Return("not-a-bool", nil)

	_, err := storage.NewS3Storage(context.Background(), cfg)
	require.Error(t, err)
	require.False(t, errors.Is(err, storage.ErrS3BucketMissing))
}