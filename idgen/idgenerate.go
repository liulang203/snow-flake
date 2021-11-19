package idgen

import (
	"errors"
	"net"
	"sync"
	"time"
)

/**
雪花算法
	 41bit timestamp | 12 bit machineID : 8bit workerID 4bit dataCenterID ｜ 10 bit sequenceBits
最多使用69年
*/

const (
	workerIDBits     = uint64(8) // 10bit 工作机器ID中的 5bit workerID
	dataCenterIDBits = uint64(4) // 10 bit 工作机器ID中的 5bit dataCenterID
	sequenceBits     = uint64(10)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits) //节点ID的最大值 用于防止溢出
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits // 时间戳向左偏移量
	dataLeft = uint8(14) // dataLeft = dataCenterIDBits + sequenceBits
	workLeft = uint8(10) // workLeft = sequenceBits // 节点IDx向左偏移量
	// 2021-11-18 08:00:00 +0000 CST
	twepoch = int64(1637193600000) // 常量时间戳(毫秒)
)

type Worker struct {
	mu           sync.Mutex
	LastStamp    int64 // 记录上一次ID的时间戳
	WorkerID     int64 // 该节点的ID
	DataCenterID int64 // 该节点的 数据中心ID
	Sequence     int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4096个ID
}

//分布式情况下,我们应通过外部配置文件或其他方式为每台机器分配独立的id
func NewWorker(workerID, dataCenterID int64) *Worker {

	return &Worker{
		WorkerID:     workerID & maxWorkerID,
		LastStamp:    0,
		Sequence:     0,
		DataCenterID: dataCenterID & maxDataCenterID,
	}
}
func getClientIp() (byte, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return 0, err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP[0], nil
			}

		}
	}
	return 0, errors.New("Can not find the client ip address!")
}

func (w *Worker) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *Worker) NextID() (int64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.nextID()
}

func (w *Worker) nextID() (int64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}

	if w.LastStamp == timeStamp {

		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) |
		(w.DataCenterID << dataLeft) |
		(w.WorkerID << workLeft) |
		w.Sequence

	return id, nil
}
