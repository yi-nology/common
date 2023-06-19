package timertask

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yi-nology/common/utils/xlogger"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_taskRedis      = "task.redis"
	TIME_TASK       = "time.honey.task.%s.%d"
	TIME_TASK_SHARD = 1
	TIME_TASK_QUEUE = "time.honey.task.%s.queue"
	_taskBucket     = 128
)

var ErrTimeout = errors.New("redis: i/o timeout, please retry")

type TaskInfo struct {
	TaskID     string `json:"task_id"`     // 任务id
	CreateTime int64  `json:"create_time"` // 开始时间
	EndTime    int64  `json:"end_time"`    // 结束时间
	Action     string `json:"action"`      // 结束执行的函数映射
	Param      string `json:"param"`       // 执行需要的参数
	Extra      string `json:"extra"`       // 额外字段
}

type TaskPool struct {
	taskCache *redis.Client
	buzType   string
	log       xlogger.Logger
}

func NewTaskPool(buzType string, client *redis.Client, log xlogger.Logger) *TaskPool {
	return &TaskPool{
		taskCache: client,
		buzType:   buzType,
		log:       log,
	}
}

func (t *TaskPool) getCurrentTaskKey() string {
	return fmt.Sprintf(TIME_TASK, t.buzType, TIME_TASK_SHARD)
}

func (t *TaskPool) getTaskKey(shard int) string {
	return fmt.Sprintf(TIME_TASK, t.buzType, shard)
}
func (t *TaskPool) getTaskQueue() string {
	return fmt.Sprintf(TIME_TASK_QUEUE, t.buzType)
}

func hashStr(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (t *TaskPool) getBucketTaskKey(taskID string) string {
	shard := TIME_TASK_SHARD
	if len(taskID) != 0 {
		shard = int(hashStr(taskID)) % _taskBucket
	}
	return fmt.Sprintf(TIME_TASK, t.buzType, shard)
}

// 根据taskId投入多个任务桶中
func (t *TaskPool) AddTimeTaskWithBucket(ctx context.Context, task *TaskInfo) error {
	key := t.getBucketTaskKey(task.TaskID)
	info, _ := json.Marshal(task)
	_, err := t.taskCache.ZAdd(ctx, key, redis.Z{
		Score:  float64(task.EndTime),
		Member: string(info),
	}).Result()
	if err != nil {
		t.log.Errorf("TaskPool|AddTimeTask|ZADD %s %+v error %s", key, info, err)
		return err
	}
	t.log.Debugf("TaskPool|AddTimeTask|ZADD ok. key:%s, task:%s", key, info)
	return nil
}

func (t *TaskPool) AddTimeTask(ctx context.Context, task *TaskInfo) error {
	key := t.getCurrentTaskKey()
	info, _ := json.Marshal(task)
	_, err := t.taskCache.ZAdd(ctx, key, redis.Z{
		Score:  float64(task.EndTime),
		Member: string(info),
	}).Result()
	if err != nil {
		t.log.Errorf("TaskPool|AddTimeTask|ZADD %s %+v error %s", key, info, err)
		return err
	}
	t.log.Debugf("TaskPool|AddTimeTask|ZADD ok. key:%s, task:%s", key, info)
	return nil
}

func (t *TaskPool) AddTimeTaskHoney(ctx context.Context, task *TaskInfo, endTime int64) error {
	key := t.getCurrentTaskKey()
	info, _ := json.Marshal(task)
	_, err := t.taskCache.ZAdd(ctx, key, redis.Z{
		Score:  float64(task.EndTime),
		Member: string(info),
	}).Result()
	if err != nil {
		t.log.Errorf("TaskPool|AddTimeTask|HSet %s %+v error %s", key, task, err)
		return err
	}
	return nil
}

func (t *TaskPool) GetAllTask(ctx context.Context, shard int) ([]string, error) {
	t.log.Debugf("TaskPool|GetAllTask|shard %d", shard)
	key := t.getTaskKey(shard)
	current := time.Now().Unix()
	reply, err := t.taskCache.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(current, 10),
		Offset: 0,
		Count:  300,
	}).Result()
	if err != nil {
		t.log.Errorf("TaskPool|GetAllTask|ZRANGEBYSCORE %s %d %d error %s", key, 0, current, err)
		return nil, err
	}

	return reply, err

}

func (t *TaskPool) GetTaskWithBatch(ctx context.Context, shard int, batch int) ([]string, error) {
	t.log.Debugf("TaskPool|GetTaskWithBatch|shard %d", shard)
	key := t.getTaskKey(shard)
	current := time.Now().Unix()
	reply, err := t.taskCache.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(current, 10),
		Offset: 0,
		Count:  int64(batch),
	}).Result()
	if err != nil {
		t.log.Errorf("TaskPool|GetTaskWithBatch|ZRANGEBYSCORE %s %d %d error %s", key, 0, current, err)
		return nil, err
	}

	return reply, err

}

func (t *TaskPool) DelTask(ctx context.Context, shard int, task string) (bool, error) {
	key := t.getTaskKey(shard)
	ret, err := t.taskCache.ZRem(ctx, key, task).Result()
	if err != nil {
		t.log.Errorf("TaskPool|DelTask|HDel %s %v error %s", key, task, err)
		return false, err
	}
	if ret == 1 {
		t.log.Debugf("TaskPool|DelTask|ZRem ok. key:%s, task:%s", key, task)
		return true, nil
	}
	t.log.Debugf("TaskPool|DelTask|ZRem fail. key:%s, task:%s", key, task)
	return false, nil
}

func (t *TaskPool) SendTaskQueue(ctx context.Context, task string) error {
	t.log.Debugf("TaskPool|SendCurrentTask|tast %+v", task)
	key := t.getTaskQueue()

	_, err := t.taskCache.RPush(ctx, key, task).Result()
	if err != nil {
		t.log.Errorf("TaskPool|SendTaskQueue|RPush %s %+v error %s", key, task, err)
		return err
	}
	t.log.Debugf("TaskPool|SendTaskQueue|RPush ok. key:%s, task:%s", key, task)
	return nil
}

func (t *TaskPool) GetTaskQueue(ctx context.Context, closech chan struct{}) chan []byte {
	key := t.getTaskQueue()
	ch := t.Receive(ctx, key, closech, 2)
	return ch
}

func (t *TaskPool) GetLen(ctx context.Context) (int64, error) {
	return t.taskCache.LLen(ctx, t.getTaskQueue()).Result()
}

func (t *TaskPool) GetTaskName(ctx context.Context) string {
	return t.buzType
}

func (t *TaskPool) Receive(ctx context.Context, name string, closech chan struct{}, bufferSize int) chan []byte {
	ch := make(chan []byte, bufferSize)
	go func() {
		defer close(ch)
		for {
			select {
			case <-closech:
				return
			default:
				data, err := t.taskCache.BLPop(ctx, time.Second, name).Result()
				if err == nil {
					if data != nil {
						ch <- []byte(data[1])
					}
				} else if err != ErrTimeout && err != redis.Nil {
					t.log.Errorf("BLPOP error %s", err)
				}
			}
		}
	}()
	return ch
}
