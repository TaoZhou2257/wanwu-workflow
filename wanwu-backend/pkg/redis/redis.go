package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host       string `json:"host" mapstructure:"host"`
	Port       string `json:"port" mapstructure:"port"`
	Username   string `json:"username" mapstructure:"username"`
	Password   string `json:"password" mapstructure:"password"`
	Channel    string `json:"channel" mapstructure:"channel"`
	Standalone bool   `json:"standalone" mapstructure:"standalone"`
	MasterName string `json:"master_name" mapstructure:"master_name"`
}

type redisImpl struct {
	client *redis.Client
}

func newClient(ctx context.Context, c Config, db int) (*redisImpl, error) {
	redisAddr := c.Host + ":" + c.Port
	var r *redis.Client
	if c.Standalone {
		r = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Username: c.Username,
			Password: c.Password,
			DB:       db,
		})
	} else {
		r = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       c.MasterName,
			SentinelAddrs:    []string{redisAddr}, // 哨兵节点地址
			DB:               db,                  // 使用的Redis数据库编号
			SentinelPassword: c.Password,
			Password:         c.Password,
		})
	}
	if _, err := r.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return &redisImpl{
		client: r,
	}, nil
}
