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

// Tree 获取编目树
//
//	@Summary		获取编目树
//	@Description	返回完整多级编目树结构
//	@Tags			编目
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		TreeNode
//	@Failure		401	{object}	httpx.ErrorVo
//	@Router			/api/v1/catalog/tree [get]
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

// createRequest 创建编目请求体。
type createRequest struct {
	Name     string  `json:"name" binding:"required" example:"新闻栏目"`
	ParentID *string `json:"parentId" example:""`
	SortCode int     `json:"sortCode" example:"0"`
}

// Create 创建编目节点
//
//	@Summary		创建编目节点
//	@Tags			编目
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createRequest	true	"编目参数"
//	@Success		201		{object}	CatalogVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/catalog [post]
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

// updateRequest 更新编目请求体。
type updateRequest struct {
	Name     *string `json:"name"`
	ParentID *string `json:"parentId"`
	SortCode *int    `json:"sortCode"`
}

// Update 更新编目节点
//
//	@Summary		更新编目节点
//	@Tags			编目
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string			true	"编目 ID"
//	@Param			body	body		updateRequest	true	"更新参数"
//	@Success		200		{object}	CatalogVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Failure		404		{object}	httpx.ErrorVo
//	@Router			/api/v1/catalog/{id} [put]
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

// Delete 删除编目节点
//
//	@Summary		删除编目节点
//	@Tags			编目
//	@Security		BearerAuth
//	@Param			id	path	string	true	"编目 ID"
//	@Success		204
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/catalog/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// ListConfig 列出编目配置
//
//	@Summary		列出编目配置
//	@Tags			编目
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"编目 ID"
//	@Success		200	{array}		ConfigVO
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/catalog/{id}/config [get]
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

// batchConfigRequest 批量更新编目配置请求体。
type batchConfigRequest struct {
	Items []struct {
		ConfigKey   string `json:"configKey" binding:"required"`
		ConfigValue string `json:"configValue"`
	} `json:"items" binding:"required"`
}

// BatchUpdateConfig 批量更新编目配置
//
//	@Summary		批量更新编目配置
//	@Tags			编目
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"编目 ID"
//	@Param			body	body		batchConfigRequest	true	"配置项列表"
//	@Success		200		{array}		ConfigVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/catalog/{id}/config [put]
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
