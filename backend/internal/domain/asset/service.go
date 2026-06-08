// Package asset 媒资领域：CRUD、预览 URL 拼装、分片上传与检索索引。
package asset

import (
	"context"
	"time"
)

const defaultTenant = "default"

const (
	// ErrCodeNotFound 媒资不存在时的业务错误码。
	ErrCodeNotFound = 40401
)

const (
	// ConfigKeyImageAccess 图片预览域名配置键。
	ConfigKeyImageAccess = "MAM_IMAGE_ACCESS_DOMAIN"
	// ConfigKeyVideoAccess 视频预览域名配置键。
	ConfigKeyVideoAccess = "MAM_VIDEO_ACCESS_DOMAIN"
	// ConfigKeyVODAccess 点播/VOD 播放域名配置键。
	ConfigKeyVODAccess = "MAM_VOD_ACCESS_DOMAIN"
	// ConfigKeyOtherAccess 其他类型媒资访问域名配置键。
	ConfigKeyOtherAccess = "MAM_OTHER_ACCESS_DOMAIN"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 媒资应用服务接口。
// 负责媒资生命周期管理，并在查询时通过 ConfigService 拼装 previewUrl。
type Service interface {
	// Page 分页查询媒资列表，支持 keyword 模糊匹配标题。
	Page(ctx context.Context, in PageInput) (*PageResult, error)
	// GetByID 按 ID 获取媒资详情，含主文件 storagePath 与 metadata。
	GetByID(ctx context.Context, id string) (*AssetVO, error)
	// Create 创建媒资及其主文件记录，可选写入 metadata。
	Create(ctx context.Context, in CreateInput) (*AssetVO, error)
	// Update 更新媒资可编辑字段；未传字段保持不变。
	Update(ctx context.Context, id string, in UpdateInput) (*AssetVO, error)
	// Delete 软删除媒资（GORM DeletedAt）。
	Delete(ctx context.Context, id string) error
	// UpdateStatus 更新媒资业务状态，供工作流等领域联动调用。
	UpdateStatus(ctx context.Context, id string, status int8) error
}

// PageInput 媒资分页查询入参。
type PageInput struct {
	// Page 页码，从 1 开始。
	Page int
	// PageSize 每页条数。
	PageSize int
	// Keyword 标题关键字，空则不过滤。
	Keyword string
}

// PageResult 媒资分页结果，JSON 字段 camelCase。
type PageResult struct {
	List     []AssetVO `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

// AssetVO 媒资对外视图对象。
type AssetVO struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Status      int8              `json:"status"`
	CatalogID   *string           `json:"catalogId,omitempty"`
	Description string            `json:"description,omitempty"`
	PreviewURL  string            `json:"previewUrl,omitempty"`
	StoragePath string            `json:"storagePath,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedBy   string            `json:"createdBy"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// CreateInput 创建媒资入参。
type CreateInput struct {
	Title       string
	Type        string
	CatalogID   *string
	Description string
	StoragePath string
	MimeType    string
	FileSize    *int64
	Metadata    map[string]string
	CreatedBy   string
}

// UpdateInput 更新媒资入参，指针字段 nil 表示不修改。
type UpdateInput struct {
	Title       *string
	CatalogID   *string
	Description *string
	Status      *int8
	Metadata    map[string]string
}
