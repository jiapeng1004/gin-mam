package workflow

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/workflow/model"
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

func (r *mysqlRepository) CreateDef(ctx context.Context, def *model.WorkflowDef) error {
	return r.conn(ctx).Create(def).Error
}

func (r *mysqlRepository) GetDefByID(ctx context.Context, tenantID, id string) (*model.WorkflowDef, error) {
	var row model.WorkflowDef
	err := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) ListDefs(ctx context.Context, tenantID string, page, pageSize int) ([]model.WorkflowDef, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	q := r.conn(ctx).Model(&model.WorkflowDef{}).Where("tenant_id = ?", tenantID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.WorkflowDef
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) CreateLevelUsers(ctx context.Context, rows []model.WorkflowLevelUser) error {
	if len(rows) == 0 {
		return nil
	}
	return r.conn(ctx).Create(&rows).Error
}

func (r *mysqlRepository) ListLevelUsersByDefID(ctx context.Context, tenantID, defID string) ([]model.WorkflowLevelUser, error) {
	var rows []model.WorkflowLevelUser
	err := r.conn(ctx).
		Where("tenant_id = ? AND workflow_def_id = ?", tenantID, defID).
		Order("level ASC, user_id ASC").
		Find(&rows).Error
	return rows, err
}

func (r *mysqlRepository) CreateInstance(ctx context.Context, inst *model.WorkflowInstance) error {
	return r.conn(ctx).Create(inst).Error
}

func (r *mysqlRepository) GetInstanceByID(ctx context.Context, tenantID, id string) (*model.WorkflowInstance, error) {
	var row model.WorkflowInstance
	err := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *mysqlRepository) UpdateInstance(ctx context.Context, inst *model.WorkflowInstance) error {
	return r.conn(ctx).Save(inst).Error
}

func (r *mysqlRepository) DeleteInstance(ctx context.Context, tenantID, id string) error {
	res := r.conn(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.WorkflowInstance{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *mysqlRepository) instanceListQuery(ctx context.Context, tenantID string) *gorm.DB {
	return r.conn(ctx).Model(&model.WorkflowInstance{}).Where("tenant_id = ?", tenantID)
}

func (r *mysqlRepository) ListInstancesAssignedToUser(ctx context.Context, tenantID, userID string, page, pageSize int) ([]model.WorkflowInstance, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	q := r.instanceListQuery(ctx, tenantID).
		Where("audit_status IN ?", []int8{AuditStatusPending, AuditStatusAuditing}).
		Where(`EXISTS (
			SELECT 1 FROM gm_workflow_instance_level_user ilu
			WHERE ilu.instance_id = gm_workflow_instance.id
			  AND ilu.tenant_id = gm_workflow_instance.tenant_id
			  AND ilu.user_id = ?
			  AND ilu.level = gm_workflow_instance.level + 1
		)`, userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.WorkflowInstance
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) ListInstancesCreatedBy(ctx context.Context, tenantID, userID string, page, pageSize int) ([]model.WorkflowInstance, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	q := r.instanceListQuery(ctx, tenantID).Where("created_by = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.WorkflowInstance
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) ListInstancesAuditedBy(ctx context.Context, tenantID, userID string, page, pageSize int) ([]model.WorkflowInstance, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	sub := r.conn(ctx).Model(&model.WorkflowOperate{}).
		Select("DISTINCT instance_id").
		Where("tenant_id = ? AND operator_id = ? AND action IN ?", tenantID, userID, []int8{OperateActionPass, OperateActionReject})
	q := r.instanceListQuery(ctx, tenantID).Where("id IN (?)", sub)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.WorkflowInstance
	offset := (page - 1) * pageSize
	err := q.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) ListAllInstances(ctx context.Context, tenantID string, page, pageSize int) ([]model.WorkflowInstance, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	q := r.instanceListQuery(ctx, tenantID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.WorkflowInstance
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}

func (r *mysqlRepository) CreateInstanceLevelUsers(ctx context.Context, rows []model.WorkflowInstanceLevelUser) error {
	if len(rows) == 0 {
		return nil
	}
	return r.conn(ctx).Create(&rows).Error
}

func (r *mysqlRepository) IsUserAssigneeAtLevel(ctx context.Context, tenantID, instanceID string, level int, userID string) (bool, error) {
	var count int64
	err := r.conn(ctx).Model(&model.WorkflowInstanceLevelUser{}).
		Where("tenant_id = ? AND instance_id = ? AND level = ? AND user_id = ?", tenantID, instanceID, level, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *mysqlRepository) UpdateInstanceLevelUserStatus(ctx context.Context, tenantID, instanceID string, level int, userID string, status int8) error {
	return r.conn(ctx).Model(&model.WorkflowInstanceLevelUser{}).
		Where("tenant_id = ? AND instance_id = ? AND level = ? AND user_id = ?", tenantID, instanceID, level, userID).
		Update("audit_status", status).Error
}

func (r *mysqlRepository) CreateOperate(ctx context.Context, row *model.WorkflowOperate) error {
	return r.conn(ctx).Create(row).Error
}
