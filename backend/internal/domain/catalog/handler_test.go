package catalog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/catalog"
	"github.com/gin-mam/backend/internal/domain/catalog/mock"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_Tree_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	parent := "p1"
	mockSvc.EXPECT().Tree(gomock.Any()).Return([]catalog.TreeNode{
		{ID: "p1", Name: "parent", ParentID: nil, Children: []catalog.TreeNode{
			{ID: "c1", Name: "child", ParentID: &parent, Children: []catalog.TreeNode{}},
		}},
	}, nil)

	h := catalog.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.GET("/api/v1/catalog/tree", h.Tree)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/tree", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []catalog.TreeNode
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp, 1)
	require.Len(t, resp[0].Children, 1)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockSvc.EXPECT().Delete(gomock.Any(), "missing").Return(
		httpx.NewBizError(404, catalog.ErrCodeNotFound, "\u680f\u76ee\u4e0d\u5b58\u5728"),
	)

	h := catalog.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.DELETE("/api/v1/catalog/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}