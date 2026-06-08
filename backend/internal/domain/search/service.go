package search

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/infra/elasticsearch"
)

const defaultTenant = "default"

type SearchInput struct {
	Keyword   string
	Page      int
	PageSize  int
	Type      string
	CatalogID string
}

type AssetDocument struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	CatalogID *string   `json:"catalogId,omitempty"`
	TenantID  string    `json:"tenantId"`
	CreatedAt time.Time `json:"createdAt"`
}

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

type Service interface {
	SearchAsset(ctx context.Context, in SearchInput) (*asset.PageResult, error)
	ScheduleIndex(doc asset.IndexDoc)
}

type service struct {
	esClient elasticsearch.Client
	assets   asset.Service
}

func NewService(esClient elasticsearch.Client, assets asset.Service) Service {
	return &service{esClient: esClient, assets: assets}
}

func (s *service) ScheduleIndex(doc asset.IndexDoc) {
	if s == nil || s.esClient == nil {
		return
	}
	go func(d asset.IndexDoc) {
		body, err := MarshalAssetDocument(d)
		if err != nil {
			return
		}
		_ = s.esClient.IndexDocument(context.Background(), elasticsearch.AssetIndexName, d.ID, body)
	}(doc)
}

func MarshalAssetDocument(doc asset.IndexDoc) ([]byte, error) {
	return json.Marshal(AssetDocument{
		ID:        doc.ID,
		Title:     doc.Title,
		Type:      doc.Type,
		CatalogID: doc.CatalogID,
		TenantID:  doc.TenantID,
		CreatedAt: doc.CreatedAt,
	})
}

func (s *service) SearchAsset(ctx context.Context, in SearchInput) (*asset.PageResult, error) {
	page, pageSize := in.Page, in.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	body, err := json.Marshal(BuildAssetSearchQuery(defaultTenant, in, page, pageSize))
	if err != nil {
		return nil, err
	}
	resp, err := s.esClient.Search(ctx, elasticsearch.AssetIndexName, body)
	if err != nil {
		return nil, err
	}
	list := make([]asset.AssetVO, 0, len(resp.Hits))
	for _, hit := range resp.Hits {
		vo, err := s.assets.GetByID(ctx, hit.ID)
		if err != nil {
			continue
		}
		list = append(list, *vo)
	}
	return &asset.PageResult{
		List:     list,
		Total:    resp.Total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func BuildAssetSearchQuery(tenantID string, in SearchInput, page, pageSize int) map[string]any {
	from := (page - 1) * pageSize
	filter := []map[string]any{
		{"term": map[string]any{"tenantId": tenantID}},
	}
	if t := strings.TrimSpace(in.Type); t != "" {
		filter = append(filter, map[string]any{"term": map[string]any{"type": t}})
	}
	if cid := strings.TrimSpace(in.CatalogID); cid != "" {
		filter = append(filter, map[string]any{"term": map[string]any{"catalogId": cid}})
	}
	var query map[string]any
	keyword := strings.TrimSpace(in.Keyword)
	if keyword == "" {
		query = map[string]any{"bool": map[string]any{"filter": filter}}
	} else {
		query = map[string]any{
			"bool": map[string]any{
				"filter": filter,
				"must": []map[string]any{
					{
						"multi_match": map[string]any{
							"query":  keyword,
							"fields": []string{"title"},
						},
					},
				},
			},
		}
	}
	return map[string]any{
		"from":  from,
		"size":  pageSize,
		"query": query,
		"sort": []map[string]any{
			{"createdAt": map[string]any{"order": "desc"}},
		},
	}
}
