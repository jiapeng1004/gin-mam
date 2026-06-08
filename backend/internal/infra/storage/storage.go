package storage

import (
	"context"
	"io"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=storage.go -destination=mock/storage_mock.go -package=mock

// Storage abstracts object storage (S3-compatible).
type Storage interface {
	Put(ctx context.Context, key string, body io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
	PresignPut(ctx context.Context, key string, expire time.Duration) (string, error)
}