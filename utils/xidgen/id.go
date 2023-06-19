package xidgen

import (
	"errors"
	"fmt"
	"github.com/yi-nology/common/utils/xlogger"
	"hash/fnv"
	"net"
	"os"
	"sync"
	"time"
)

///////////////////////////////
// 分布式int64唯一id
// |-- 1位符号位 --|-- 41位毫秒时间戳 --|-- 10位机器id --|-- 12位自增序列id --|
//
///////////////////////////////

const (
	_timeStart  = 1577808000000 // 2020-01-01 00:00:00
	_seqIDMask  = 4095          // 12位掩码
	_workIDMask = 1023          // 10位掩码
	_timeMask   = 2199023255551 // 41位掩码
)

// 唯一ID生成
type SequenceID int64

type IDGenerator struct {
	ip     string
	seqId  int64
	mutex  sync.Mutex
	logger xlogger.Logger
}

func NewIDGenerator(logger xlogger.Logger) *IDGenerator {
	ip, err := GetIP()
	if err != nil {
		fmt.Printf("IDGenerator|get ip error. err:%v\n", err)
		os.Exit(1)
	}

	logger.Infof("IDGenerator|NewIDGenerator ip:%s", ip)

	return &IDGenerator{
		ip:    ip,
		seqId: 0,
	}
}

func (g *IDGenerator) GenID() (SequenceID, error) {
	// 毫秒时间戳:41
	timestamp := (time.Now().UnixNano() / 1000000)
	tid := (timestamp - _timeStart) & _timeMask
	// 机器id:10位
	workId, err := g.getWorkID()
	if err != nil {
		return 0, err
	}
	// 自增id:12位
	seqId := g.getSeqID()

	id := (tid << 22) | (workId << 12) | (seqId)

	g.logger.Infof("IDGenerator|id:%d, timestamp:%d, tid:%d, workId:%d, seqId:%d", id, timestamp, tid, workId, seqId)

	return SequenceID(id), nil
}

func (g *IDGenerator) getWorkID() (int64, error) {
	h, err := HashIP(g.ip)
	if err != nil {
		return 0, err
	}

	wid := h & _workIDMask
	return wid, nil
}

func (g *IDGenerator) getSeqID() int64 {
	g.mutex.Lock()

	defer g.mutex.Unlock()

	seqId := (g.seqId + 1) & _seqIDMask
	g.seqId = seqId

	return seqId
}

func GetIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("get ip error")
}

func HashIP(ip string) (int64, error) {
	h := fnv.New32a()
	if _, err := h.Write([]byte(ip)); err != nil {
		return 0, err
	}
	a := h.Sum32()
	return int64(a), nil
}

func (id SequenceID) GetTime() time.Time {
	timestamp := (int64(id) >> 22) + _timeStart
	t := time.Unix(timestamp/1000, 0)
	return t
}

func (id SequenceID) GetWorkID() int64 {
	wid := (int64(id) >> 12) & _workIDMask
	return wid
}

func (id SequenceID) GetSeqID() int64 {
	seqId := int64(id) & _seqIDMask
	return seqId
}
