package lock

import (
	"context"
	"errors"
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
	if err != nil && !errors.Is(err, redis.Nil) {
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

func (l *Lock) UpWait(ctx context.Context, function string, key string, aliveSeconds int64) error {
	for {
		success := l.Up(ctx, key, aliveSeconds)
		if success {
			l.log.Infof("Tag=%v; func=%v; key=%+v; info=%s;", "UpWait", function, key, "加锁成功")
			return nil
		} else {
			l.log.Infof("Tag=%v; func=%v; key=%+v; info=%s;", "UpWait", function, key, "有业务在执行,休眠100ms")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
