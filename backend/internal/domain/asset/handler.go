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

// pageRequest 媒资分页请求体。
type pageRequest struct {
	Page     int    `json:"page" example:"1"`
	PageSize int    `json:"pageSize" example:"20"`
	Keyword  string `json:"keyword" example:""`
}

// Page 分页查询媒资
//
//	@Summary		分页查询媒资
//	@Tags			媒资
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		pageRequest	true	"分页参数"
//	@Success		200		{object}	PageResult
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/page [post]
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

// Get 获取媒资详情
//
//	@Summary		获取媒资详情
//	@Tags			媒资
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"媒资 ID"
//	@Success		200	{object}	AssetVO
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

// createRequest 创建媒资请求体。
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

// Create 创建媒资
//
//	@Summary		创建媒资
//	@Description	直接创建媒资（已有 storagePath，非分片上传流程）
//	@Tags			媒资
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createRequest	true	"媒资参数"
//	@Success		201		{object}	AssetVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset [post]
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

// updateRequest 更新媒资请求体。
type updateRequest struct {
	Title       *string           `json:"title"`
	CatalogID   *string           `json:"catalogId"`
	Description *string           `json:"description"`
	Status      *int8             `json:"status"`
	Metadata    map[string]string `json:"metadata"`
}

// Update 更新媒资
//
//	@Summary		更新媒资
//	@Tags			媒资
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string			true	"媒资 ID"
//	@Param			body	body		updateRequest	true	"更新参数"
//	@Success		200		{object}	AssetVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Failure		404		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/{id} [put]
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

// Delete 删除媒资
//
//	@Summary		删除媒资
//	@Description	软删除媒资（GORM DeletedAt）
//	@Tags			媒资
//	@Security		BearerAuth
//	@Param			id	path	string	true	"媒资 ID"
//	@Success		204
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}

// uploadInitRequest 初始化上传请求体。
type uploadInitRequest struct {
	FileName  string `json:"fileName" binding:"required"`
	FileSize  int64  `json:"fileSize" binding:"required"`
	MimeType  string `json:"mimeType"`
	ChunkSize *int64 `json:"chunkSize"`
}

// UploadInit 初始化分片上传
//
//	@Summary		初始化分片上传
//	@Tags			媒资上传
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		uploadInitRequest	true	"文件信息"
//	@Success		200		{object}	UploadInitResult
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/upload/init [post]
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

// UploadChunk 上传分片
//
//	@Summary		上传分片
//	@Tags			媒资上传
//	@Accept			multipart/form-data
//	@Security		BearerAuth
//	@Param			uploadId	formData	string	true	"上传会话 ID"
//	@Param			chunkIndex	formData	int		true	"分片序号（从 0 开始）"
//	@Param			file		formData	file	true	"分片二进制"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorVo
//	@Failure		401	{object}	httpx.ErrorVo
//	@Failure		404	{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/upload/chunk [post]
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

// uploadCompleteRequest 完成上传请求体。
type uploadCompleteRequest struct {
	UploadID   string  `json:"uploadId" binding:"required"`
	AssetTitle string  `json:"assetTitle" binding:"required"`
	CatalogID  *string `json:"catalogId"`
	Type       string  `json:"type" binding:"required"`
}

// UploadComplete 完成分片上传
//
//	@Summary		完成分片上传
//	@Description	合并分片、上传对象存储并创建媒资记录
//	@Tags			媒资上传
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		uploadCompleteRequest	true	"完成参数"
//	@Success		201		{object}	AssetVO
//	@Failure		400		{object}	httpx.ErrorVo
//	@Failure		401		{object}	httpx.ErrorVo
//	@Router			/api/v1/asset/upload/complete [post]
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
