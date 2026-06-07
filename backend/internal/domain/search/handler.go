package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type searchAssetRequest struct {
	Keyword   string `json:"keyword"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	Type      string `json:"type"`
	CatalogID string `json:"catalogId"`
}

func (h *Handler) SearchAsset(c *gin.Context) {
	var req searchAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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
