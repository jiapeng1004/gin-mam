// Package transcode HTTP 契约层：转码组/任务管理与外部回调 API。
package transcode

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 转码 HTTP 处理器。
type Handler struct {
	svc Service
}

// NewHandler 构造转码 Handler。
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// pageRequest 转码分页 POST 接口通用请求体。
type pageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// GroupPage 分页查询转码组列表。
//
// 路由：POST /api/v1/transcode/group/page
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) GroupPage(c *gin.Context) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := h.svc.GroupPage(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

// CreateGroup 创建转码组。
//
// 路由：POST /api/v1/transcode/group
// 鉴权：JWT Bearer
// 请求体：{ name, message?, groupType?, param?, strategyType?, defaultFlag?, availableFlag? }
// 成功：201 GroupVO
func (h *Handler) CreateGroup(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Message       string `json:"message"`
		GroupType     int    `json:"groupType"`
		Param         string `json:"param"`
		StrategyType  int    `json:"strategyType"`
		DefaultFlag   int8   `json:"defaultFlag"`
		AvailableFlag int8   `json:"availableFlag"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := h.svc.CreateGroup(c.Request.Context(), httpx.GetUserID(c), CreateGroupInput{
		Name:          req.Name,
		Message:       req.Message,
		GroupType:     req.GroupType,
		Param:         req.Param,
		StrategyType:  req.StrategyType,
		DefaultFlag:   req.DefaultFlag,
		AvailableFlag: req.AvailableFlag,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.Created(c, out)
}

// UpdateGroup 更新转码组。
//
// 路由：PUT /api/v1/transcode/group/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 请求体：各字段可选，未传不修改
// 成功：200 GroupVO
// 失败：404 { err_code: 40403 }
func (h *Handler) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name          *string `json:"name"`
		Message       *string `json:"message"`
		GroupType     *int    `json:"groupType"`
		Param         *string `json:"param"`
		StrategyType  *int    `json:"strategyType"`
		DefaultFlag   *int8   `json:"defaultFlag"`
		AvailableFlag *int8   `json:"availableFlag"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := h.svc.UpdateGroup(c.Request.Context(), id, UpdateGroupInput{
		Name:          req.Name,
		Message:       req.Message,
		GroupType:     req.GroupType,
		Param:         req.Param,
		StrategyType:  req.StrategyType,
		DefaultFlag:   req.DefaultFlag,
		AvailableFlag: req.AvailableFlag,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

// DeleteGroup 删除转码组。
//
// 路由：DELETE /api/v1/transcode/group/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 成功：204 无响应体
func (h *Handler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteGroup(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// CreateTask 创建转码任务并提交外部转码服务。
//
// 路由：POST /api/v1/transcode/task
// 鉴权：JWT Bearer
// 请求体：{ assetId, assetFileId?, transcodeGroupId }
// 成功：201 TaskVO
func (h *Handler) CreateTask(c *gin.Context) {
	var req struct {
		AssetID          string  `json:"assetId" binding:"required"`
		AssetFileID      *string `json:"assetFileId"`
		TranscodeGroupID string  `json:"transcodeGroupId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := h.svc.CreateTask(c.Request.Context(), httpx.GetUserID(c), CreateTaskInput{
		AssetID:          req.AssetID,
		AssetFileID:      req.AssetFileID,
		TranscodeGroupID: req.TranscodeGroupID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.Created(c, out)
}

// GetTask 按 ID 获取转码任务详情。
//
// 路由：GET /api/v1/transcode/task/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 成功：200 TaskVO
// 失败：404 { err_code: 40403 }
func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetTask(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

// TaskPage 分页查询转码任务列表。
//
// 路由：POST /api/v1/transcode/task/page
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) TaskPage(c *gin.Context) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := h.svc.TaskPage(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

// Callback 接收外部转码服务 HMAC 签名回调。
//
// 路由：POST /api/v1/transcode/callback
// 鉴权：无 JWT；校验请求头 X-Transcode-Sign
// 请求体：原始 JSON（taskId、status、outputPath、errorMsg）
// 成功：204 无响应体
// 失败：403 { err_code: 40300 } 签名校验失败
func (h *Handler) Callback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	if err := h.svc.HandleCallback(c.Request.Context(), body, c.GetHeader("X-Transcode-Sign")); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}
