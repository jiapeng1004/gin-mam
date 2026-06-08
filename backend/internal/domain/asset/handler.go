// Package asset HTTP 契约层：媒资 CRUD 与分片上传 API。
package asset

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

// Handler 媒资 HTTP 处理器。
type Handler struct {
	svc    Service
	upload UploadService
}

// NewHandler 构造媒资 Handler。
func NewHandler(svc Service, upload UploadService) *Handler {
	return &Handler{svc: svc, upload: upload}
}

// pageRequest POST /api/v1/asset/page 请求体。
type pageRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
}

// Page 分页查询媒资列表。
//
// 路由：POST /api/v1/asset/page
// 鉴权：JWT Bearer
// 请求体：{ page, pageSize, keyword? }
// 成功：200 { list, total, page, pageSize }
func (h *Handler) Page(c *gin.Context) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.Page(c.Request.Context(), PageInput{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

// Get 按 ID 获取媒资详情。
//
// 路由：GET /api/v1/asset/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 成功：200 AssetVO（含 previewUrl、metadata）
// 失败：404 { err_code: 40401 }
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

// createRequest POST /api/v1/asset 请求体。
type createRequest struct {
	Title       string            `json:"title" binding:"required"`
	Type        string            `json:"type" binding:"required"`
	CatalogID   *string           `json:"catalogId"`
	Description string            `json:"description"`
	StoragePath string            `json:"storagePath" binding:"required"`
	MimeType    string            `json:"mimeType"`
	FileSize    *int64            `json:"fileSize"`
	Metadata    map[string]string `json:"metadata"`
	CreatedBy   string            `json:"createdBy"`
}

// Create 直接创建媒资（已有 storagePath 时使用，非分片上传流程）。
//
// 路由：POST /api/v1/asset
// 鉴权：JWT Bearer
// 请求体：createRequest
// 成功：201 AssetVO
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.Create(c.Request.Context(), CreateInput{
		Title:       req.Title,
		Type:        req.Type,
		CatalogID:   req.CatalogID,
		Description: req.Description,
		StoragePath: req.StoragePath,
		MimeType:    req.MimeType,
		FileSize:    req.FileSize,
		Metadata:    req.Metadata,
		CreatedBy:   req.CreatedBy,
	})
	if err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			_ = c.Error(biz)
			return
		}
		_ = c.Error(err)
		return
	}
	httpx.Created(c, result)
}

// updateRequest PUT /api/v1/asset/:id 请求体。
type updateRequest struct {
	Title       *string           `json:"title"`
	CatalogID   *string           `json:"catalogId"`
	Description *string           `json:"description"`
	Status      *int8             `json:"status"`
	Metadata    map[string]string `json:"metadata"`
}

// Update 更新媒资。
//
// 路由：PUT /api/v1/asset/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 请求体：updateRequest（指针字段 nil 表示不修改）
// 成功：200 AssetVO
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.svc.Update(c.Request.Context(), id, UpdateInput{
		Title:       req.Title,
		CatalogID:   req.CatalogID,
		Description: req.Description,
		Status:      req.Status,
		Metadata:    req.Metadata,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

// Delete 软删除媒资。
//
// 路由：DELETE /api/v1/asset/:id
// 鉴权：JWT Bearer
// 路径参数：id
// 成功：204 无响应体
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// uploadInitRequest POST /api/v1/asset/upload/init 请求体。
type uploadInitRequest struct {
	FileName  string `json:"fileName" binding:"required"`
	FileSize  int64  `json:"fileSize" binding:"required"`
	MimeType  string `json:"mimeType"`
	ChunkSize *int64 `json:"chunkSize"`
}

// UploadInit 初始化分片上传会话。
//
// 路由：POST /api/v1/asset/upload/init
// 鉴权：JWT Bearer
// 请求体：{ fileName, fileSize, mimeType?, chunkSize? }
// 成功：200 { uploadId, chunkSize }
func (h *Handler) UploadInit(c *gin.Context) {
	var req uploadInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.upload.Init(c.Request.Context(), UploadInitInput{
		FileName:  req.FileName,
		FileSize:  req.FileSize,
		MimeType:  req.MimeType,
		ChunkSize: req.ChunkSize,
	})
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

// UploadChunk 上传单个分片。
//
// 路由：POST /api/v1/asset/upload/chunk
// 鉴权：JWT Bearer
// 表单字段：uploadId、chunkIndex、file（二进制分片）
// 成功：204 无响应体
// 失败：404 { err_code: 40402 } 上传会话不存在
func (h *Handler) UploadChunk(c *gin.Context) {
	uploadID := c.PostForm("uploadId")
	chunkIndexStr := c.PostForm("chunkIndex")
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || uploadID == "" {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	f, err := file.Open()
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer f.Close()
	if err := h.upload.SaveChunk(c.Request.Context(), UploadChunkInput{
		UploadID:   uploadID,
		ChunkIndex: chunkIndex,
		Body:       f,
	}); err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			_ = c.Error(biz)
			return
		}
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// uploadCompleteRequest POST /api/v1/asset/upload/complete 请求体。
type uploadCompleteRequest struct {
	UploadID   string  `json:"uploadId" binding:"required"`
	AssetTitle string  `json:"assetTitle" binding:"required"`
	CatalogID  *string `json:"catalogId"`
	Type       string  `json:"type" binding:"required"`
}

// UploadComplete 完成分片上传并创建媒资。
//
// 路由：POST /api/v1/asset/upload/complete
// 鉴权：JWT Bearer
// 请求体：{ uploadId, assetTitle, catalogId?, type }
// 成功：201 AssetVO（分片合并上传 S3 后入库）
func (h *Handler) UploadComplete(c *gin.Context) {
	var req uploadCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "参数错误")
		return
	}
	result, err := h.upload.Complete(c.Request.Context(), UploadCompleteInput{
		UploadID:   req.UploadID,
		AssetTitle: req.AssetTitle,
		CatalogID:  req.CatalogID,
		Type:       req.Type,
	})
	if err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			_ = c.Error(biz)
			return
		}
		_ = c.Error(err)
		return
	}
	httpx.Created(c, result)
}
