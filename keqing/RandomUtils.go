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
	mu  sync.Mutex
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

// ConfigureRandomID 在首次生成 ID 前配置 workerID（0..1023）。
// 多进程/多主机必须由部署方分配不同 workerID；默认 MAC 哈希只适用于单进程。
func ConfigureRandomID(workerID int64) error {
	sfMu.Lock()
	defer sfMu.Unlock()
	if workerID < 0 || workerID > maxWorkerId {
		return errWorkerIDInvalid
	}
	if sfStarted {
		return errors.New("RandomId has already started; workerID cannot be changed")
	}
	sf = &snowflake{workerId: workerID, lastTimestamp: -1}
	return nil
}

// RandomId 根据雪花算法生成 ID；初始化失败或时钟回拨时 panic。
// 需要恢复错误的调用方可以使用 RandomIdWithError。
func RandomId() int64 {
	id, err := RandomIdWithError()
	if err != nil {
		panic(err)
	}
	return id
}

// RandomIdWithError 与 RandomId 使用同一生成器，但将失败返回给调用方。
func RandomIdWithError() (int64, error) {
	sfMu.Lock()
	if sf == nil {
		instance, err := newSnowflake()
		if err != nil {
			sfMu.Unlock()
			return 0, err
		}
		sf = instance
	}
	sfStarted = true
	instance := sf
	sfMu.Unlock()
	return instance.nextID()
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
	sfMu                sync.Mutex
	sfStarted           bool
)

type snowflake struct {
	sync.Mutex
	workerId      int64
	lastTimestamp int64
	sequence      int64
	// Set only when constructing an instance; nil uses the wall clock.
	nowMillis func() int64
}

func (s *snowflake) currentMillis() int64 {
	if s.nowMillis != nil {
		return s.nowMillis()
	}
	return time.Now().UnixNano() / 1e6
}

// NextID 生成下一个ID
func (s *snowflake) nextID() (int64, error) {
	s.Lock()
	defer s.Unlock()

	timestamp := s.currentMillis()

	if timestamp < s.lastTimestamp {
		return 0, errTimeBackwards
	}

	sequence := int64(0)
	if timestamp == s.lastTimestamp {
		sequence = (s.sequence + 1) & sequenceMask
		if sequence == 0 {
			var err error
			timestamp, err = s.tilNextMillis(s.lastTimestamp)
			if err != nil {
				return 0, err
			}
		}
	}

	// Do not reset an exhausted sequence if waiting fails: retrying in the same
	// millisecond must wait again, rather than reuse already-issued sequence values.
	s.lastTimestamp = timestamp
	s.sequence = sequence

	return (timestamp << timestampShift) |
		(s.workerId << workerShift) |
		s.sequence, nil
}

func (s *snowflake) tilNextMillis(lastTimestamp int64) (int64, error) {
	for {
		timestamp := s.currentMillis()
		if timestamp < lastTimestamp {
			return 0, errTimeBackwards
		}
		if timestamp > lastTimestamp {
			return timestamp, nil
		}
	}
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
	mu.Lock()
	defer mu.Unlock()
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
		if int64(int(val)) != val {
			panic("random integer argument is outside int range")
		}
		return true, int(val)
	default:
		return false, 0
	}
}

func nextInt2(min, max int) int {
	if min >= max {
		panic("random interval must have min < max")
	}
	mu.Lock()
	defer mu.Unlock()
	// Unsigned subtraction also represents intervals wider than MaxInt.
	width := uint(max) - uint(min)
	if width <= uint(^uint(0)>>1) {
		return rng.Intn(int(width)) + min
	}
	// Reject the incomplete modulo bucket to preserve uniform sampling.
	threshold := -width % width
	for {
		n := uint(rng.Uint64())
		if n >= threshold {
			return int(uint(min) + n%width)
		}
	}
}

func nextInt1(num int) int {
	mu.Lock()
	defer mu.Unlock()
	return rng.Intn(num)
}
