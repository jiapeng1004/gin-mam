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

// pageRequest 工作流分页类 POST 接口通用请求体。
type pageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// Submit 提交媒资进入审核流程。
//
// 路由：POST /api/v1/asset-workflow/submit
// 鉴权：JWT Bearer
// 请求体：{ assetId, workflowDefId }
// 成功：201 InstanceVO
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

// Audit 对单个流程实例执行通过或打回。
//
// 路由：POST /api/v1/asset-workflow/audit
// 鉴权：JWT Bearer（须为当前层级审核人）
// 请求体：{ instanceId, pass, remark? }
// 成功：204 无响应体
// 失败：403 无权审核、409 工作流锁定
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

// MultiAudit 批量审核多个流程实例。
//
// 路由：POST /api/v1/asset-workflow/multi-audit
// 鉴权：JWT Bearer
// 请求体：{ items: [{ instanceId, pass, remark? }] }
// 成功：204 无响应体
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

// Revoke 撤回本人发起的未完成审核实例。
//
// 路由：POST /api/v1/asset-workflow/revoke
// 鉴权：JWT Bearer
// 请求体：{ instanceIds: string[] }
// 成功：204 无响应体
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

// postPage 工作流分页 POST 接口通用处理：解析 page/pageSize 并委托查询函数。
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

// AssignedToMe 分页查询待当前用户审核的实例。
//
// 路由：POST /api/v1/asset-workflow/assigned-to-me
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) AssignedToMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAssignedToMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

// CreatedByMe 分页查询当前用户发起的审核实例。
//
// 路由：POST /api/v1/asset-workflow/created-by-me
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) CreatedByMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListCreatedByMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

// AuditedByMe 分页查询当前用户已审核过的实例。
//
// 路由：POST /api/v1/asset-workflow/audited-by-me
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) AuditedByMe(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAuditedByMe(c.Request.Context(), httpx.GetUserID(c), page, pageSize)
	})
}

// All 分页查询全部审核实例（管理视角）。
//
// 路由：POST /api/v1/asset-workflow/all
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) All(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.ListAll(c.Request.Context(), page, pageSize)
	})
}

// CreateDef 创建审核流程定义。
//
// 路由：POST /api/v1/workflow/def
// 鉴权：JWT Bearer
// 请求体：{ name, auditLevel, levelUsers: [{ level, userId }] }
// 成功：201 DefVO
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

// DefPage 分页查询流程定义列表。
//
// 路由：POST /api/v1/workflow/def/page
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) DefPage(c *gin.Context) {
	h.postPage(c, func(c *gin.Context, page, pageSize int) (any, error) {
		return h.svc.DefPage(c.Request.Context(), page, pageSize)
	})
}
