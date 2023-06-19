package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/common/tools/timertask"
	"github.com/yi-nology/common/utils/xuuid"
	"go.uber.org/zap"
	"time"
)

type TaskExtra struct {
	Uid   int64 `json:"uid"`
	Times int64 `json:"times"`
}

var rds redis.Client

func init() {
	flag.Parse()
	rds = *redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "zhanghao123",
		DB:       5,
	})
}

func main() {
	logger, _ := zap.NewProduction()
	taskBck := timertask.NewTaskBackend("test", ProcessTask, &rds, logger.Sugar())

	// 添加任务
	go func() {
		for i := 0; i < 10; i = i + 1 {
			AddTask(context.Background(), logger.Sugar())
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(5 * time.Second)

	// 异步处理任务
	taskBck.StartBackend(context.Background())

	time.Sleep(30 * time.Second)
	taskBck.StartBackendByCheckAndWork(context.Background(), 2, 2)
	time.Sleep(100000 * time.Hour)
}

func AddTask(ctx context.Context, logger *zap.SugaredLogger) {
	taskMgr := timertask.NewTimeTaskManager("test", &rds, logger)

	taskExtra := &TaskExtra{
		Uid:   1001,
		Times: 10,
	}
	extra, _ := json.Marshal(taskExtra)

	info := &timertask.TaskInfo{
		TaskID:     xuuid.Uuid(),
		CreateTime: time.Now().Unix(),
		EndTime:    time.Now().Unix() + 2,
		Action:     "key1",
		Param:      "",
		Extra:      string(extra),
	}

	// 添加任务
	taskMgr.SendTask(context.Background(), info)
	fmt.Printf("AddTask|info:%+v\n", info)
}

func ProcessTask(ctx context.Context, task *timertask.TaskInfo) error {
	fmt.Printf("ProcessTask|task:%+v\n", task)

	switch task.Action {
	case "key1":
		extra := &TaskExtra{}
		_ = json.Unmarshal([]byte(task.Extra), extra)
		fmt.Printf("ProcessTask|do key1 task. extra:%+v\n", extra)
	default:
	}

	return nil
}
