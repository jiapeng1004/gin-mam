package search_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	assetmock "github.com/gin-mam/backend/internal/domain/asset/mock"
	"github.com/gin-mam/backend/internal/domain/search"
	searchmock "github.com/gin-mam/backend/internal/domain/search/mock"
	"github.com/gin-mam/backend/internal/infra/elasticsearch"
	esmock "github.com/gin-mam/backend/internal/infra/elasticsearch/mock"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_SearchAsset_BuildsESRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	es := esmock.NewMockClient(ctrl)
	assets := assetmock.NewMockService(ctrl)

	es.EXPECT().Search(gomock.Any(), elasticsearch.AssetIndexName, gomock.Any()).DoAndReturn(
		func(_ context.Context, _ string, body []byte) (*elasticsearch.SearchResponse, error) {
			var q map[string]any
			require.NoError(t, json.Unmarshal(body, &q))
			require.Equal(t, float64(0), q["from"])
			require.Equal(t, float64(20), q["size"])
			return &elasticsearch.SearchResponse{Total: 1, Hits: []elasticsearch.SearchHit{{ID: "a1"}}}, nil
		},
	)
	assets.EXPECT().GetByID(gomock.Any(), "a1").Return(&asset.AssetVO{ID: "a1", Title: "t"}, nil)

	svc := search.NewService(es, assets)
	result, err := svc.SearchAsset(context.Background(), search.SearchInput{Keyword: "kw", Page: 1, PageSize: 20})
	require.NoError(t, err)
	require.Equal(t, int64(1), result.Total)
	require.Len(t, result.List, 1)
	require.Equal(t, "a1", result.List[0].ID)
}

func TestHandler_SearchAsset_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := searchmock.NewMockService(ctrl)
	mockSvc.EXPECT().SearchAsset(gomock.Any(), search.SearchInput{
		Keyword: "x", Page: 1, PageSize: 10, Type: "video", CatalogID: "c1",
	}).Return(&asset.PageResult{List: []asset.AssetVO{}, Total: 0, Page: 1, PageSize: 10}, nil)

	h := search.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/search/asset", h.SearchAsset)

	body := []byte(`{"keyword":"x","page":1,"pageSize":10,"type":"video","catalogId":"c1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/search/asset", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
