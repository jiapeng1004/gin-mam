// Package catalog HTTP 契约层：编目树与节点配置 API。
package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 编目 HTTP 处理器。
type Handler struct {
	svc Service
}

// NewHandler 构造编目 Handler。
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Tree 获取完整编目树。
//
// 路由：GET /api/v1/catalog/tree
// 鉴权：JWT Bearer
// 成功：200 TreeNode[]（空树返回 []）
func (h *Handler) Tree(c *gin.Context) {
	result, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	if result == nil {
		result = []TreeNode{}
	}
	httpx.OK(c, result)
}

// createRequest POST /api/v1/catalog 请求体。
type createRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parentId"`
	SortCode int     `json:"sortCode"`
}

// Create 创建编目节点。
//
// 路由：POST /api/v1/catalog
// 鉴权：JWT Bearer
// 请求体：{ name, parentId?, sortCode? }
// 成功：201 CatalogVO
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.Create(c.Request.Context(), CreateInput{
		Name:     req.Name,
		ParentID: req.ParentID,
		SortCode: req.SortCode,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.Created(c, result)
}

// updateRequest PUT /api/v1/catalog/:id 请求体。
type updateRequest struct {
	Name     *string `json:"name"`
	ParentID *string `json:"parentId"`
	SortCode *int    `json:"sortCode"`
}

// Update 更新编目节点。
//
// 路由：PUT /api/v1/catalog/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 请求体：{ name?, parentId?, sortCode? }（未传字段不修改）
// 成功：200 CatalogVO
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.Update(c.Request.Context(), id, UpdateInput{
		Name:     req.Name,
		ParentID: req.ParentID,
		SortCode: req.SortCode,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

// Delete 删除编目节点。
//
// 路由：DELETE /api/v1/catalog/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 成功：204 无响应体
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// ListConfig 列出编目节点配置项。
//
// 路由：GET /api/v1/catalog/:id/config
// 鉴权：JWT Bearer
// 路径参数：id（编目节点 ID）
// 成功：200 ConfigVO[]
func (h *Handler) ListConfig(c *gin.Context) {
	catalogID := c.Param("id")
	result, err := h.svc.ListConfig(c.Request.Context(), catalogID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if result == nil {
		result = []ConfigVO{}
	}
	httpx.OK(c, result)
}

// batchConfigRequest PUT /api/v1/catalog/:id/config 请求体。
type batchConfigRequest struct {
	Items []struct {
		ConfigKey   string `json:"configKey" binding:"required"`
		ConfigValue string `json:"configValue"`
	} `json:"items" binding:"required"`
}

// BatchUpdateConfig 批量更新编目节点配置。
//
// 路由：PUT /api/v1/catalog/:id/config
// 鉴权：JWT Bearer
// 路径参数：id（编目节点 ID）
// 请求体：{ items: [{ configKey, configValue }] }
// 成功：200 ConfigVO[]
func (h *Handler) BatchUpdateConfig(c *gin.Context) {
	catalogID := c.Param("id")
	var req batchConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	items := make([]ConfigItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, ConfigItemInput{
			ConfigKey:   item.ConfigKey,
			ConfigValue: item.ConfigValue,
		})
	}
	result, err := h.svc.BatchUpdateConfig(c.Request.Context(), catalogID, items)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}
