package catalog

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-mam/backend/internal/domain/catalog/model"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/gin-mam/backend/internal/pkg/id"
	"gorm.io/gorm"
)

type service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return &service{repo: repo, now: time.Now}
}

func (s *service) notFound() *httpx.BizError {
	return httpx.NewBizError(404, ErrCodeNotFound, "\u680f\u76ee\u4e0d\u5b58\u5728")
}

func catalogToVO(row *model.Catalog) *CatalogVO {
	return &CatalogVO{
		ID:       row.ID,
		Name:     row.Name,
		ParentID: row.ParentID,
		SortCode: row.SortCode,
	}
}

func configToVO(row *model.CatalogConfig) ConfigVO {
	return ConfigVO{
		ID:          row.ID,
		ConfigKey:   row.ConfigKey,
		ConfigValue: row.ConfigValue,
	}
}

func buildTree(rows []model.Catalog) []TreeNode {
	index := make(map[string]TreeNode, len(rows))
	for _, row := range rows {
		index[row.ID] = TreeNode{
			ID:       row.ID,
			Name:     row.Name,
			ParentID: row.ParentID,
			Children: []TreeNode{},
		}
	}
	rootIDs := make([]string, 0)
	childLinks := make(map[string][]string)
	for _, row := range rows {
		if row.ParentID == nil || strings.TrimSpace(*row.ParentID) == "" {
			rootIDs = append(rootIDs, row.ID)
			continue
		}
		pid := strings.TrimSpace(*row.ParentID)
		if _, ok := index[pid]; !ok {
			rootIDs = append(rootIDs, row.ID)
			continue
		}
		childLinks[pid] = append(childLinks[pid], row.ID)
	}
	var build func(id string) TreeNode
	build = func(id string) TreeNode {
		n := index[id]
		for _, cid := range childLinks[id] {
			n.Children = append(n.Children, build(cid))
		}
		return n
	}
	roots := make([]TreeNode, 0, len(rootIDs))
	for _, id := range rootIDs {
		roots = append(roots, build(id))
	}
	return roots
}

func (s *service) Tree(ctx context.Context) ([]TreeNode, error) {
	rows, err := s.repo.ListCatalogs(ctx, defaultTenant)
	if err != nil {
		return nil, err
	}
	return buildTree(rows), nil
}

func (s *service) validateParent(ctx context.Context, parentID *string, selfID string) error {
	if parentID == nil || strings.TrimSpace(*parentID) == "" {
		return nil
	}
	pid := strings.TrimSpace(*parentID)
	if selfID != "" && pid == selfID {
		return httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	_, err := s.repo.GetCatalogByID(ctx, defaultTenant, pid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.notFound()
		}
		return err
	}
	return nil
}

func (s *service) Create(ctx context.Context, in CreateInput) (*CatalogVO, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	if err := s.validateParent(ctx, in.ParentID, ""); err != nil {
		return nil, err
	}
	now := s.now()
	row := &model.Catalog{
		ID:        id.New(),
		TenantID:  defaultTenant,
		ParentID:  in.ParentID,
		Name:      name,
		SortCode:  in.SortCode,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.CreateCatalog(ctx, row); err != nil {
		return nil, err
	}
	return catalogToVO(row), nil
}

func (s *service) Update(ctx context.Context, id string, in UpdateInput) (*CatalogVO, error) {
	row, err := s.repo.GetCatalogByID(ctx, defaultTenant, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, s.notFound()
		}
		return nil, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
		}
		row.Name = name
	}
	if in.ParentID != nil {
		if err := s.validateParent(ctx, in.ParentID, id); err != nil {
			return nil, err
		}
		row.ParentID = in.ParentID
	}
	if in.SortCode != nil {
		row.SortCode = *in.SortCode
	}
	row.UpdatedAt = s.now()
	if err := s.repo.UpdateCatalog(ctx, row); err != nil {
		return nil, err
	}
	return catalogToVO(row), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.DeleteCatalog(ctx, defaultTenant, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.notFound()
		}
		return err
	}
	return nil
}

func (s *service) ensureCatalog(ctx context.Context, catalogID string) error {
	_, err := s.repo.GetCatalogByID(ctx, defaultTenant, catalogID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.notFound()
		}
		return err
	}
	return nil
}

func (s *service) ListConfig(ctx context.Context, catalogID string) ([]ConfigVO, error) {
	if err := s.ensureCatalog(ctx, catalogID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListConfigs(ctx, defaultTenant, catalogID)
	if err != nil {
		return nil, err
	}
	out := make([]ConfigVO, 0, len(rows))
	for i := range rows {
		out = append(out, configToVO(&rows[i]))
	}
	return out, nil
}

func (s *service) BatchUpdateConfig(ctx context.Context, catalogID string, items []ConfigItemInput) ([]ConfigVO, error) {
	if err := s.ensureCatalog(ctx, catalogID); err != nil {
		return nil, err
	}
	now := s.now()
	keys := make([]string, 0, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item.ConfigKey)
		if key == "" {
			continue
		}
		keys = append(keys, key)
		if err := s.repo.UpsertConfig(ctx, &model.CatalogConfig{
			ID:          id.New(),
			TenantID:    defaultTenant,
			CatalogID:   catalogID,
			ConfigKey:   key,
			ConfigValue: item.ConfigValue,
			CreatedAt:   now,
			UpdatedAt:   now,
		}); err != nil {
			return nil, err
		}
	}
	if err := s.repo.DeleteConfigsNotInKeys(ctx, defaultTenant, catalogID, keys); err != nil {
		return nil, err
	}
	return s.ListConfig(ctx, catalogID)
}