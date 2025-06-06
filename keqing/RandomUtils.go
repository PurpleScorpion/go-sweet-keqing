package keqing

import (
	"errors"
	"hash/fnv"
	"math/rand"
	"net"
	"sync"
	"time"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu  sync.RWMutex
)

// WinningRate 随机概率
// num: 0~100之间的数，表示概率,可以为小数 , 最高支持0.000001%精度
// 当num >= 100时，一定会中奖 , 当num <= 0时 , 一定不会中奖
// 返回值: true表示命中，false表示未命中
func WinningRate(num float64) bool {
	if num >= 100 {
		return true
	}

	if num <= 0 {
		return false
	}

	const scale = 1_0000_0000           // 百万分之一精度（1e8），即支持 0.000001% 精度
	threshold := int(num * scale / 100) // 转换为百万分比的整数阈值
	if threshold < 0 {
		threshold = 0
	}
	if threshold > scale {
		threshold = scale
	}

	n := RandomNextInt(scale)
	return n < threshold
}

// RandomId 根据雪花算法生成id
func RandomId() int64 {
	if sf == nil {
		tmp, err := newSnowflake()
		if err == nil {
			sf = tmp
			id, err2 := sf.nextID()
			if err2 == nil {
				return id
			}
			return 0
		}
		return 0
	}
	id, err2 := sf.nextID()
	if err2 == nil {
		return id
	}
	return 0
}

const (
	workerBits     = 10
	sequenceBits   = 12
	maxWorkerId    = -1 ^ (-1 << workerBits)   // 支持的最大机器ID
	maxSequence    = -1 ^ (-1 << sequenceBits) // 同一毫秒内支持的最大序列数量
	workerShift    = sequenceBits
	timestampShift = sequenceBits + workerBits
	sequenceMask   = maxSequence
)

var (
	errTimeBackwards               = errors.New("时钟回拨")
	errSequenceOverflow            = errors.New("序列号超过最大值")
	errWorkerIDInvalid             = errors.New("无效的机器ID")
	sf                  *snowflake = nil
)

type snowflake struct {
	sync.Mutex
	workerId      int64
	lastTimestamp int64
	sequence      int64
}

// NextID 生成下一个ID
func (s *snowflake) nextID() (int64, error) {
	s.Lock()
	defer s.Unlock()

	timestamp := time.Now().UnixNano() / 1e6 // 毫秒级时间戳

	if timestamp < s.lastTimestamp {
		return 0, errTimeBackwards
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			timestamp = tilNextMillis(s.lastTimestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return (timestamp << timestampShift) |
		(s.workerId << workerShift) |
		s.sequence, nil
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixNano() / 1e6
	for timestamp <= lastTimestamp {
		timestamp = time.Now().UnixNano() / 1e6
	}
	return timestamp
}

func newSnowflake() (*snowflake, error) {
	workerId, err := getDefaultWorkerID()
	if err != nil {
		return nil, err
	}
	if workerId < 0 || workerId > maxWorkerId {
		return nil, errWorkerIDInvalid
	}
	return &snowflake{
		workerId:      workerId,
		lastTimestamp: -1,
		sequence:      0,
	}, nil
}

func getDefaultWorkerID() (int64, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}

	for _, i := range interfaces {
		if i.HardwareAddr != nil {
			// 使用 MAC 地址哈希计算 workerId
			hwAddrStr := i.HardwareAddr.String()
			h := fnv.New32a()
			_, _ = h.Write([]byte(hwAddrStr))
			return int64(h.Sum32() % 1024), nil // 限制在 10 位以内（0~1023）
		}
	}
	return 0, errors.New("no MAC address found")
}

/*
RandomNextInt
1个参数: 0~该参数随机 , [0,arg)
2个参数: [arg1,arg2)
*/
func RandomNextInt(args ...interface{}) int {
	if IsEmpty(args) {
		panic("Parameter cannot be empty")
	}
	parsedArgs, ok := parseArgsToInts(args...)
	if !ok {
		return 0
	}
	switch len(parsedArgs) {
	case 1:
		return nextInt1(parsedArgs[0])
	case 2:
		return nextInt2(parsedArgs[0], parsedArgs[1])
	default:
		return 0
	}
}

// Random 生成0-1之间的随机数 [0,1)
func Random() float64 {
	return rng.Float64()
}

func parseArgsToInts(args ...interface{}) ([]int, bool) {
	var result []int
	for _, arg := range args {
		if f, v := isInt(arg); f {
			result = append(result, v)
		} else {
			return nil, false
		}
	}
	return result, true
}

func isInt(v interface{}) (bool, int) {
	switch val := v.(type) {
	case int:
		return true, val
	case int8:
		return true, int(val)
	case int16:
		return true, int(val)
	case int32:
		return true, int(val)
	case int64:
		return true, int(val)
	default:
		return false, 0
	}
}

func nextInt2(min, max int) int {
	mu.RLock()
	defer mu.RUnlock()
	return rng.Intn(max-min) + min
}

func nextInt1(num int) int {
	mu.RLock()
	defer mu.RUnlock()
	return rng.Intn(num)
}
