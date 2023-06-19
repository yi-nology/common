package corntask

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/common/utils/xlogger"
	"github.com/yi-nology/common/utils/xnet"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LOCAL_QUEUE_LENGTH = 100000
)

func NewTimer(log xlogger.Logger, redisClient *redis.Client, remoteQueueKey string, remoteQueueNum int, localQueueMaxDelayNum int, callBack func(shardingKey string, data string, triggerTime int64) error) *Timer {
	if redisClient == nil {
		panic("redisClient is nil")
	}

	t := &Timer{log: log, redisClient: redisClient, remoteQueueKey: remoteQueueKey, remoteQueueNum: remoteQueueNum, callBack: callBack, localQueueMaxDelayNum: localQueueMaxDelayNum}
	t.init()
	return t
}

type Timer struct {
	redisClient *redis.Client

	remoteQueueKey string
	remoteQueueNum int

	callBack func(shardingKey string, data string, triggerTime int64) error

	assignedRemoteQueueKeys []string
	clientID                string

	localQueue            chan *localQueueItem
	localQueueMaxDelayNum int

	isStop bool
	log    xlogger.Logger
}

func (t *Timer) init() {
	if t.remoteQueueNum <= 2 {
		t.remoteQueueNum = 2
	}
	if t.remoteQueueNum >= 128 {
		t.remoteQueueNum = 128
	}

	t.assignedRemoteQueueKeys = make([]string, 0)
	t.localQueue = make(chan *localQueueItem, LOCAL_QUEUE_LENGTH)

	localIP, _ := xnet.ServerIP()
	t.clientID = localIP + "_" + strconv.Itoa(int(time.Now().Unix()))

	t.cleanInvalidClientID()
	t.clientHeartbeat()

	t.remoteQueueAssign()

	localQueueExecutorNum := len(t.assignedRemoteQueueKeys) * 10
	if localQueueExecutorNum > 500 {
		localQueueExecutorNum = 500
	}
	for i := 0; i < localQueueExecutorNum; i++ {
		go t.localQueueExecutor(i)
	}

	t.consumeRemoteQueue()
}

func (t *Timer) Add(ctx context.Context, shardingKey string, data string, delaySeconds int) (bool, error) {
	triggerTime := time.Now().Unix() + int64(delaySeconds)

	num, err := t.redisClient.ZAdd(ctx, t.selectRemoteQueue(shardingKey),
		redis.Z{Score: float64(triggerTime), Member: t.encodeMember(shardingKey, data)},
	).Result()

	return num == 0, err
}

func (t *Timer) Del(ctx context.Context, shardingKey string, data string) error {
	return t.redisClient.ZRem(ctx, t.selectRemoteQueue(shardingKey),
		t.encodeMember(shardingKey, data)).Err()
}

func (t *Timer) Stop() {
	t.isStop = true
	t.stopClientHeartbeat()

	for {
		if len(t.localQueue) > 0 {
			time.Sleep(time.Millisecond * 10)
		} else {
			break
		}
	}
}

func (t *Timer) clientHeartbeat() {
	if t.isStop {
		return
	}
	defer time.AfterFunc(time.Millisecond*300, t.clientHeartbeat)

	err := t.redisClient.ZAdd(context.Background(), t.getClientKey(),
		redis.Z{Score: float64(time.Now().Unix()), Member: t.clientID}).Err()
	if err != nil {
		t.log.Errorf("timer=remoteQueueAssign, RemoteQueueKey=%s, err=%v", t.remoteQueueKey, err)
		return
	}
}

func (t *Timer) stopClientHeartbeat() {
	err := t.redisClient.ZRem(context.Background(), t.getClientKey(),
		redis.Z{Score: float64(time.Now().Unix()), Member: t.clientID}).Err()
	if err != nil {
		t.log.Errorf("timer=Stop, RemoteQueueKey=%s, err=%v", t.remoteQueueKey, err)
		return
	}
}

func (t *Timer) cleanInvalidClientID() {
	err := t.redisClient.ZRemRangeByScore(context.Background(),
		t.getClientKey(), "0", strconv.Itoa(int(time.Now().Unix()-3600))).Err()
	if err != nil {
		t.log.Errorf("timer=cleanInvalidClientID, RemoteQueueKey=%s, err=%v", t.remoteQueueKey, err)
	}
}

func (t *Timer) delWithTriggerTime(selectedRemoteQueueKey string, shardingKey string, data string, triggerTime int64) error {
	script := redis.NewScript(`
if redis.call('zscore', KEYS[1], KEYS[2]) == ARGV[1] then
	return redis.call('zrem', KEYS[1], KEYS[2])
	else
	return 0
end
`)
	return script.Run(context.Background(), t.redisClient,
		[]string{selectedRemoteQueueKey, t.encodeMember(shardingKey, data)},
		triggerTime).Err()
}

func (t *Timer) remoteQueueAssign() {
	if t.isStop {
		return
	}
	defer time.AfterFunc(time.Second, t.remoteQueueAssign)

	curTime := time.Now().Unix()

	clients, err := t.redisClient.ZRangeByScore(context.Background(),
		t.getClientKey(),
		&redis.ZRangeBy{Min: strconv.Itoa(int(curTime - 1)), Max: strconv.Itoa(int(curTime))}).
		Result()
	if err != nil {
		t.log.Errorf("timer=remoteQueueAssign, RemoteQueueKey=%s, err=%v", t.remoteQueueKey, err)
		return
	}
	if len(clients) == 0 {
		return
	}

	sort.Strings(clients)
	idx := 0
	for _, one := range clients {
		if one == t.clientID {
			break
		}

		idx++
	}

	assignedKeys := make([]string, 0)
	for i := 0; i < t.remoteQueueNum; i++ {
		if i%len(clients) == idx {
			assignedKeys = append(assignedKeys, t.getRemoteKey(i))
		}
	}

	t.assignedRemoteQueueKeys = assignedKeys
}

func (t *Timer) consumeRemoteQueue() {
	if t.isStop {
		return
	}

	defer time.AfterFunc(time.Millisecond*500, t.consumeRemoteQueue)

	ctx := context.Background()

	var wg sync.WaitGroup
	for _, key := range t.assignedRemoteQueueKeys {
		wg.Add(1)

		go func(queueKey string) {
			defer wg.Done()

			offset, step := int64(0), int64(50)
			curTime := time.Now().Unix()

			for {
				members, err := t.redisClient.ZRangeByScoreWithScores(ctx, queueKey, &redis.ZRangeBy{
					Min:    "0",
					Max:    strconv.Itoa(int(curTime)),
					Offset: offset,
					Count:  step,
				}).Result()
				if err != nil {
					t.log.Errorf("timer=consumeRemoteQueue; queueKey=%s,err=%v", queueKey, members)
					break
				}

				for _, one := range members {
					if len(t.localQueue) > t.localQueueMaxDelayNum {
						t.log.Errorf("timer=consumeRemoteQueue, queueKey=%s, localQueue > 1/2", t.remoteQueueKey)
					}

					shardingKey, data := t.decodeMember(one.Member)
					t.localQueue <- &localQueueItem{
						selectedRemoteQueueKey: queueKey,
						shardingKey:            shardingKey,
						data:                   data,
						triggerTime:            int64(one.Score),
					}
				}

				if len(members) < int(step) {
					break
				}
				offset += step
			}
		}(key)
	}

	wg.Wait()
}

func (t *Timer) localQueueExecutor(workID int) {
	for item := range t.localQueue {
		t.log.Debugf("timer=localQueueExecutor, remoteQueueKey=%s, workID=%d,item=%+v",
			t.remoteQueueKey, workID, item)

		err := t.callBack(item.shardingKey, item.data, item.triggerTime)
		if err != nil {
			t.log.Debugf("timer=localQueueExecutor, remoteQueueKey=%s, workID=%d,item=%+v，err=%+v",
				t.remoteQueueKey, workID, item, err)
			continue
		}

		err = t.delWithTriggerTime(item.selectedRemoteQueueKey, item.shardingKey, item.data, item.triggerTime)
		if err != nil {
			t.log.Errorf("queueKey=%s,shardingKey=%s, data=%s",
				t.remoteQueueKey, item.shardingKey, item.data)
		}
	}
}

func (t *Timer) selectRemoteQueue(shardingKey string) string {
	h := fnv.New32a()
	h.Write([]byte(shardingKey))

	idx := int(h.Sum32()) % t.remoteQueueNum
	return t.getRemoteKey(idx)
}

func (t *Timer) getRemoteKey(idx int) string {
	return t.remoteQueueKey + "_" + strconv.Itoa(idx)
}

func (t *Timer) getClientKey() string {
	return t.remoteQueueKey + "_client"
}

func (t *Timer) encodeMember(shardingKey string, data string) string {
	return shardingKey + "@#@" + data
}

func (t *Timer) decodeMember(member interface{}) (string, string) {
	val, ok := member.(string)
	if !ok {
		return "", ""
	}

	data := strings.Split(val, "@#@")
	if len(data) != 2 {
		return "", ""
	}

	return data[0], data[1]
}

type localQueueItem struct {
	selectedRemoteQueueKey string
	shardingKey            string
	data                   string
	triggerTime            int64
}
