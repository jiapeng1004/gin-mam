// Package sys HTTP 契约层：系统认证与用户管理 API。
package sys

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 系统管理 HTTP 处理器，将 Gin 请求委托给 Service。
type Handler struct {
	svc Service
}

// NewHandler 构造系统管理 Handler。
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// loginRequest POST /api/v1/auth/login 请求体。
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录。
//
// 路由：POST /api/v1/auth/login
// 鉴权：无
// 请求体：{ username, password }
// 成功：200 { token, expiresAt }
// 失败：401 { err_code: 40100, err_msg }
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

// UserPage 分页查询用户列表。
//
// 路由：GET /api/v1/sys/user/page
// 鉴权：JWT Bearer
// 查询参数：page（默认 1）、pageSize（默认 20）
// 成功：200 { list, total, page, pageSize }
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
