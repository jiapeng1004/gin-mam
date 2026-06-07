package asset

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

type Handler struct {
	svc    Service
	upload UploadService
}

func NewHandler(svc Service, upload UploadService) *Handler {
	return &Handler{svc: svc, upload: upload}
}

type pageRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
}

func (h *Handler) Page(c *gin.Context) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	httpx.OK(c, result)
}

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

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

type updateRequest struct {
	Title       *string           `json:"title"`
	CatalogID   *string           `json:"catalogId"`
	Description *string           `json:"description"`
	Status      *int8             `json:"status"`
	Metadata    map[string]string `json:"metadata"`
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	httpx.NoContent(c)
}
type uploadInitRequest struct {
	FileName  string `json:"fileName" binding:"required"`
	FileSize  int64  `json:"fileSize" binding:"required"`
	MimeType  string `json:"mimeType"`
	ChunkSize *int64 `json:"chunkSize"`
}

func (h *Handler) UploadInit(c *gin.Context) {
	var req uploadInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

func (h *Handler) UploadChunk(c *gin.Context) {
	uploadID := c.PostForm("uploadId")
	chunkIndexStr := c.PostForm("chunkIndex")
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || uploadID == "" {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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

type uploadCompleteRequest struct {
	UploadID    string  `json:"uploadId" binding:"required"`
	AssetTitle  string  `json:"assetTitle" binding:"required"`
	CatalogID   *string `json:"catalogId"`
	Type        string  `json:"type" binding:"required"`
}

func (h *Handler) UploadComplete(c *gin.Context) {
	var req uploadCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40000, "\u53c2\u6570\u9519\u8bef")
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
