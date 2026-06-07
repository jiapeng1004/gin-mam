package search_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/search"
	"github.com/stretchr/testify/require"
)

func TestMarshalAssetDocument(t *testing.T) {
	catalogID := "cat1"
	doc := asset.IndexDoc{
		ID:        "a1",
		Title:     "hello",
		Type:      "video",
		CatalogID: &catalogID,
		TenantID:  "default",
		CreatedAt: time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC),
	}
	raw, err := search.MarshalAssetDocument(doc)
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	require.Equal(t, "a1", m["id"])
	require.Equal(t, "hello", m["title"])
	require.Equal(t, "video", m["type"])
	require.Equal(t, "cat1", m["catalogId"])
	require.Equal(t, "default", m["tenantId"])
	require.NotEmpty(t, m["createdAt"])
}

func TestBuildAssetSearchQuery_KeywordAndFilters(t *testing.T) {
	q := search.BuildAssetSearchQuery("default", search.SearchInput{
		Keyword:   "news",
		Page:      2,
		PageSize:  10,
		Type:      "image",
		CatalogID: "c1",
	}, 2, 10)
	require.Equal(t, 10, q["from"])
	require.Equal(t, 10, q["size"])
	raw, err := json.Marshal(q)
	require.NoError(t, err)
	s := string(raw)
	require.Contains(t, s, `"tenantId"`)
	require.Contains(t, s, `"multi_match"`)
	require.Contains(t, s, "news")
	require.Contains(t, s, `"type"`)
	require.Contains(t, s, `"catalogId"`)
}

func TestBuildAssetSearchQuery_NoKeyword(t *testing.T) {
	q := search.BuildAssetSearchQuery("default", search.SearchInput{Page: 1, PageSize: 20}, 1, 20)
	boolQ := q["query"].(map[string]any)["bool"].(map[string]any)
	_, hasMust := boolQ["must"]
	require.False(t, hasMust)
	require.NotNil(t, boolQ["filter"])
}
