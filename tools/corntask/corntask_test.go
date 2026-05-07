package corntask

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "1",
		DB:       11,
	})
	m.Run()
}

func TestTimer(t *testing.T) {

	tt := NewTimer(zap.NewNop().Sugar(), redisClient, "golib1122:test", 5, 100, callback)
	//

	for i := 0; i < 3; i++ {
		exist, err := tt.Add(context.Background(), strconv.Itoa(i), "fff", 3)
		fmt.Println(exist, err)
	}

	time.Sleep(time.Hour)
}

func TestEncodeDecode(t *testing.T) {
	tt := NewTimer(zap.NewNop().Sugar(), redisClient, "golib1122:test", 5, 100, callback)
	//shutdown.AddShutDownJobWithTimeout("timer", func() error {
	//	tt.Stop()
	//	return nil
	//}, time.Second*10)

	s := tt.encodeMember("shadingkey1111", "abcdefg_123123")
	fmt.Println(s)

	fmt.Println(tt.decodeMember(s))

	// 测试旧版
	fmt.Println(tt.decodeMember("shadingkey1111_abcdefg123123"))
}

func callback(shardingKey string, data string, triggerTime int64) error {
	println(fmt.Sprintf("shardingKey=%s, data=%s, triggerTime=%d", shardingKey, data, triggerTime))
	return nil
}
