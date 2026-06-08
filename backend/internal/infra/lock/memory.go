package lock

import (
	"context"
	"time"
)

type memoryLocker struct{}

func NewMemoryLocker() Locker {
	return &memoryLocker{}
}

func (memoryLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
	_ = ctx
	_ = key
	_ = ttl
	return "memory-token", true, nil
}

func (memoryLocker) Unlock(ctx context.Context, key string, token string) error {
	_ = ctx
	_ = key
	_ = token
	return nil
}
