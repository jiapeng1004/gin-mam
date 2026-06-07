package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
)

func TestFail_ReturnsSnakeCaseErrorVo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	httpx.Fail(c, http.StatusNotFound, 40401, "媒资不存在")
	require.Equal(t, 404, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, float64(40401), body["err_code"])
	require.Equal(t, "媒资不存在", body["err_msg"])
	require.NotContains(t, body, "errCode")
}

func TestOK_ReturnsDataDirectly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	httpx.OK(c, map[string]string{"previewUrl": "http://x"})
	require.Equal(t, 200, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, "http://x", body["previewUrl"])
	require.NotContains(t, body, "code")
	require.NotContains(t, body, "data")
}