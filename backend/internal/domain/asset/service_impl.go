package asset

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-mam/backend/internal/domain/asset/model"
	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/gin-mam/backend/internal/pkg/id"
	"gorm.io/gorm"
)

type service struct {
	repo   Repository
	config infraconfig.Service
	txMgr  tx.Manager
	now    func() time.Time
}

func NewService(repo Repository, config infraconfig.Service, txMgr tx.Manager) Service {
	return &service{
		repo:   repo,
		config: config,
		txMgr:  txMgr,
		now:    time.Now,
	}
}

func (s *service) notFound() *httpx.BizError {
	return httpx.NewBizError(404, ErrCodeNotFound, "\u5a92\u8d44\u4e0d\u5b58\u5728")
}

func previewConfigKey(assetType string) string {
	switch strings.ToLower(strings.TrimSpace(assetType)) {
	case "image":
		return ConfigKeyImageAccess
	case "video":
		return ConfigKeyVideoAccess
	case "vod":
		return ConfigKeyVODAccess
	default:
		return ConfigKeyOtherAccess
	}
}

func joinPreviewURL(domain, storagePath string) string {
	domain = strings.TrimSpace(domain)
	storagePath = strings.TrimSpace(storagePath)
	if domain == "" || storagePath == "" {
		return ""
	}
	return strings.TrimRight(domain, "/") + "/" + strings.TrimLeft(storagePath, "/")
}

func (s *service) loadPreviewDomain(ctx context.Context, assetType string) (string, error) {
	key := previewConfigKey(assetType)
	var domain string
	err := s.txMgr.Run(ctx, tx.NotSupported, func(ctx context.Context) error {
		val, err := s.config.Get(ctx, key)
		if errors.Is(err, infraconfig.ErrNotFound) {
			domain = ""
			return nil
		}
		if err != nil {
			return err
		}
		domain = val
		return nil
	})
	return domain, err
}

func metadataToMap(rows []model.AssetMetadata) map[string]string {
	if len(rows) == 0 {
		return nil
	}
	m := make(map[string]string, len(rows))
	for _, row := range rows {
		m[row.MetaKey] = row.MetaValue
	}
	return m
}

func (s *service) buildVO(ctx context.Context, asset *model.Asset, master *model.AssetFile, meta []model.AssetMetadata) (*AssetVO, error) {
	vo := &AssetVO{
		ID:          asset.ID,
		Title:       asset.Title,
		Type:        asset.Type,
		Status:      asset.Status,
		CatalogID:   asset.CatalogID,
		Description: asset.Description,
		Metadata:    metadataToMap(meta),
		CreatedBy:   asset.CreatedBy,
		CreatedAt:   asset.CreatedAt,
		UpdatedAt:   asset.UpdatedAt,
	}
	if master != nil {
		vo.StoragePath = master.StoragePath
		domain, err := s.loadPreviewDomain(ctx, asset.Type)
		if err != nil {
			return nil, err
		}
		vo.PreviewURL = joinPreviewURL(domain, master.StoragePath)
	}
	return vo, nil
}

func (s *service) Page(ctx context.Context, in PageInput) (*PageResult, error) {
	page, pageSize := in.Page, in.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	assets, total, err := s.repo.ListAssets(ctx, defaultTenant, page, pageSize, in.Keyword)
	if err != nil {
		return nil, err
	}
	list := make([]AssetVO, 0, len(assets))
	for i := range assets {
		master, err := s.repo.GetMasterFile(ctx, defaultTenant, assets[i].ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		var masterPtr *model.AssetFile
		if err == nil {
			masterPtr = master
		}
		meta, err := s.repo.ListMetadata(ctx, defaultTenant, assets[i].ID)
		if err != nil {
			return nil, err
		}
		vo, err := s.buildVO(ctx, &assets[i], masterPtr, meta)
		if err != nil {
			return nil, err
		}
		list = append(list, *vo)
	}
	return &PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*AssetVO, error) {
	asset, err := s.repo.GetAssetByID(ctx, defaultTenant, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, s.notFound()
		}
		return nil, err
	}
	master, err := s.repo.GetMasterFile(ctx, defaultTenant, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var masterPtr *model.AssetFile
	if err == nil {
		masterPtr = master
	}
	meta, err := s.repo.ListMetadata(ctx, defaultTenant, id)
	if err != nil {
		return nil, err
	}
	return s.buildVO(ctx, asset, masterPtr, meta)
}

func (s *service) Create(ctx context.Context, in CreateInput) (*AssetVO, error) {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Type) == "" || strings.TrimSpace(in.StoragePath) == "" {
		return nil, httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	createdBy := strings.TrimSpace(in.CreatedBy)
	if createdBy == "" {
		createdBy = "system"
	}
	now := s.now()
	assetID := id.New()
	var created *model.Asset
	err := s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		created = &model.Asset{
			ID:          assetID,
			TenantID:    defaultTenant,
			Title:       strings.TrimSpace(in.Title),
			Type:        strings.TrimSpace(in.Type),
			Status:      0,
			CatalogID:   in.CatalogID,
			Description: strings.TrimSpace(in.Description),
			CreatedBy:   createdBy,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := s.repo.CreateAsset(ctx, created); err != nil {
			return err
		}
		file := &model.AssetFile{
			ID:          id.New(),
			TenantID:    defaultTenant,
			AssetID:     assetID,
			Version:     1,
			StoragePath: strings.TrimSpace(in.StoragePath),
			MimeType:    strings.TrimSpace(in.MimeType),
			FileSize:    in.FileSize,
			IsMaster:    1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := s.repo.CreateAssetFile(ctx, file); err != nil {
			return err
		}
		for k, v := range in.Metadata {
			key := strings.TrimSpace(k)
			if key == "" {
				continue
			}
			if err := s.repo.UpsertMetadata(ctx, &model.AssetMetadata{
				ID:        id.New(),
				TenantID:  defaultTenant,
				AssetID:   assetID,
				MetaKey:   key,
				MetaValue: v,
				CreatedAt: now,
				UpdatedAt: now,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, created.ID)
}

func (s *service) Update(ctx context.Context, assetID string, in UpdateInput) (*AssetVO, error) {
	err := s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		asset, err := s.repo.GetAssetByID(ctx, defaultTenant, assetID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return s.notFound()
			}
			return err
		}
		if in.Title != nil {
			asset.Title = strings.TrimSpace(*in.Title)
		}
		if in.CatalogID != nil {
			asset.CatalogID = in.CatalogID
		}
		if in.Description != nil {
			asset.Description = strings.TrimSpace(*in.Description)
		}
		if in.Status != nil {
			asset.Status = *in.Status
		}
		asset.UpdatedAt = s.now()
		if err := s.repo.UpdateAsset(ctx, asset); err != nil {
			return err
		}
		if in.Metadata != nil {
			now := s.now()
			keys := make([]string, 0, len(in.Metadata))
			for k, v := range in.Metadata {
				key := strings.TrimSpace(k)
				if key == "" {
					continue
				}
				keys = append(keys, key)
				if err := s.repo.UpsertMetadata(ctx, &model.AssetMetadata{
					ID:        id.New(),
					TenantID:  defaultTenant,
					AssetID:   assetID,
					MetaKey:   key,
					MetaValue: v,
					CreatedAt: now,
					UpdatedAt: now,
				}); err != nil {
					return err
				}
			}
			if err := s.repo.DeleteMetadataNotInKeys(ctx, defaultTenant, assetID, keys); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			return nil, biz
		}
		return nil, err
	}
	return s.GetByID(ctx, assetID)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.txMgr.Run(ctx, tx.Required, func(ctx context.Context) error {
		if err := s.repo.DeleteAsset(ctx, defaultTenant, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return s.notFound()
			}
			return err
		}
		return nil
	})
}