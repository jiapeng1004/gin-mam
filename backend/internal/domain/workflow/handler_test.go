package workflow_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/workflow"
	"github.com/gin-mam/backend/internal/domain/workflow/mock"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_AssignedToMe_Returns200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockSvc.EXPECT().
		ListAssignedToMe(gomock.Any(), "user-1", 1, 20).
		Return(&workflow.PageResult{List: []workflow.InstanceVO{}, Total: 0, Page: 1, PageSize: 20}, nil)

	h := workflow.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/asset-workflow/assigned-to-me", func(c *gin.Context) {
		c.Set("userId", "user-1")
		h.AssignedToMe(c)
	})

	body, _ := json.Marshal(map[string]any{"page": 1, "pageSize": 20})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/asset-workflow/assigned-to-me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, float64(0), resp["total"])
}
