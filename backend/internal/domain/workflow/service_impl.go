package workflow

import (
	"context"
	"errors"
	"time"

	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/workflow/model"
	"github.com/gin-mam/backend/internal/infra/lock"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/gin-mam/backend/internal/pkg/id"
	"gorm.io/gorm"
)

const lockTTL = 30 * time.Second

type service struct {
	repo      Repository
	assetRepo asset.Repository
	assetSvc  asset.Service
	txMgr     tx.Manager
	locker    lock.Locker
	now       func() time.Time
}

func NewService(repo Repository, assetRepo asset.Repository, assetSvc asset.Service, txMgr tx.Manager, locker lock.Locker) Service {
	return &service{
		repo:      repo,
		assetRepo: assetRepo,
		assetSvc:  assetSvc,
		txMgr:     txMgr,
		locker:    locker,
		now:       time.Now,
	}
}

func (s *service) notFound(msg string) *httpx.BizError {
	if msg == "" {
		msg = "工作流不存在"
	}
	return httpx.NewBizError(404, ErrCodeNotFound, msg)
}

func (s *service) lockConflict() *httpx.BizError {
	return httpx.NewBizError(409, ErrCodeLockConflict, "工作流正在处理中")
}

func (s *service) invalidState(msg string) *httpx.BizError {
	return httpx.NewBizError(422, ErrCodeInvalidState, msg)
}

func instanceToVO(row *model.WorkflowInstance) *InstanceVO {
	return &InstanceVO{
		ID:            row.ID,
		WorkflowDefID: row.WorkflowDefID,
		AssetID:       row.AssetID,
		Level:         row.Level,
		AuditLevel:    row.AuditLevel,
		AuditStatus:   row.AuditStatus,
		CreatedBy:     row.CreatedBy,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func (s *service) writeOperate(ctx context.Context, instanceID, operatorID string, action int8, remark string) error {
	return s.txMgr.Run(ctx, tx.RequiresNew, func(ctx context.Context) error {
		now := s.now()
		return s.repo.CreateOperate(ctx, &model.WorkflowOperate{
			ID:         id.New(),
			TenantID:   defaultTenant,
			InstanceID: instanceID,
			OperatorID: operatorID,
			Action:     action,
			Remark:     remark,
			CreatedAt:  now,
		})
	})
}

func (s *service) Submit(ctx context.Context, userID string, in SubmitRequest) (*InstanceVO, error) {
	if userID == "" || in.AssetID == "" || in.WorkflowDefID == "" {
		return nil, httpx.NewBizError(400, 40000, "参数错误")
	}
	var created *model.WorkflowInstance
	err := s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		assetRow, err := s.assetRepo.GetAssetByID(ctx, defaultTenant, in.AssetID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return httpx.NewBizError(404, asset.ErrCodeNotFound, "媒资不存在")
			}
			return err
		}
		if assetRow.Status != AssetStatusDraft && assetRow.Status != AssetStatusRejected {
			return s.invalidState("媒资状态不允许送审")
		}
		def, err := s.repo.GetDefByID(ctx, defaultTenant, in.WorkflowDefID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return s.notFound("工作流定义不存在")
			}
			return err
		}
		levelUsers, err := s.repo.ListLevelUsersByDefID(ctx, defaultTenant, def.ID)
		if err != nil {
			return err
		}
		now := s.now()
		inst := &model.WorkflowInstance{
			ID:            id.New(),
			TenantID:      defaultTenant,
			WorkflowDefID: def.ID,
			AssetID:       in.AssetID,
			Level:         0,
			AuditLevel:    def.AuditLevel,
			AuditStatus:   AuditStatusPending,
			CreatedBy:     userID,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.repo.CreateInstance(ctx, inst); err != nil {
			return err
		}
		ilu := make([]model.WorkflowInstanceLevelUser, 0, len(levelUsers))
		for _, lu := range levelUsers {
			ilu = append(ilu, model.WorkflowInstanceLevelUser{
				ID:          id.New(),
				TenantID:    defaultTenant,
				InstanceID:  inst.ID,
				Level:       lu.Level,
				UserID:      lu.UserID,
				AuditStatus: AuditStatusPending,
			})
		}
		if err := s.repo.CreateInstanceLevelUsers(ctx, ilu); err != nil {
			return err
		}
		if err := s.assetSvc.UpdateStatus(ctx, in.AssetID, AssetStatusAuditing); err != nil {
			return err
		}
		created = inst
		return nil
	})
	if err != nil {
		return nil, err
	}
	_ = s.writeOperate(ctx, created.ID, userID, OperateActionSubmit, "")
	return instanceToVO(created), nil
}

func (s *service) Audit(ctx context.Context, userID string, in AuditRequest) error {
	if userID == "" || in.InstanceID == "" {
		return httpx.NewBizError(400, 40000, "参数错误")
	}
	lockKey := lock.InstanceLockKey(in.InstanceID)
	token, ok, err := s.locker.TryLock(ctx, lockKey, lockTTL)
	if err != nil {
		return err
	}
	if !ok {
		return s.lockConflict()
	}
	defer func() { _ = s.locker.Unlock(context.Background(), lockKey, token) }()

	var action int8
	if in.Pass {
		action = OperateActionPass
	} else {
		action = OperateActionReject
	}

	err = s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		inst, err := s.repo.GetInstanceByID(ctx, defaultTenant, in.InstanceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return s.notFound("")
			}
			return err
		}
		if inst.AuditStatus != AuditStatusPending && inst.AuditStatus != AuditStatusAuditing {
			return s.invalidState("当前状态不可审核")
		}
		assigneeLevel := AssigneeLevel(inst.Level)
		isAssignee, err := s.repo.IsUserAssigneeAtLevel(ctx, defaultTenant, inst.ID, assigneeLevel, userID)
		if err != nil {
			return err
		}
		if !isAssignee {
			return httpx.NewBizError(403, 40300, "无权审核")
		}
		locker := userID
		inst.LockerID = &locker
		inst.UpdatedAt = s.now()
		if err := s.repo.UpdateInstance(ctx, inst); err != nil {
			return err
		}

		newLevel, newStatus := ApplyAuditAction(inst.Level, inst.AuditLevel, in.Pass)
		inst.Level = newLevel
		inst.AuditStatus = newStatus
		inst.LockerID = nil
		inst.UpdatedAt = s.now()
		if err := s.repo.UpdateInstance(ctx, inst); err != nil {
			return err
		}
		if err := s.repo.UpdateInstanceLevelUserStatus(ctx, defaultTenant, inst.ID, assigneeLevel, userID, newStatus); err != nil {
			return err
		}

		var assetStatus int8
		switch newStatus {
		case AuditStatusPassed:
			assetStatus = AssetStatusPassed
		case AuditStatusRejected:
			assetStatus = AssetStatusRejected
		default:
			assetStatus = AssetStatusAuditing
		}
		if newStatus == AuditStatusPassed || newStatus == AuditStatusRejected || newStatus == AuditStatusAuditing {
			if err := s.assetSvc.UpdateStatus(ctx, inst.AssetID, assetStatus); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return s.writeOperate(ctx, in.InstanceID, userID, action, in.Remark)
}

func (s *service) MultiAudit(ctx context.Context, userID string, items []AuditRequest) error {
	for _, item := range items {
		if err := s.Audit(ctx, userID, item); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Revoke(ctx context.Context, userID string, instanceIDs []string) error {
	if userID == "" || len(instanceIDs) == 0 {
		return httpx.NewBizError(400, 40000, "参数错误")
	}
	for _, instanceID := range instanceIDs {
		err := s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
			inst, err := s.repo.GetInstanceByID(ctx, defaultTenant, instanceID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return s.notFound("")
				}
				return err
			}
			if inst.CreatedBy != userID {
				return httpx.NewBizError(403, 40300, "无权撤回")
			}
			if inst.AuditStatus == AuditStatusPassed {
				return s.invalidState("已通过不可撤回")
			}
			if err := s.assetSvc.UpdateStatus(ctx, inst.AssetID, AssetStatusDraft); err != nil {
				return err
			}
			return s.repo.DeleteInstance(ctx, defaultTenant, instanceID)
		})
		if err != nil {
			return err
		}
		_ = s.writeOperate(ctx, instanceID, userID, OperateActionRevoke, "")
	}
	return nil
}

func (s *service) pageResult(rows []model.WorkflowInstance, total int64, page, pageSize int) *PageResult {
	list := make([]InstanceVO, 0, len(rows))
	for i := range rows {
		list = append(list, *instanceToVO(&rows[i]))
	}
	return &PageResult{List: list, Total: total, Page: page, PageSize: pageSize}
}

func (s *service) ListAssignedToMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListInstancesAssignedToUser(ctx, defaultTenant, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return s.pageResult(rows, total, page, pageSize), nil
}

func (s *service) ListCreatedByMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListInstancesCreatedBy(ctx, defaultTenant, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return s.pageResult(rows, total, page, pageSize), nil
}

func (s *service) ListAuditedByMe(ctx context.Context, userID string, page, pageSize int) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListInstancesAuditedBy(ctx, defaultTenant, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return s.pageResult(rows, total, page, pageSize), nil
}

func (s *service) ListAll(ctx context.Context, page, pageSize int) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListAllInstances(ctx, defaultTenant, page, pageSize)
	if err != nil {
		return nil, err
	}
	return s.pageResult(rows, total, page, pageSize), nil
}

func (s *service) CreateDef(ctx context.Context, in CreateDefInput) (*DefVO, error) {
	if in.Name == "" || in.AuditLevel < 1 {
		return nil, httpx.NewBizError(400, 40000, "参数错误")
	}
	now := s.now()
	defID := id.New()
	def := &model.WorkflowDef{
		ID:         defID,
		TenantID:   defaultTenant,
		Name:       in.Name,
		AuditLevel: in.AuditLevel,
		Status:     1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	err := s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		if err := s.repo.CreateDef(ctx, def); err != nil {
			return err
		}
		rows := make([]model.WorkflowLevelUser, 0, len(in.LevelUsers))
		for _, lu := range in.LevelUsers {
			rows = append(rows, model.WorkflowLevelUser{
				ID:            id.New(),
				TenantID:      defaultTenant,
				WorkflowDefID: defID,
				Level:         lu.Level,
				UserID:        lu.UserID,
			})
		}
		return s.repo.CreateLevelUsers(ctx, rows)
	})
	if err != nil {
		return nil, err
	}
	return s.defToVO(def, in.LevelUsers), nil
}

func (s *service) DefPage(ctx context.Context, page, pageSize int) (*DefPageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListDefs(ctx, defaultTenant, page, pageSize)
	if err != nil {
		return nil, err
	}
	list := make([]DefVO, 0, len(rows))
	for i := range rows {
		lu, err := s.repo.ListLevelUsersByDefID(ctx, defaultTenant, rows[i].ID)
		if err != nil {
			return nil, err
		}
		voUsers := make([]DefLevelUserVO, 0, len(lu))
		for _, u := range lu {
			voUsers = append(voUsers, DefLevelUserVO{Level: u.Level, UserID: u.UserID})
		}
		list = append(list, *s.defToVO(&rows[i], voUsers))
	}
	return &DefPageResult{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *service) defToVO(def *model.WorkflowDef, levelUsers []DefLevelUserVO) *DefVO {
	return &DefVO{
		ID:         def.ID,
		Name:       def.Name,
		AuditLevel: def.AuditLevel,
		Status:     def.Status,
		LevelUsers: levelUsers,
		CreatedAt:  def.CreatedAt,
		UpdatedAt:  def.UpdatedAt,
	}
}
