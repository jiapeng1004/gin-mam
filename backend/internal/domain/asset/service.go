package asset

import (
	"context"
	"time"
)

const defaultTenant = "default"

const (
	ErrCodeNotFound = 40401
)

const (
	ConfigKeyImageAccess = "MAM_IMAGE_ACCESS_DOMAIN"
	ConfigKeyVideoAccess = "MAM_VIDEO_ACCESS_DOMAIN"
	ConfigKeyVODAccess   = "MAM_VOD_ACCESS_DOMAIN"
	ConfigKeyOtherAccess = "MAM_OTHER_ACCESS_DOMAIN"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Page(ctx context.Context, in PageInput) (*PageResult, error)
	GetByID(ctx context.Context, id string) (*AssetVO, error)
	Create(ctx context.Context, in CreateInput) (*AssetVO, error)
	Update(ctx context.Context, id string, in UpdateInput) (*AssetVO, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int8) error
}

type PageInput struct {
	Page     int
	PageSize int
	Keyword  string
}

type PageResult struct {
	List     []AssetVO `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

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

type UpdateInput struct {
	Title       *string
	CatalogID   *string
	Description *string
	Status      *int8
	Metadata    map[string]string
}
