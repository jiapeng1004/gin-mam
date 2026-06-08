package transcode_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/transcode"
	"github.com/gin-mam/backend/internal/domain/transcode/mock"
	transcodemodel "github.com/gin-mam/backend/internal/domain/transcode/model"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_Callback_InvalidSign403(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockSvc.EXPECT().HandleCallback(gomock.Any(), gomock.Any(), "bad").Return(httpx.NewBizError(403, transcode.ErrCodeForbidden, "签名校验失败"))

	h := transcode.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/transcode/callback", h.Callback)

	body := []byte(`{"taskId":"t1","status":2}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transcode/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transcode-Sign", "bad")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_GroupPage_Returns200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockSvc.EXPECT().GroupPage(gomock.Any(), 1, 20).Return(&transcode.GroupPageResult{
		List: []transcode.GroupVO{}, Total: 0, Page: 1, PageSize: 20,
	}, nil)

	h := transcode.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/transcode/group/page", h.GroupPage)

	payload, _ := json.Marshal(map[string]any{"page": 1, "pageSize": 20})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transcode/group/page", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_CreateTask_Returns201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	mockSvc.EXPECT().CreateTask(gomock.Any(), "user-1", transcode.CreateTaskInput{
		AssetID: "a1", TranscodeGroupID: "g1",
	}).Return(&transcode.TaskVO{ID: "t1", AssetID: "a1", TranscodeGroupID: "g1", Status: transcodemodel.TaskStatusRunning}, nil)

	h := transcode.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/transcode/task", func(c *gin.Context) {
		c.Set("userId", "user-1")
		h.CreateTask(c)
	})

	payload, _ := json.Marshal(map[string]any{"assetId": "a1", "transcodeGroupId": "g1"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transcode/task", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}