package asset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	"github.com/gin-mam/backend/internal/infra/storage"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/gin-mam/backend/internal/pkg/id"
	"github.com/redis/go-redis/v9"
)

const (
	// ConfigKeyUploadChunkSize 分片上传块大小（字节）配置键。
	ConfigKeyUploadChunkSize = "MAM_UPLOAD_CHUNK_SIZE"
	// DefaultUploadChunkSize 默认分片大小 5MB。
	DefaultUploadChunkSize = 5242880
	// uploadSessionTTL 上传会话在 Redis 中的过期时间。
	uploadSessionTTL = 24 * time.Hour
	// ErrCodeUploadNotFound 上传会话不存在时的业务错误码。
	ErrCodeUploadNotFound = 40402
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=upload.go -destination=mock/upload_mock.go -package=mock

// UploadService 媒资分片上传应用服务接口。
// 流程：Init 创建会话 → SaveChunk 逐片写入 Redis → Complete 合并上传 S3 并创建媒资。
type UploadService interface {
	// Init 初始化上传会话，返回 uploadId 与建议分片大小。
	Init(ctx context.Context, in UploadInitInput) (*UploadInitResult, error)
	// SaveChunk 保存单个分片二进制数据至 Redis。
	SaveChunk(ctx context.Context, in UploadChunkInput) error
	// Complete 合并分片、上传对象存储并创建媒资记录。
	Complete(ctx context.Context, in UploadCompleteInput) (*AssetVO, error)
}

// UploadInitInput 初始化上传入参。
type UploadInitInput struct {
	FileName  string
	FileSize  int64
	MimeType  string
	ChunkSize *int64
}

// UploadInitResult 初始化上传响应。
type UploadInitResult struct {
	UploadID  string `json:"uploadId"`
	ChunkSize int64  `json:"chunkSize"`
}

// UploadChunkInput 上传单个分片入参。
type UploadChunkInput struct {
	UploadID   string
	ChunkIndex int
	Body       io.Reader
}

// UploadCompleteInput 完成上传并入参。
type UploadCompleteInput struct {
	UploadID   string
	AssetTitle string
	CatalogID  *string
	Type       string
	CreatedBy  string
}

// UploadSessionStore 上传会话持久化接口，默认 Redis 实现。
type UploadSessionStore interface {
	// SaveMeta 保存上传会话元信息。
	SaveMeta(ctx context.Context, uploadID string, meta UploadMeta) error
	// GetMeta 读取上传会话元信息。
	GetMeta(ctx context.Context, uploadID string) (*UploadMeta, error)
	// SaveChunk 保存指定序号的分片数据。
	SaveChunk(ctx context.Context, uploadID string, index int, data []byte) error
	// GetChunk 读取指定序号的分片数据。
	GetChunk(ctx context.Context, uploadID string, index int) ([]byte, error)
	// DeleteSession 删除上传会话及全部分片。
	DeleteSession(ctx context.Context, uploadID string) error
}

// UploadMeta 上传会话元信息，存 Redis JSON。
type UploadMeta struct {
	FileName  string `json:"fileName"`
	FileSize  int64  `json:"fileSize"`
	MimeType  string `json:"mimeType"`
	ChunkSize int64  `json:"chunkSize"`
}

type uploadService struct {
	store   UploadSessionStore
	config  infraconfig.Service
	storage storage.Storage
	assets  Service
}

func NewUploadService(store UploadSessionStore, config infraconfig.Service, stor storage.Storage, assets Service) UploadService {
	return &uploadService{store: store, config: config, storage: stor, assets: assets}
}

func UploadMetaKey(uploadID string) string {
	return fmt.Sprintf("upload:%s:meta", uploadID)
}

func uploadChunkKey(uploadID string, index int) string {
	return fmt.Sprintf("upload:%s:chunk:%d", uploadID, index)
}

type redisUploadStore struct {
	rdb redis.Cmdable
}

func NewRedisUploadStore(rdb redis.Cmdable) UploadSessionStore {
	return &redisUploadStore{rdb: rdb}
}

func (s *redisUploadStore) SaveMeta(ctx context.Context, uploadID string, meta UploadMeta) error {
	b, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, UploadMetaKey(uploadID), b, uploadSessionTTL).Err()
}

func (s *redisUploadStore) GetMeta(ctx context.Context, uploadID string) (*UploadMeta, error) {
	val, err := s.rdb.Get(ctx, UploadMetaKey(uploadID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, httpx.NewBizError(404, ErrCodeUploadNotFound, "\u4e0a\u4f20\u4f1a\u8bdd\u4e0d\u5b58\u5728")
		}
		return nil, err
	}
	var meta UploadMeta
	if err := json.Unmarshal(val, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func (s *redisUploadStore) SaveChunk(ctx context.Context, uploadID string, index int, data []byte) error {
	return s.rdb.Set(ctx, uploadChunkKey(uploadID, index), data, uploadSessionTTL).Err()
}

func (s *redisUploadStore) GetChunk(ctx context.Context, uploadID string, index int) ([]byte, error) {
	val, err := s.rdb.Get(ctx, uploadChunkKey(uploadID, index)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("chunk %d missing", index)
		}
		return nil, err
	}
	return val, nil
}

func (s *redisUploadStore) DeleteSession(ctx context.Context, uploadID string) error {
	pattern := fmt.Sprintf("upload:%s:chunk:*", uploadID)
	var cursor uint64
	var keys []string
	for {
		k, next, err := s.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, k...)
		cursor = next
		if cursor == 0 {
			break
		}
	}
	keys = append(keys, UploadMetaKey(uploadID))
	if len(keys) == 0 {
		return nil
	}
	return s.rdb.Del(ctx, keys...).Err()
}

func (s *uploadService) resolveChunkSize(ctx context.Context, override *int64) (int64, error) {
	if override != nil && *override > 0 {
		return *override, nil
	}
	chunkSize := int64(DefaultUploadChunkSize)
	val, err := s.config.Get(ctx, ConfigKeyUploadChunkSize)
	if err != nil {
		if errors.Is(err, infraconfig.ErrNotFound) {
			return chunkSize, nil
		}
		return 0, err
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
	if err != nil || parsed <= 0 {
		return chunkSize, nil
	}
	return parsed, nil
}

func uploadChunkTotal(fileSize, chunkSize int64) int {
	if chunkSize <= 0 {
		return 0
	}
	n := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		n++
	}
	return int(n)
}

func (s *uploadService) Init(ctx context.Context, in UploadInitInput) (*UploadInitResult, error) {
	fileName := strings.TrimSpace(in.FileName)
	if fileName == "" || in.FileSize <= 0 {
		return nil, httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	chunkSize, err := s.resolveChunkSize(ctx, in.ChunkSize)
	if err != nil {
		return nil, err
	}
	uploadID := id.New()
	meta := UploadMeta{
		FileName:  path.Base(fileName),
		FileSize:  in.FileSize,
		MimeType:  strings.TrimSpace(in.MimeType),
		ChunkSize: chunkSize,
	}
	if err := s.store.SaveMeta(ctx, uploadID, meta); err != nil {
		return nil, err
	}
	return &UploadInitResult{UploadID: uploadID, ChunkSize: chunkSize}, nil
}

func (s *uploadService) SaveChunk(ctx context.Context, in UploadChunkInput) error {
	uploadID := strings.TrimSpace(in.UploadID)
	if uploadID == "" || in.ChunkIndex < 0 {
		return httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	meta, err := s.store.GetMeta(ctx, uploadID)
	if err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			return biz
		}
		return err
	}
	total := uploadChunkTotal(meta.FileSize, meta.ChunkSize)
	if in.ChunkIndex >= total {
		return httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	data, err := io.ReadAll(in.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	expected := meta.ChunkSize
	if in.ChunkIndex == total-1 {
		if rem := meta.FileSize % meta.ChunkSize; rem != 0 {
			expected = rem
		}
	}
	if int64(len(data)) > expected {
		return httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	return s.store.SaveChunk(ctx, uploadID, in.ChunkIndex, data)
}

func (s *uploadService) Complete(ctx context.Context, in UploadCompleteInput) (*AssetVO, error) {
	uploadID := strings.TrimSpace(in.UploadID)
	title := strings.TrimSpace(in.AssetTitle)
	assetType := strings.TrimSpace(in.Type)
	if uploadID == "" || title == "" || assetType == "" {
		return nil, httpx.NewBizError(400, 40000, "\u53c2\u6570\u9519\u8bef")
	}
	meta, err := s.store.GetMeta(ctx, uploadID)
	if err != nil {
		if biz, ok := err.(*httpx.BizError); ok {
			return nil, biz
		}
		return nil, err
	}
	total := uploadChunkTotal(meta.FileSize, meta.ChunkSize)
	merged := make([]byte, 0, meta.FileSize)
	for i := 0; i < total; i++ {
		part, err := s.store.GetChunk(ctx, uploadID, i)
		if err != nil {
			return nil, httpx.NewBizError(400, 40000, "\u5206\u7247\u672a\u5b8c\u6574")
		}
		merged = append(merged, part...)
	}
	if int64(len(merged)) != meta.FileSize {
		return nil, httpx.NewBizError(400, 40000, "\u5206\u7247\u672a\u5b8c\u6574")
	}
	storageKey := fmt.Sprintf("uploads/%s/%s", uploadID, meta.FileName)
	if err := s.storage.Put(ctx, storageKey, bytes.NewReader(merged), meta.FileSize, meta.MimeType); err != nil {
		return nil, err
	}
	createdBy := strings.TrimSpace(in.CreatedBy)
	if createdBy == "" {
		createdBy = "system"
	}
	fileSize := meta.FileSize
	vo, err := s.assets.Create(ctx, CreateInput{
		Title:       title,
		Type:        assetType,
		CatalogID:   in.CatalogID,
		StoragePath: storageKey,
		MimeType:    meta.MimeType,
		FileSize:    &fileSize,
		CreatedBy:   createdBy,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, storageKey)
		return nil, err
	}
	_ = s.store.DeleteSession(ctx, uploadID)
	return vo, nil
}