package catalog

import (
	"context"
	"errors"

	"github.com/gin-mam/backend/internal/domain/catalog/model"
	"github.com/gin-mam/backend/internal/infra/tx"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) conn(ctx context.Context) *gorm.DB {
	if db, ok := tx.TxFromContext(ctx); ok {
		return db.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *mysqlRepository) ListCatalogs(ctx context.Context, tenantID string) ([]model.Catalog, error) {
	var rows []model.Catalog
	err := r.conn(ctx).
		Where("tenant_id = ?", tenantID).
		Order("sort_code ASC, created_at ASC").
		Find(&rows).Error
	return rows, err
}

func (r *mysqlRepository) GetCatalogByID(ctx context.Context, tenantID, id string) (*model.Catalog, error) {
	var row model.Catalog
	err := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) CreateCatalog(ctx context.Context, row *model.Catalog) error {
	return r.conn(ctx).Create(row).Error
}

func (r *mysqlRepository) UpdateCatalog(ctx context.Context, row *model.Catalog) error {
	return r.conn(ctx).Save(row).Error
}

func (r *mysqlRepository) DeleteCatalog(ctx context.Context, tenantID, id string) error {
	res := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.Catalog{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *mysqlRepository) ListConfigs(ctx context.Context, tenantID, catalogID string) ([]model.CatalogConfig, error) {
	var rows []model.CatalogConfig
	err := r.conn(ctx).
		Where("tenant_id = ? AND catalog_id = ?", tenantID, catalogID).
		Order("config_key ASC").
		Find(&rows).Error
	return rows, err
}

func (r *mysqlRepository) UpsertConfig(ctx context.Context, row *model.CatalogConfig) error {
	var existing model.CatalogConfig
	err := r.conn(ctx).
		Where("tenant_id = ? AND catalog_id = ? AND config_key = ?", row.TenantID, row.CatalogID, row.ConfigKey).
		First(&existing).Error
	switch {
	case err == nil:
		existing.ConfigValue = row.ConfigValue
		existing.UpdatedAt = row.UpdatedAt
		return r.conn(ctx).Save(&existing).Error
	case errors.Is(err, gorm.ErrRecordNotFound):
		return r.conn(ctx).Create(row).Error
	default:
		return err
	}
}

func (r *mysqlRepository) DeleteConfigsNotInKeys(ctx context.Context, tenantID, catalogID string, keys []string) error {
	q := r.conn(ctx).Where("tenant_id = ? AND catalog_id = ?", tenantID, catalogID)
	if len(keys) > 0 {
		q = q.Where("config_key NOT IN ?", keys)
	}
	return q.Delete(&model.CatalogConfig{}).Error
}