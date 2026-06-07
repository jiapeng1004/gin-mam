package sys

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			_ = c.Error(biz)
			return
		}
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

func (h *Handler) UserPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.svc.UserPage(c.Request.Context(), page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}
