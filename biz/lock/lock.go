package lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/common/utils/xlogger"
	"time"
)

type Lock struct {
	rds redis.Client
	log xlogger.Logger
}

func NewLock(rds redis.Client, log xlogger.Logger) *Lock {
	return &Lock{rds: rds, log: log}
}

func (l *Lock) Up(ctx context.Context, key string, aliveSeconds int64) bool {
	success, err := l.rds.SetNX(ctx, key, 1, time.Duration(aliveSeconds)*time.Second).Result()
	if err != nil && err != redis.Nil {
		l.log.Errorf("Lock error key =%v, err =%v ", key, err)
		return false
	}
	return success
}

func (l *Lock) Down(ctx context.Context, key string) {
	err := l.rds.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		l.log.Errorf("UnLock error key =%v, err =%v ", key, err)
	}
}
