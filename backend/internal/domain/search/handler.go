// Package search HTTP 契约层：媒资 Elasticsearch 检索 API。
package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 检索 HTTP 处理器。
type Handler struct {
	svc Service
}

// NewHandler 构造检索 Handler。
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// searchAssetRequest POST /api/v1/search/asset 请求体。
type searchAssetRequest struct {
	Keyword   string `json:"keyword"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	Type      string `json:"type"`
	CatalogID string `json:"catalogId"`
}

// SearchAsset 全文检索媒资。
//
// 路由：POST /api/v1/search/asset
// 鉴权：JWT Bearer
// 请求体：{ keyword?, page, pageSize, type?, catalogId? }
// 成功：200 { list: AssetVO[], total, page, pageSize }（ES 命中后回源拼装完整媒资）
func (h *Handler) SearchAsset(c *gin.Context) {
	var req searchAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.SearchAsset(c.Request.Context(), SearchInput{
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Type:      req.Type,
		CatalogID: req.CatalogID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}
