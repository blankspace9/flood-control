package flood

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisFloodControl struct {
	redisClient *redis.Client
	key         string
	interval    time.Duration
	maxCalls    int
}

func NewRedisFloodControl(addr, key string, interval time.Duration, maxCall int) *RedisFloodControl {
	return &RedisFloodControl{
		redisClient: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		key:      key,
		interval: interval,
		maxCalls: maxCall,
	}
}

func (rfc *RedisFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	return true, nil
}
