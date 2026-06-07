package catalog

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

type createRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parentId"`
	SortCode int     `json:"sortCode"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

type updateRequest struct {
	Name     *string `json:"name"`
	ParentID *string `json:"parentId"`
	SortCode *int    `json:"sortCode"`
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

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

type batchConfigRequest struct {
	Items []struct {
		ConfigKey   string `json:"configKey" binding:"required"`
		ConfigValue string `json:"configValue"`
	} `json:"items" binding:"required"`
}

func (h *Handler) BatchUpdateConfig(c *gin.Context) {
	catalogID := c.Param("id")
	var req batchConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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