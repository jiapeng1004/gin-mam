// Package transcode 转码领域：转码组/任务管理与外部回调处理。
package transcode

import (
	"context"
	"time"
)

const defaultTenant = "default"

// ErrCodeNotFound 转码组或任务不存在时的业务错误码。
const ErrCodeNotFound = 40403

// ErrCodeForbidden 转码回调签名校验失败等业务错误码。
const ErrCodeForbidden = 40300

const (
	// ConfigKeyTranscodeURL 外部转码服务提交地址配置键。
	ConfigKeyTranscodeURL = "MAM_TRANSCODE_URL"
	// ConfigKeyTranscodeCallbackURL 转码回调完整 URL 配置键。
	ConfigKeyTranscodeCallbackURL = "MAM_TRANSCODE_CALLBACK_URL"
	// ConfigKeyDomainURL 系统对外域名，用于拼装回调地址。
	ConfigKeyDomainURL = "MAM_DOMAIN_URL"
	// ConfigKeyTranscodeCallbackSecret 转码回调 HMAC 签名密钥配置键。
	ConfigKeyTranscodeCallbackSecret = "MAM_TRANSCODE_CALLBACK_SECRET"
)

const callbackPath = "/api/v1/transcode/callback"

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 转码应用服务接口。
type Service interface {
	// GroupPage 分页查询转码组列表。
	GroupPage(ctx context.Context, page, pageSize int) (*GroupPageResult, error)
	// CreateGroup 创建转码组。
	CreateGroup(ctx context.Context, userID string, in CreateGroupInput) (*GroupVO, error)
	// UpdateGroup 更新转码组可编辑字段。
	UpdateGroup(ctx context.Context, id string, in UpdateGroupInput) (*GroupVO, error)
	// DeleteGroup 软删除转码组。
	DeleteGroup(ctx context.Context, id string) error

	// CreateTask 创建转码任务并提交至外部转码服务。
	CreateTask(ctx context.Context, userID string, in CreateTaskInput) (*TaskVO, error)
	// GetTask 按 ID 获取转码任务详情。
	GetTask(ctx context.Context, id string) (*TaskVO, error)
	// TaskPage 分页查询转码任务列表。
	TaskPage(ctx context.Context, page, pageSize int) (*TaskPageResult, error)
	// HandleCallback 处理外部转码服务 HMAC 签名回调，更新任务状态与输出路径。
	HandleCallback(ctx context.Context, body []byte, signHeader string) error
}

// GroupPageResult 转码组分页结果。
type GroupPageResult struct {
	List     []GroupVO `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

// GroupVO 转码组视图对象。
type GroupVO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Message       string    `json:"message,omitempty"`
	GroupType     int       `json:"groupType"`
	Param         string    `json:"param,omitempty"`
	StrategyType  int       `json:"strategyType"`
	DefaultFlag   int8      `json:"defaultFlag"`
	AvailableFlag int8      `json:"availableFlag"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// CreateGroupInput 创建转码组入参。
type CreateGroupInput struct {
	Name          string
	Message       string
	GroupType     int
	Param         string
	StrategyType  int
	DefaultFlag   int8
	AvailableFlag int8
}

// UpdateGroupInput 更新转码组入参，指针字段 nil 表示不修改。
type UpdateGroupInput struct {
	Name          *string
	Message       *string
	GroupType     *int
	Param         *string
	StrategyType  *int
	DefaultFlag   *int8
	AvailableFlag *int8
}

// TaskPageResult 转码任务分页结果。
type TaskPageResult struct {
	List     []TaskVO `json:"list"`
	Total    int64    `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

// TaskVO 转码任务视图对象。
type TaskVO struct {
	ID               string    `json:"id"`
	AssetID          string    `json:"assetId"`
	AssetFileID      *string   `json:"assetFileId,omitempty"`
	TranscodeGroupID string    `json:"transcodeGroupId"`
	ProfileID        *string   `json:"profileId,omitempty"`
	Status           int8      `json:"status"`
	ExternalJobID    string    `json:"externalJobId,omitempty"`
	OutputPath       string    `json:"outputPath,omitempty"`
	ErrorMsg         string    `json:"errorMsg,omitempty"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// CreateTaskInput 创建转码任务入参。
type CreateTaskInput struct {
	AssetID          string
	AssetFileID      *string
	TranscodeGroupID string
}

// CallbackInput 转码回调 JSON 体结构（供解析参考）。
type CallbackInput struct {
	TaskID     string `json:"taskId"`
	Status     int8   `json:"status"`
	OutputPath string `json:"outputPath"`
	ErrorMsg   string `json:"errorMsg"`
}
