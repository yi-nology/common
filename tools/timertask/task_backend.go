package timertask

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/common/utils/xlogger"
	"time"
)

const TIME_TASK_DAY_COUNT = "queue.data.%s"

// 定时任务处理
type ProcessTimerTaskFunc = func(ctx context.Context, taskInfo *TaskInfo) error

type TaskBackend struct {
	TimeTaskManager  *TimeTaskManager
	processTimerFunc ProcessTimerTaskFunc
	log              xlogger.Logger
}

func NewTaskBackend(buzType string, process ProcessTimerTaskFunc, client *redis.Client, log xlogger.Logger) *TaskBackend {
	return &TaskBackend{
		TimeTaskManager:  NewTimeTaskManager(buzType, client, log),
		processTimerFunc: process,
		log:              log,
	}
}

func (m *TaskBackend) StartBackend(ctx context.Context) error {
	closech := make(chan struct{})
	//go m.falconQueueUpload(ctx)
	// 获取所有到期任务 放到 队列中
	go m.TaskBackendCheck(ctx)

	// 队列处理到期任务
	go m.TaskWorkBeckend(ctx, closech)

	return nil
}

func (m *TaskBackend) StartBackendByCheckAndWork(ctx context.Context, check int, work int) error {
	closech := make(chan struct{})
	//go m.falconQueueUpload(ctx)
	if check <= 0 || work <= 0 {
		return errors.New("check is >=1 and work >=1")
	}
	// 获取所有到期任务 放到 队列中
	for i := 0; i < check; i++ {
		go m.TaskBackendCheck(ctx)
	}

	// 队列处理到期任务
	for i := 0; i < work; i++ {
		go m.TaskWorkBeckend(ctx, closech)
	}
	return nil
}

func (m *TaskBackend) StartBackendBatch(ctx context.Context) error {
	closech := make(chan struct{})
	//go m.falconQueueUpload(ctx)
	// 获取所有到期任务 放到 队列中
	go m.TaskBackendCheckBatch(ctx)

	// 队列处理到期任务
	go m.TaskWorkBeckend(ctx, closech)

	return nil
}
func (m *TaskBackend) StartBackendBatchByCheckAndWork(ctx context.Context, check int, work int) error {
	closech := make(chan struct{})
	//go m.falconQueueUpload(ctx)
	// 获取所有到期任务 放到 队列中
	if check <= 0 || work <= 0 {
		return errors.New("check is >=1 and work >=1")
	}
	for i := 0; i < check; i++ {
		go m.TaskBackendCheckBatch(ctx)
	}
	// 队列处理到期任务
	for i := 0; i < work; i++ {
		go m.TaskWorkBeckend(ctx, closech)
	}
	return nil
}

func (m *TaskBackend) TaskBackendCheck(ctx context.Context) error {
	m.log.Debugf("TimeTaskManager|TaskBackendCheck")
	for range time.Tick(1 * time.Second) {
		shard := 1
		// 获取所有到期任务
		tasks, err := m.TimeTaskManager.TaskPool.GetAllTask(ctx, shard)
		if err != nil {
			m.log.Errorf("TimeTaskManager|TaskBackendCheck error %s", err)
			continue
		}
		for _, task := range tasks {
			succ, err := m.TimeTaskManager.TaskPool.DelTask(ctx, shard, task)
			if err != nil {
				m.log.Errorf("TimeTaskManager|TaskBackendCheck|DelTask %+v error %", task, err)
				continue
			}
			m.log.Debugf("TimeTaskManager|TaskBackendCheck|DelTask %+v succ %+v", task, succ)
			if succ {
				m.log.Debugf("TimeTaskManager|TaskBackendCheck|SendTaskQueue task %+v", task)
				m.TimeTaskManager.TaskPool.SendTaskQueue(ctx, task)
			}
		}
	}
	return nil
}

// 每次只获取到期的1000个
func (m *TaskBackend) TaskBackendCheckBatch(ctx context.Context) error {
	m.log.Debugf("TimeTaskManager|TaskBackendCheckBatch")
	for {
		for i := 0; i < _taskBucket; i++ {
			m.GetTaskTaskWithShard(ctx, i)
			<-time.NewTimer(time.Second * 1).C
		}
	}
	return nil
}

func (m *TaskBackend) GetTaskTaskWithShard(ctx context.Context, shard int) {
	// 获取所有到期任务
	tasks, err := m.TimeTaskManager.TaskPool.GetTaskWithBatch(ctx, shard, 1000)
	if err != nil {
		m.log.Errorf("TimeTaskManager|TaskBackendCheckBatch error %s", err)
		return
	}
	for _, task := range tasks {
		succ, err := m.TimeTaskManager.TaskPool.DelTask(ctx, shard, task)
		if err != nil {
			m.log.Errorf("TimeTaskManager|TaskBackendCheckBatch|DelTask %+v error %", task, err)
			continue
		}
		m.log.Debugf("TimeTaskManager|TaskBackendCheckBatch|DelTask %+v succ %+v", task, succ)
		if succ {
			m.log.Debugf("TimeTaskManager|TaskBackendCheckBatch|SendTaskQueue task %+v", task)
			m.TimeTaskManager.TaskPool.SendTaskQueue(ctx, task)
		}
	}
	return
}

func (m *TaskBackend) TaskWorkBeckend(ctx context.Context, closech chan struct{}) error {
	m.log.Debugf("TimeTaskManager|TaskWrokBackend")
	for v := range m.TimeTaskManager.TaskPool.GetTaskQueue(ctx, closech) {
		m.log.Debugf("TimeTaskManager|TaskWorkBackend|Do Work v %+v", v)
		task := new(TaskInfo)
		err := json.Unmarshal(v, task)
		if err != nil {
			m.log.Errorf("TimeTaskManager|TaskWordkBeckend|Unamrshal %s error %s", string(v), err)
			continue
		}
		m.log.Debugf("TimeTaskManager|TaskWorkBackend|Do Work %+v", task)

		m.processTimerFunc(ctx, task)
	}
	return nil
}

//func (m *TaskBackend) falconQueueUpload(ctx context.Context) {
//	for range time.Tick(1 * time.Minute) {
//		count, _ := m.TimeTaskManager.TaskPool.GetLen(ctx)
//		metrics.Meter(fmt.Sprintf(TIME_TASK_DAY_COUNT, m.TimeTaskManager.TaskPool.GetTaskName(ctx)), int(count))
//	}
//}
