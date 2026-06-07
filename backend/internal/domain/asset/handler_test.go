package asset_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/asset/mock"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_Get_NotFound_Returns40401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockSvc.EXPECT().
		GetByID(gomock.Any(), "missing").
		Return(nil, httpx.NewBizError(404, asset.ErrCodeNotFound, "\u5a92\u8d44\u4e0d\u5b58\u5728"))

	h := asset.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.GET("/api/v1/asset/:id", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/asset/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	var resp httpx.ErrorVo
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, asset.ErrCodeNotFound, resp.ErrCode)
}