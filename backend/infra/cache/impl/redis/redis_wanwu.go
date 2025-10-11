package redis

import (
	"github.com/coze-dev/coze-studio/backend/infra/contract/cache"
	"github.com/redis/go-redis/v9"
)

func NewWithRedisCli(rdb *redis.Client) cache.Cmdable {
	cache.SetDefaultNilError(redis.Nil)

	return &redisImpl{client: rdb}
}
