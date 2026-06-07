package sys_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/domain/sys/mock"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_LoginSuccess_ReturnsTokenField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockService(ctrl)
	expires := time.Date(2026, 6, 8, 12, 0, 0, 0, time.UTC)
	mockSvc.EXPECT().
		Login(gomock.Any(), "admin", "admin123").
		Return(&sys.LoginResult{Token: "jwt-token", ExpiresAt: expires}, nil)

	h := sys.NewHandler(mockSvc)
	r := gin.New()
	r.Use(httpx.ErrorMiddleware())
	r.POST("/api/v1/auth/login", h.Login)

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "admin123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "jwt-token", resp["token"])
	require.NotEmpty(t, resp["expiresAt"])
}
