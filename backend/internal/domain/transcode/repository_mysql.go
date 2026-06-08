package transcode

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/transcode/model"
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

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return page, pageSize
}

func (r *mysqlRepository) CreateGroup(ctx context.Context, row *model.TranscodeGroup) error {
	return r.conn(ctx).Create(row).Error
}

func (r *mysqlRepository) GetGroupByID(ctx context.Context, tenantID, id string) (*model.TranscodeGroup, error) {
	var row model.TranscodeGroup
	err := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) UpdateGroup(ctx context.Context, row *model.TranscodeGroup) error {
	return r.conn(ctx).Save(row).Error
}

func (r *mysqlRepository) DeleteGroup(ctx context.Context, tenantID, id string) error {
	res := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.TranscodeGroup{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *mysqlRepository) ListGroups(ctx context.Context, tenantID string, page, pageSize int) ([]model.TranscodeGroup, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	q := r.conn(ctx).Model(&model.TranscodeGroup{}).Where("tenant_id = ?", tenantID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.TranscodeGroup
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) CreateTask(ctx context.Context, row *model.TranscodeTask) error {
	return r.conn(ctx).Create(row).Error
}

func (r *mysqlRepository) GetTaskByID(ctx context.Context, tenantID, id string) (*model.TranscodeTask, error) {
	var row model.TranscodeTask
	err := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) UpdateTask(ctx context.Context, row *model.TranscodeTask) error {
	return r.conn(ctx).Save(row).Error
}

func (r *mysqlRepository) ListTasks(ctx context.Context, tenantID string, page, pageSize int) ([]model.TranscodeTask, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	q := r.conn(ctx).Model(&model.TranscodeTask{}).Where("tenant_id = ?", tenantID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.TranscodeTask
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}