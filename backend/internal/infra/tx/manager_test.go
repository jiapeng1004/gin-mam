package tx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testItem struct {
	ID string `gorm:"primaryKey"`
}

func (testItem) TableName() string { return "test_items" }

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&testItem{}))
	return db
}

func dbForCtx(root *gorm.DB, ctx context.Context) *gorm.DB {
	if gdb, ok := tx.TxFromContext(ctx); ok {
		return gdb.WithContext(ctx)
	}
	return root.WithContext(ctx)
}

func TestMandatory_WithoutTx_ReturnsError(t *testing.T) {
	db := openTestDB(t)
	mgr := tx.NewManager(db)

	err := mgr.Run(context.Background(), tx.Mandatory, func(ctx context.Context) error {
		return nil
	})
	require.ErrorIs(t, err, tx.ErrTxMandatory)
}

func TestNever_WithTx_ReturnsError(t *testing.T) {
	db := openTestDB(t)
	mgr := tx.NewManager(db)

	ctx := tx.ContextWithTx(context.Background(), db)
	err := mgr.Run(ctx, tx.Never, func(ctx context.Context) error {
		return nil
	})
	require.ErrorIs(t, err, tx.ErrTxNever)
}

func TestRequired_CommitsOnSuccess(t *testing.T) {
	db := openTestDB(t)
	mgr := tx.NewManager(db)

	err := mgr.Run(context.Background(), tx.Required, func(ctx context.Context) error {
		return dbForCtx(db, ctx).Create(&testItem{ID: "committed"}).Error
	})
	require.NoError(t, err)

	var item testItem
	require.NoError(t, db.First(&item, "id = ?", "committed").Error)
	require.Equal(t, "committed", item.ID)
}

func TestRequiresNew_InnerRollbackOuterCommits(t *testing.T) {
	db := openTestDB(t)
	mgr := tx.NewManager(db)
	innerFail := errors.New("rollback inner")

	err := mgr.Run(context.Background(), tx.Required, func(ctx context.Context) error {
		outerTx, ok := tx.TxFromContext(ctx)
		require.True(t, ok)

		innerErr := mgr.Run(ctx, tx.RequiresNew, func(innerCtx context.Context) error {
			innerTx, ok := tx.TxFromContext(innerCtx)
			require.True(t, ok)
			require.NotSame(t, outerTx, innerTx)
			return innerFail
		})
		require.ErrorIs(t, innerErr, innerFail)

		return dbForCtx(db, ctx).Create(&testItem{ID: "outer"}).Error
	})
	require.NoError(t, err)

	var count int64
	require.NoError(t, db.Model(&testItem{}).Where("id = ?", "outer").Count(&count).Error)
	require.Equal(t, int64(1), count)
	require.NoError(t, db.Model(&testItem{}).Where("id = ?", "inner").Count(&count).Error)
	require.Equal(t, int64(0), count)
}

