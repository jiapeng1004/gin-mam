package asset

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-mam/backend/internal/domain/asset/model"
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

func (r *mysqlRepository) GetAssetByID(ctx context.Context, tenantID, id string) (*model.Asset, error) {
	var row model.Asset
	err := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) ListAssets(ctx context.Context, tenantID string, page, pageSize int, keyword string) ([]model.Asset, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.conn(ctx).Model(&model.Asset{}).Where("tenant_id = ?", tenantID)
	if kw := strings.TrimSpace(keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("title LIKE ?", like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Asset
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) CreateAsset(ctx context.Context, asset *model.Asset) error {
	return r.conn(ctx).Create(asset).Error
}

func (r *mysqlRepository) UpdateAsset(ctx context.Context, asset *model.Asset) error {
	return r.conn(ctx).Save(asset).Error
}

func (r *mysqlRepository) DeleteAsset(ctx context.Context, tenantID, id string) error {
	res := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.Asset{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *mysqlRepository) GetMasterFile(ctx context.Context, tenantID, assetID string) (*model.AssetFile, error) {
	var row model.AssetFile
	err := r.conn(ctx).
		Where("tenant_id = ? AND asset_id = ? AND is_master = 1", tenantID, assetID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) CreateAssetFile(ctx context.Context, file *model.AssetFile) error {
	return r.conn(ctx).Create(file).Error
}

func (r *mysqlRepository) ListMetadata(ctx context.Context, tenantID, assetID string) ([]model.AssetMetadata, error) {
	var rows []model.AssetMetadata
	err := r.conn(ctx).
		Where("tenant_id = ? AND asset_id = ?", tenantID, assetID).
		Order("meta_key ASC").
		Find(&rows).Error
	return rows, err
}

func (r *mysqlRepository) UpsertMetadata(ctx context.Context, row *model.AssetMetadata) error {
	var existing model.AssetMetadata
	err := r.conn(ctx).
		Where("tenant_id = ? AND asset_id = ? AND meta_key = ?", row.TenantID, row.AssetID, row.MetaKey).
		First(&existing).Error
	switch {
	case err == nil:
		existing.MetaValue = row.MetaValue
		existing.UpdatedAt = row.UpdatedAt
		return r.conn(ctx).Save(&existing).Error
	case errors.Is(err, gorm.ErrRecordNotFound):
		return r.conn(ctx).Create(row).Error
	default:
		return err
	}
}

func (r *mysqlRepository) DeleteMetadataNotInKeys(ctx context.Context, tenantID, assetID string, keys []string) error {
	q := r.conn(ctx).Where("tenant_id = ? AND asset_id = ?", tenantID, assetID)
	if len(keys) > 0 {
		q = q.Where("meta_key NOT IN ?", keys)
	}
	return q.Delete(&model.AssetMetadata{}).Error
}