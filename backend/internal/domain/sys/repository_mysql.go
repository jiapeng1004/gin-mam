package sys

import (
	"context"
	"errors"

	"github.com/gin-mam/backend/internal/domain/sys/model"
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

func (r *mysqlRepository) FindUserByUsername(ctx context.Context, tenantID, username string) (*model.User, error) {
	var user model.User
	err := r.conn(ctx).
		Where("tenant_id = ? AND username = ?", tenantID, username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mysqlRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.conn(ctx).Create(user).Error
}

func (r *mysqlRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.conn(ctx).Save(user).Error
}

func (r *mysqlRepository) DeleteUser(ctx context.Context, tenantID, id string) error {
	res := r.conn(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&model.User{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *mysqlRepository) GetUserByID(ctx context.Context, tenantID, id string) (*model.User, error) {
	var user model.User
	err := r.conn(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mysqlRepository) ListUsers(ctx context.Context, tenantID string, page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.conn(ctx).Model(&model.User{}).Where("tenant_id = ?", tenantID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []model.User
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

var _ Repository = (*mysqlRepository)(nil)

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
