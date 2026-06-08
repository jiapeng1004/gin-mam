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

// pageRequest 转码分页请求体。
type pageRequest struct {
	Page     int `json:"page" example:"1"`
	PageSize int `json:"pageSize" example:"20"`
}

// createGroupRequest 创建转码组请求体。
type createGroupRequest struct {
	Name          string `json:"name" binding:"required"`
	Message       string `json:"message"`
	GroupType     int    `json:"groupType"`
	Param         string `json:"param"`
	StrategyType  int    `json:"strategyType"`
	DefaultFlag   int8   `json:"defaultFlag"`
	AvailableFlag int8   `json:"availableFlag"`
}

// updateGroupRequest 更新转码组请求体。
type updateGroupRequest struct {
	Name          *string `json:"name"`
	Message       *string `json:"message"`
	GroupType     *int    `json:"groupType"`
	Param         *string `json:"param"`
	StrategyType  *int    `json:"strategyType"`
	DefaultFlag   *int8   `json:"defaultFlag"`
	AvailableFlag *int8   `json:"availableFlag"`
}

// createTaskRequest 创建转码任务请求体。
type createTaskRequest struct {
	AssetID          string  `json:"assetId" binding:"required"`
	AssetFileID      *string `json:"assetFileId"`
	TranscodeGroupID string  `json:"transcodeGroupId" binding:"required"`
}

// GroupPage 转码组分页
//
//	@Summary		转码组分页
//	@Tags			转码
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	GroupPageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/group/page [post]
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

// CreateGroup 创建转码组
//
//	@Summary		创建转码组
//	@Tags			转码
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createGroupRequest	true	"转码组参数"
//	@Success		201		{object}	GroupVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/group [post]
func (h *Handler) CreateGroup(c *gin.Context) {
	var req createGroupRequest
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

// UpdateGroup 更新转码组
//
//	@Summary		更新转码组
//	@Tags			转码
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"转码组 ID"
//	@Param			body	body		updateGroupRequest	true	"更新参数"
//	@Success		200		{object}	GroupVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Failure		404		{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/group/{id} [put]
func (h *Handler) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req updateGroupRequest
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

// DeleteGroup 删除转码组
//
//	@Summary		删除转码组
//	@Tags			转码
//	@Security		BearerAuth
//	@Param			id	path	string	true	"转码组 ID"
//	@Success		204
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/group/{id} [delete]
func (h *Handler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteGroup(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// CreateTask 创建转码任务
//
//	@Summary		创建转码任务
//	@Description	创建任务并提交外部转码服务
//	@Tags			转码
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createTaskRequest	true	"任务参数"
//	@Success		201		{object}	TaskVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/task [post]
func (h *Handler) CreateTask(c *gin.Context) {
	var req createTaskRequest
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

// GetTask 获取转码任务
//
//	@Summary		获取转码任务
//	@Tags			转码
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"任务 ID"
//	@Success		200	{object}	TaskVO
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/task/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetTask(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

// TaskPage 转码任务分页
//
//	@Summary		转码任务分页
//	@Tags			转码
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	TaskPageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/task/page [post]
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

// Callback 转码回调
//
//	@Summary		转码回调
//	@Description	外部转码服务 HMAC 签名回调，无需 JWT
//	@Tags			转码
//	@Accept			json
//	@Param			X-Transcode-Sign	header		string			true	"HMAC 签名"
//	@Param			body				body		CallbackInput	true	"回调体"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorVo
//	@Failure		403	{object}	httpx.ErrorVo
//	@Router			/api/v1/transcode/callback [post]
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
