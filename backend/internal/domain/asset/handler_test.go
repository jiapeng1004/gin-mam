package asset_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	mockUpload := mock.NewMockUploadService(ctrl)
	h := asset.NewHandler(mockSvc, mockUpload)
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
func TestHandler_UploadInit_ReturnsUploadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockUpload := mock.NewMockUploadService(ctrl)
	mockUpload.EXPECT().Init(gomock.Any(), gomock.Any()).Return(&asset.UploadInitResult{
		UploadID:  "upload123",
		ChunkSize: asset.DefaultUploadChunkSize,
	}, nil)

	h := asset.NewHandler(mockSvc, mockUpload)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/asset/upload/init", h.UploadInit)

	body := `{"fileName":"a.mp4","fileSize":100,"mimeType":"video/mp4"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/asset/upload/init", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp asset.UploadInitResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "upload123", resp.UploadID)
	require.Equal(t, int64(asset.DefaultUploadChunkSize), resp.ChunkSize)
}