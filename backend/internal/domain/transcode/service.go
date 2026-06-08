package transcode

import (
	"context"
	"time"
)

const defaultTenant = "default"

const ErrCodeNotFound = 40403
const ErrCodeForbidden = 40300

const (
	ConfigKeyTranscodeURL           = "MAM_TRANSCODE_URL"
	ConfigKeyTranscodeCallbackURL   = "MAM_TRANSCODE_CALLBACK_URL"
	ConfigKeyDomainURL              = "MAM_DOMAIN_URL"
	ConfigKeyTranscodeCallbackSecret = "MAM_TRANSCODE_CALLBACK_SECRET"
)

const callbackPath = "/api/v1/transcode/callback"

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	GroupPage(ctx context.Context, page, pageSize int) (*GroupPageResult, error)
	CreateGroup(ctx context.Context, userID string, in CreateGroupInput) (*GroupVO, error)
	UpdateGroup(ctx context.Context, id string, in UpdateGroupInput) (*GroupVO, error)
	DeleteGroup(ctx context.Context, id string) error

	CreateTask(ctx context.Context, userID string, in CreateTaskInput) (*TaskVO, error)
	GetTask(ctx context.Context, id string) (*TaskVO, error)
	TaskPage(ctx context.Context, page, pageSize int) (*TaskPageResult, error)
	HandleCallback(ctx context.Context, body []byte, signHeader string) error
}

type GroupPageResult struct {
	List     []GroupVO `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

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

type CreateGroupInput struct {
	Name          string
	Message       string
	GroupType     int
	Param         string
	StrategyType  int
	DefaultFlag   int8
	AvailableFlag int8
}

type UpdateGroupInput struct {
	Name          *string
	Message       *string
	GroupType     *int
	Param         *string
	StrategyType  *int
	DefaultFlag   *int8
	AvailableFlag *int8
}

type TaskPageResult struct {
	List     []TaskVO `json:"list"`
	Total    int64    `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

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

type CreateTaskInput struct {
	AssetID          string
	AssetFileID      *string
	TranscodeGroupID string
}

type CallbackInput struct {
	TaskID     string `json:"taskId"`
	Status     int8   `json:"status"`
	OutputPath string `json:"outputPath"`
	ErrorMsg   string `json:"errorMsg"`
}