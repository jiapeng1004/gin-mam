// Package sys HTTP 契约层：系统认证与用户管理 API。
package sys

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 系统管理 HTTP 处理器。
type Handler struct {
	svc Service
}

// NewHandler 构造系统管理 Handler。
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// loginRequest 登录请求体。
type loginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	用户名密码登录，返回 JWT Token
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginRequest	true	"登录参数"
//	@Success		200		{object}	LoginResult
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/auth/login [post]
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

// UserPage 分页查询用户
//
//	@Summary		分页查询用户
//	@Description	获取系统用户分页列表
//	@Tags			系统用户
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"页码"	default(1)
//	@Param			pageSize	query		int	false	"每页条数"	default(20)
//	@Success		200			{object}	UserPageResult
//	@Failure		401			{object}	httpx.ErrorVo
//	@Router			/api/v1/sys/user/page [get]
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
