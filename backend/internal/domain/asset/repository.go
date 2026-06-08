package asset

import (
	"context"

	"github.com/gin-mam/backend/internal/domain/asset/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=repository.go -destination=mock/repository_mock.go -package=mock
type Repository interface {
	GetAssetByID(ctx context.Context, tenantID, id string) (*model.Asset, error)
	ListAssets(ctx context.Context, tenantID string, page, pageSize int, keyword string) ([]model.Asset, int64, error)
	CreateAsset(ctx context.Context, asset *model.Asset) error
	UpdateAsset(ctx context.Context, asset *model.Asset) error
	DeleteAsset(ctx context.Context, tenantID, id string) error

	GetMasterFile(ctx context.Context, tenantID, assetID string) (*model.AssetFile, error)
	CreateAssetFile(ctx context.Context, file *model.AssetFile) error

	ListMetadata(ctx context.Context, tenantID, assetID string) ([]model.AssetMetadata, error)
	UpsertMetadata(ctx context.Context, row *model.AssetMetadata) error
	DeleteMetadataNotInKeys(ctx context.Context, tenantID, assetID string, keys []string) error
}