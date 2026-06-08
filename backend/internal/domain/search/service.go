// Package search 检索领域：基于 Elasticsearch 的媒资全文检索。
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

// SearchInput 媒资检索入参。
type SearchInput struct {
	// Keyword 标题关键字，空则仅按 filter 条件查询。
	Keyword string
	// Page 页码，从 1 开始。
	Page int
	// PageSize 每页条数。
	PageSize int
	// Type 媒资类型过滤，空则不过滤。
	Type string
	// CatalogID 编目 ID 过滤，空则不过滤。
	CatalogID string
}

// AssetDocument Elasticsearch 索引文档结构。
type AssetDocument struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	CatalogID *string   `json:"catalogId,omitempty"`
	TenantID  string    `json:"tenantId"`
	CreatedAt time.Time `json:"createdAt"`
}

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=service.go -destination=mock/service_mock.go -package=mock

// Service 检索应用服务接口。
type Service interface {
	// SearchAsset 在 ES 中检索媒资 ID，再回源 asset.Service 拼装完整 AssetVO 列表。
	SearchAsset(ctx context.Context, in SearchInput) (*asset.PageResult, error)
	// ScheduleIndex 异步写入/更新 ES 索引文档，媒资变更后由 asset 层调用。
	ScheduleIndex(doc asset.IndexDoc)
}

type service struct {
	esClient elasticsearch.Client
	assets   asset.Service
}

// NewService 构造检索服务，依赖 ES 客户端与媒资服务。
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

// MarshalAssetDocument 将 IndexDoc 序列化为 ES 索引 JSON。
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

// BuildAssetSearchQuery 构造 ES bool 查询 DSL。
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
