package transcode

import (
	"io"
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

type pageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

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

func (h *Handler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteGroup(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

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

func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetTask(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

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