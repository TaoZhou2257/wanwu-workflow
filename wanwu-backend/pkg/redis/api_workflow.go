package redis

import (
	"context"
	"fmt"
)

const (
	_dbWorkflow = 1
)

var (
	_redisWorkflow *redisImpl
)

func InitWorkflow(ctx context.Context, cfg Config) error {
	if _redisWorkflow != nil {
		return fmt.Errorf("redis workflow client already init")
	}
	c, err := newClient(ctx, cfg, _dbWorkflow)
	if err != nil {
		return err
	}
	_redisWorkflow = c
	return nil
}

func Workflow() *redisImpl {
	return _redisWorkflow
}
