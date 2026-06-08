package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker interface {
	TryLock(ctx context.Context, key string, ttl time.Duration) (token string, ok bool, err error)
	Unlock(ctx context.Context, key string, token string) error
}

type redisLocker struct {
	rdb *redis.Client
}

func NewRedisLocker(rdb *redis.Client) Locker {
	return &redisLocker{rdb: rdb}
}

var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`)

func (l *redisLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
	token, err := randomToken()
	if err != nil {
		return "", false, err
	}
	ok, err := l.rdb.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return token, true, nil
}

func (l *redisLocker) Unlock(ctx context.Context, key string, token string) error {
	if key == "" || token == "" {
		return errors.New("lock: empty key or token")
	}
	return unlockScript.Run(ctx, l.rdb, []string{key}, token).Err()
}

func randomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func InstanceLockKey(instanceID string) string {
	return "wf:lock:" + instanceID
}
