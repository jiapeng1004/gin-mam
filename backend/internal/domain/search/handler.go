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

// searchAssetRequest 媒资检索请求体。
type searchAssetRequest struct {
	Keyword   string `json:"keyword" example:"宣传片"`
	Page      int    `json:"page" example:"1"`
	PageSize  int    `json:"pageSize" example:"20"`
	Type      string `json:"type" example:"video"`
	CatalogID string `json:"catalogId" example:""`
}

// SearchAsset 检索媒资
//
//	@Summary		检索媒资
//	@Description	基于 Elasticsearch 全文检索媒资，命中后回源拼装 AssetVO
//	@Tags			检索
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		searchAssetRequest	true	"检索参数"
//	@Success		200		{object}	asset.PageResult
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/search/asset [post]
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
