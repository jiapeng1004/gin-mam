package workflow

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

type pageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (h *Handler) Submit(c *gin.Context) {
	var req struct {
		AssetID       string `json:"assetId" binding:"required"`
		WorkflowDefID string `json:"workflowDefId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := h.svc.Submit(c.Request.Context(), httpx.GetUserID(c), SubmitRequest{
		AssetID:       req.AssetID,
		WorkflowDefID: req.WorkflowDefID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.Created(c, out)
}

func (h *Handler) Audit(c *gin.Context) {
	var req struct {
		InstanceID string `json:"instanceId" binding:"required"`
		Pass       bool   `json:"pass"`
		Remark     string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	if err := h.svc.Audit(c.Request.Context(), httpx.GetUserID(c), AuditRequest{
		InstanceID: req.InstanceID,
		Pass:       req.Pass,
		Remark:     req.Remark,
	}); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

func (h *Handler) MultiAudit(c *gin.Context) {
	var req struct {
		Items []struct {
			InstanceID string `json:"instanceId" binding:"required"`
			Pass       bool   `json:"pass"`
			Remark     string `json:"remark"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	items := make([]AuditRequest, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, AuditRequest{InstanceID: it.InstanceID, Pass: it.Pass, Remark: it.Remark})
	}
	if err := h.svc.MultiAudit(c.Request.Context(), httpx.GetUserID(c), items); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

func (h *Handler) Revoke(c *gin.Context) {
	var req struct {
		InstanceIDs []string `json:"instanceIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	if err := h.svc.Revoke(c.Request.Context(), httpx.GetUserID(c), req.InstanceIDs); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

func (h *Handler) postPage(c *gin.Context, fn func(ctx *gin.Context, page, pageSize int) (any, error)) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	out, err := fn(c, req.Page, req.PageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, out)
}

func (h *Handler) AssignedToMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAssignedToMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

func (h *Handler) CreatedByMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListCreatedByMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

func (h *Handler) AuditedByMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAuditedByMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

func (h *Handler) All(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAll(c.Request.Context(), page, pageSize)
	})
}

func (h *Handler) CreateDef(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		AuditLevel int    `json:"auditLevel" binding:"required"`
		LevelUsers []struct {
			Level  int    `json:"level" binding:"required"`
			UserID string `json:"userId" binding:"required"`
		} `json:"levelUsers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	lu := make([]DefLevelUserVO, 0, len(req.LevelUsers))
	for _, row := range req.LevelUsers {
		lu = append(lu, DefLevelUserVO{Level: row.Level, UserID: row.UserID})
	}
	out, err := h.svc.CreateDef(c.Request.Context(), CreateDefInput{
		Name:       req.Name,
		AuditLevel: req.AuditLevel,
		LevelUsers: lu,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.Created(c, out)
}

func (h *Handler) DefPage(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.DefPage(c.Request.Context(), page, pageSize)
	})
}
