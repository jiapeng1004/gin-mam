package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-mam/backend/internal/domain/catalog"
	"github.com/gin-mam/backend/internal/domain/catalog/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openCatalogTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{
		NowFunc: func() time.Time { return time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Catalog{}, &model.CatalogConfig{}))
	return db
}

func TestTree_BuildsNestedChildren(t *testing.T) {
	db := openCatalogTestDB(t)
	now := db.NowFunc()
	require.NoError(t, db.Create(&model.Catalog{
		ID: "root", TenantID: "default", Name: "root", SortCode: 0, CreatedAt: now, UpdatedAt: now,
	}).Error)
	parent := "root"
	require.NoError(t, db.Create(&model.Catalog{
		ID: "child", TenantID: "default", ParentID: &parent, Name: "child", SortCode: 0, CreatedAt: now, UpdatedAt: now,
	}).Error)

	svc := catalog.NewService(catalog.NewMySQLRepository(db))
	tree, err := svc.Tree(context.Background())
	require.NoError(t, err)
	require.Len(t, tree, 1)
	require.Equal(t, "root", tree[0].ID)
	require.Len(t, tree[0].Children, 1)
	require.Equal(t, "child", tree[0].Children[0].ID)
}

func TestBatchUpdateConfig_UpsertAndPrune(t *testing.T) {
	db := openCatalogTestDB(t)
	now := db.NowFunc()
	require.NoError(t, db.Create(&model.Catalog{
		ID: "cat1", TenantID: "default", Name: "c", CreatedAt: now, UpdatedAt: now,
	}).Error)
	require.NoError(t, db.Create(&model.CatalogConfig{
		ID: "cfg1", TenantID: "default", CatalogID: "cat1", ConfigKey: "old", ConfigValue: "x", CreatedAt: now, UpdatedAt: now,
	}).Error)

	svc := catalog.NewService(catalog.NewMySQLRepository(db))
	out, err := svc.BatchUpdateConfig(context.Background(), "cat1", []catalog.ConfigItemInput{
		{ConfigKey: "new", ConfigValue: "y"},
	})
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "new", out[0].ConfigKey)
	require.Equal(t, "y", out[0].ConfigValue)
}