// Package workflow HTTP 契约层：媒资审核工作流 API。
package workflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 工作流 HTTP 处理器。
type Handler struct {
	svc Service
}

// NewHandler 构造工作流 Handler。
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// pageRequest 工作流分页请求体。
type pageRequest struct {
	Page     int `json:"page" example:"1"`
	PageSize int `json:"pageSize" example:"20"`
}

// submitRequest 提交审核请求体。
type submitRequest struct {
	AssetID       string `json:"assetId" binding:"required"`
	WorkflowDefID string `json:"workflowDefId" binding:"required"`
}

// auditRequest 单条审核请求体。
type auditRequest struct {
	InstanceID string `json:"instanceId" binding:"required"`
	Pass       bool   `json:"pass"`
	Remark     string `json:"remark"`
}

// auditItem 批量审核单项。
type auditItem struct {
	InstanceID string `json:"instanceId" binding:"required"`
	Pass       bool   `json:"pass"`
	Remark     string `json:"remark"`
}

// multiAuditRequest 批量审核请求体。
type multiAuditRequest struct {
	Items []auditItem `json:"items" binding:"required"`
}

// revokeRequest 撤回审核请求体。
type revokeRequest struct {
	InstanceIDs []string `json:"instanceIds" binding:"required"`
}

// levelUserItem 流程定义层级审核人。
type levelUserItem struct {
	Level  int    `json:"level" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}

// createDefRequest 创建流程定义请求体。
type createDefRequest struct {
	Name       string          `json:"name" binding:"required"`
	AuditLevel int             `json:"auditLevel" binding:"required"`
	LevelUsers []levelUserItem `json:"levelUsers"`
}

// Submit 提交审核
//
//	@Summary		提交审核
//	@Description	将媒资提交至指定流程定义，创建审核实例
//	@Tags			工作流
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		submitRequest	true	"提交参数"
//	@Success		201		{object}	InstanceVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/submit [post]
func (h *Handler) Submit(c *gin.Context) {
	var req submitRequest
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

// Audit 审核实例
//
//	@Summary		审核实例
//	@Description	对单个流程实例执行通过或打回
//	@Tags			工作流
//	@Accept			json
//	@Security		BearerAuth
//	@Param			body	body	auditRequest	true	"审核参数"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorVo
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		403	{object}	httpx.ErrorVo
//	@Failure		409	{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/audit [post]
func (h *Handler) Audit(c *gin.Context) {
	var req auditRequest
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

// MultiAudit 批量审核
//
//	@Summary		批量审核
//	@Tags			工作流
//	@Accept			json
//	@Security		BearerAuth
//	@Param			body	body	multiAuditRequest	true	"批量审核参数"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorVo
//	@Failure		401	{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/multi-audit [post]
func (h *Handler) MultiAudit(c *gin.Context) {
	var req multiAuditRequest
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

// Revoke 撤回审核
//
//	@Summary		撤回审核
//	@Description	撤回本人发起的未完成审核实例
//	@Tags			工作流
//	@Accept			json
//	@Security		BearerAuth
//	@Param			body	body	revokeRequest	true	"实例 ID 列表"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorVo
//	@Failure		401	{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/revoke [post]
func (h *Handler) Revoke(c *gin.Context) {
	var req revokeRequest
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

// AssignedToMe 待我审核列表
//
//	@Summary		待我审核列表
//	@Tags			工作流
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	PageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/assigned-to-me [post]
func (h *Handler) AssignedToMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAssignedToMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

// CreatedByMe 我发起的审核
//
//	@Summary		我发起的审核
//	@Tags			工作流
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	PageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/created-by-me [post]
func (h *Handler) CreatedByMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListCreatedByMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

// AuditedByMe 我审过的
//
//	@Summary		我审过的
//	@Tags			工作流
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	PageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/audited-by-me [post]
func (h *Handler) AuditedByMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAuditedByMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

// All 全部审核实例
//
//	@Summary		全部审核实例
//	@Tags			工作流
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	PageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset-workflow/all [post]
func (h *Handler) All(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAll(c.Request.Context(), page, pageSize)
	})
}

// CreateDef 创建流程定义
//
//	@Summary		创建流程定义
//	@Tags			工作流定义
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createDefRequest	true	"流程定义"
//	@Success		201		{object}	DefVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/workflow/def [post]
func (h *Handler) CreateDef(c *gin.Context) {
	var req createDefRequest
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

// DefPage 流程定义分页
//
//	@Summary		流程定义分页
//	@Tags			工作流定义
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	DefPageResult
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/workflow/def/page [post]
func (h *Handler) DefPage(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.DefPage(c.Request.Context(), page, pageSize)
	})
}
