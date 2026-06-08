package transcode

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gin-mam/backend/internal/domain/transcode/client"
	"github.com/gin-mam/backend/internal/domain/transcode/model"
	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/gin-mam/backend/internal/pkg/id"
	"gorm.io/gorm"
)

type service struct {
	repo      Repository
	configSvc infraconfig.Service
	txMgr     tx.Manager
	client    client.Client
	now       func() time.Time
}

func NewService(repo Repository, configSvc infraconfig.Service, txMgr tx.Manager, c client.Client) Service {
	return &service{
		repo:      repo,
		configSvc: configSvc,
		txMgr:     txMgr,
		client:    c,
		now:       time.Now,
	}
}

func (s *service) notFoundGroup() *httpx.BizError {
	return httpx.NewBizError(404, ErrCodeNotFound, "转码组不存在")
}

func (s *service) notFoundTask() *httpx.BizError {
	return httpx.NewBizError(404, ErrCodeNotFound, "转码任务不存在")
}

func groupToVO(row *model.TranscodeGroup) *GroupVO {
	return &GroupVO{
		ID:            row.ID,
		Name:          row.Name,
		Message:       row.Message,
		GroupType:     row.GroupType,
		Param:         row.Param,
		StrategyType:  row.StrategyType,
		DefaultFlag:   row.DefaultFlag,
		AvailableFlag: row.AvailableFlag,
		CreatedBy:     row.CreatedBy,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func taskToVO(row *model.TranscodeTask) *TaskVO {
	return &TaskVO{
		ID:               row.ID,
		AssetID:          row.AssetID,
		AssetFileID:      row.AssetFileID,
		TranscodeGroupID: row.TranscodeGroupID,
		ProfileID:        row.ProfileID,
		Status:           row.Status,
		ExternalJobID:    row.ExternalJobID,
		OutputPath:       row.OutputPath,
		ErrorMsg:         row.ErrorMsg,
		CreatedBy:        row.CreatedBy,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func (s *service) GroupPage(ctx context.Context, page, pageSize int) (*GroupPageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListGroups(ctx, defaultTenant, page, pageSize)
	if err != nil {
		return nil, err
	}
	list := make([]GroupVO, 0, len(rows))
	for i := range rows {
		list = append(list, *groupToVO(&rows[i]))
	}
	return &GroupPageResult{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *service) CreateGroup(ctx context.Context, userID string, in CreateGroupInput) (*GroupVO, error) {
	if strings.TrimSpace(in.Name) == "" || userID == "" {
		return nil, httpx.NewBizError(400, 40000, "参数错误")
	}
	now := s.now()
	row := &model.TranscodeGroup{
		ID:            id.New(),
		TenantID:      defaultTenant,
		Name:          in.Name,
		Message:       in.Message,
		GroupType:     in.GroupType,
		Param:         in.Param,
		StrategyType:  in.StrategyType,
		DefaultFlag:   in.DefaultFlag,
		AvailableFlag: in.AvailableFlag,
		CreatedBy:     userID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	row.AvailableFlag = in.AvailableFlag
	if row.AvailableFlag == 0 {
		row.AvailableFlag = 1
	}
	if err := s.repo.CreateGroup(ctx, row); err != nil {
		return nil, err
	}
	return groupToVO(row), nil
}

func (s *service) UpdateGroup(ctx context.Context, id string, in UpdateGroupInput) (*GroupVO, error) {
	row, err := s.repo.GetGroupByID(ctx, defaultTenant, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, s.notFoundGroup()
		}
		return nil, err
	}
	if in.Name != nil {
		row.Name = *in.Name
	}
	if in.Message != nil {
		row.Message = *in.Message
	}
	if in.GroupType != nil {
		row.GroupType = *in.GroupType
	}
	if in.Param != nil {
		row.Param = *in.Param
	}
	if in.StrategyType != nil {
		row.StrategyType = *in.StrategyType
	}
	if in.DefaultFlag != nil {
		row.DefaultFlag = *in.DefaultFlag
	}
	if in.AvailableFlag != nil {
		row.AvailableFlag = *in.AvailableFlag
	}
	row.UpdatedAt = s.now()
	if err := s.repo.UpdateGroup(ctx, row); err != nil {
		return nil, err
	}
	return groupToVO(row), nil
}

func (s *service) DeleteGroup(ctx context.Context, id string) error {
	err := s.repo.DeleteGroup(ctx, defaultTenant, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.notFoundGroup()
		}
		return err
	}
	return nil
}

func (s *service) resolveCallbackURL(ctx context.Context) (string, error) {
	if v, err := s.configSvc.Get(ctx, ConfigKeyTranscodeCallbackURL); err == nil && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v), nil
	} else if err != nil && !errors.Is(err, infraconfig.ErrNotFound) {
		return "", err
	}
	domain, err := s.configSvc.Get(ctx, ConfigKeyDomainURL)
	if err != nil {
		if errors.Is(err, infraconfig.ErrNotFound) {
			return "", httpx.NewBizError(500, 50000, "未配置转码回调地址")
		}
		return "", err
	}
	domain = strings.TrimRight(strings.TrimSpace(domain), "/")
	if domain == "" {
		return "", httpx.NewBizError(500, 50000, "未配置转码回调地址")
	}
	return domain + callbackPath, nil
}

func (s *service) CreateTask(ctx context.Context, userID string, in CreateTaskInput) (*TaskVO, error) {
	if userID == "" || strings.TrimSpace(in.AssetID) == "" || strings.TrimSpace(in.TranscodeGroupID) == "" {
		return nil, httpx.NewBizError(400, 40000, "参数错误")
	}
	group, err := s.repo.GetGroupByID(ctx, defaultTenant, in.TranscodeGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, s.notFoundGroup()
		}
		return nil, err
	}
	if group.AvailableFlag != 1 {
		return nil, httpx.NewBizError(400, 40000, "转码组不可用")
	}
	transcodeURL, err := s.configSvc.Get(ctx, ConfigKeyTranscodeURL)
	if err != nil {
		if errors.Is(err, infraconfig.ErrNotFound) {
			return nil, httpx.NewBizError(500, 50000, "未配置转码服务地址")
		}
		return nil, err
	}
	transcodeURL = strings.TrimSpace(transcodeURL)
	if transcodeURL == "" {
		return nil, httpx.NewBizError(500, 50000, "未配置转码服务地址")
	}
	callbackURL, err := s.resolveCallbackURL(ctx)
	if err != nil {
		return nil, err
	}

	now := s.now()
	taskID := id.New()
	task := &model.TranscodeTask{
		ID:               taskID,
		TenantID:         defaultTenant,
		AssetID:          in.AssetID,
		AssetFileID:      in.AssetFileID,
		TranscodeGroupID: in.TranscodeGroupID,
		Status:           model.TaskStatusPending,
		CreatedBy:        userID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	err = s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		return s.repo.CreateTask(ctx, task)
	})
	if err != nil {
		return nil, err
	}

	assetFileID := ""
	if in.AssetFileID != nil {
		assetFileID = *in.AssetFileID
	}
	resp, submitErr := s.client.SubmitJob(ctx, transcodeURL, client.SubmitJobRequest{
		TaskID:           taskID,
		CallbackURL:      callbackURL,
		AssetID:          in.AssetID,
		AssetFileID:      assetFileID,
		TranscodeGroupID: in.TranscodeGroupID,
		GroupParam:       group.Param,
	})
	task.UpdatedAt = s.now()
	if submitErr != nil {
		task.Status = model.TaskStatusFailed
		task.ErrorMsg = submitErr.Error()
		_ = s.repo.UpdateTask(ctx, task)
		return nil, submitErr
	}
	task.Status = model.TaskStatusRunning
	if resp != nil {
		task.ExternalJobID = resp.JobID
	}
	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}
	return taskToVO(task), nil
}

func (s *service) GetTask(ctx context.Context, id string) (*TaskVO, error) {
	row, err := s.repo.GetTaskByID(ctx, defaultTenant, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, s.notFoundTask()
		}
		return nil, err
	}
	return taskToVO(row), nil
}

func (s *service) TaskPage(ctx context.Context, page, pageSize int) (*TaskPageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, total, err := s.repo.ListTasks(ctx, defaultTenant, page, pageSize)
	if err != nil {
		return nil, err
	}
	list := make([]TaskVO, 0, len(rows))
	for i := range rows {
		list = append(list, *taskToVO(&rows[i]))
	}
	return &TaskPageResult{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *service) HandleCallback(ctx context.Context, body []byte, signHeader string) error {
	secret, err := s.configSvc.Get(ctx, ConfigKeyTranscodeCallbackSecret)
	if err != nil {
		if errors.Is(err, infraconfig.ErrNotFound) {
			return httpx.NewBizError(403, ErrCodeForbidden, "签名校验失败")
		}
		return err
	}
	if !VerifyCallbackSign(body, secret, signHeader) {
		return httpx.NewBizError(403, ErrCodeForbidden, "签名校验失败")
	}
	var in CallbackInput
	if err := json.Unmarshal(body, &in); err != nil {
		return httpx.NewBizError(400, 40000, "参数错误")
	}
	if strings.TrimSpace(in.TaskID) == "" {
		return httpx.NewBizError(400, 40000, "参数错误")
	}
	task, err := s.repo.GetTaskByID(ctx, defaultTenant, in.TaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.notFoundTask()
		}
		return err
	}
	payload := string(body)
	task.CallbackPayload = &payload
	task.Status = in.Status
	task.OutputPath = in.OutputPath
	task.ErrorMsg = in.ErrorMsg
	task.UpdatedAt = s.now()
	return s.repo.UpdateTask(ctx, task)
}