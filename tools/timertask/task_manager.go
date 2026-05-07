package timertask

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/common/utils/xlogger"
)

type TimeTaskManager struct {
	TaskPool *TaskPool
	log      xlogger.Logger
}

func NewTimeTaskManager(buzType string, client *redis.Client, log xlogger.Logger) *TimeTaskManager {
	return &TimeTaskManager{
		TaskPool: NewTaskPool(buzType, client, log),
		log:      log,
	}
}

func (m *TimeTaskManager) SendTask(ctx context.Context, task *TaskInfo) error {
	err := m.TaskPool.AddTimeTask(ctx, task)
	if err != nil {
		m.log.Errorf("TimeTaskManager|SendTask|req %+v AddTimeTask error %s", task, err)
		return err
	}
	return nil
}

func (m *TimeTaskManager) SendTaskById(ctx context.Context, task *TaskInfo) error {
	err := m.TaskPool.AddTimeTaskWithBucket(ctx, task)
	if err != nil {
		m.log.Errorf("TimeTaskManager|SendTask|req %+v AddTimeTaskWithBucket error %s", task, err)
		return err
	}
	return nil
}

func (m *TimeTaskManager) SendTaskHoney(ctx context.Context, task *TaskInfo, endTime int64) error {
	m.log.Debugf("TimeTaskManager|SendTask|task %+v", task)
	err := m.TaskPool.AddTimeTaskHoney(ctx, task, endTime)
	if err != nil {
		m.log.Errorf("TimeTaskManager|SendTask|req %+v AddTimeTask error %s", task, err)
		return err
	}
	return nil
}
