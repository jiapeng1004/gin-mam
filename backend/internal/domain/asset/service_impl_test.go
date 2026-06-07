package asset_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/asset/model"
	configmock "github.com/gin-mam/backend/internal/infra/config/mock"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openAssetTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{
		NowFunc: func() time.Time { return time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Asset{}, &model.AssetFile{}, &model.AssetMetadata{}))
	return db
}

func seedAssetWithFile(t *testing.T, db *gorm.DB) string {
	t.Helper()
	now := db.NowFunc()
	assetID := "asset001"
	require.NoError(t, db.Create(&model.Asset{
		ID:        assetID,
		TenantID:  "default",
		Title:     "demo",
		Type:      "image",
		Status:    0,
		CreatedBy: "u1",
		CreatedAt: now,
		UpdatedAt: now,
	}).Error)
	require.NoError(t, db.Create(&model.AssetFile{
		ID:          "file001",
		TenantID:    "default",
		AssetID:     assetID,
		Version:     1,
		StoragePath: "bucket/path/photo.jpg",
		IsMaster:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}).Error)
	return assetID
}

func TestGetByID_EnrichesPreviewURL_WhenConfigReturnsDomain(t *testing.T) {
	db := openAssetTestDB(t)
	assetID := seedAssetWithFile(t, db)

	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	cfg.EXPECT().
		Get(gomock.Any(), asset.ConfigKeyImageAccess).
		Return("https://img.example.com/", nil)

	repo := asset.NewMySQLRepository(db)
	svc := asset.NewService(repo, cfg, tx.NewManager(db))

	vo, err := svc.GetByID(context.Background(), assetID)
	require.NoError(t, err)
	require.Equal(t, "https://img.example.com/bucket/path/photo.jpg", vo.PreviewURL)
}