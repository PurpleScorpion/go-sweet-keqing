package keqing

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
	"time"
)

// 使用前请看
// 本UUID功能是在google的uuid实现基础上修改的, 去除了大部分用不到的函数 , 仅保留一个UUID生成入口
// 并对使用方式进行简约化
// 仅供学习参考，更多详细使用请查看google的uuid实现: https://github.com/google/uuid

type uuID [16]byte

const randPoolSize = 16 * 16
const nanoPerMilli = 1000000

var (
	rander      = rand.Reader // random function
	poolEnabled = false
	poolMu      sync.Mutex
	poolPos     = randPoolSize     // protected with poolMu
	pool        [randPoolSize]byte // protected with poolMu
	Nil         uuID               // empty uuID, all zeros
	timeMu      sync.Mutex
	timeNow     = time.Now // for testing
)
var lastUUIDtime int64

func UUID() string {
	uuid := must(newUUID())
	return uuid.string()
}

func newUUID() (uuID, error) {
	uuid, err := newRandom()
	if err != nil {
		return uuid, err
	}
	makeUUID(uuid[:])
	return uuid, nil
}

func newRandom() (uuID, error) {
	if !poolEnabled {
		return newRandomFromReader(rander)
	}
	return newRandomFromPool()
}

func newRandomFromPool() (uuID, error) {
	var uuid uuID
	poolMu.Lock()
	if poolPos == randPoolSize {
		_, err := io.ReadFull(rander, pool[:])
		if err != nil {
			poolMu.Unlock()
			return Nil, err
		}
		poolPos = 0
	}
	copy(uuid[:], pool[poolPos:(poolPos+16)])
	poolPos += 16
	poolMu.Unlock()

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}

func makeUUID(uuid []byte) {
	_ = uuid[15] // bounds check
	t, s := getUUIDTime()

	uuid[0] = byte(t >> 40)
	uuid[1] = byte(t >> 32)
	uuid[2] = byte(t >> 24)
	uuid[3] = byte(t >> 16)
	uuid[4] = byte(t >> 8)
	uuid[5] = byte(t)
	uuid[6] = 0x70 | (0x0F & byte(s>>8))
	uuid[7] = byte(s)
}

func getUUIDTime() (milli, seq int64) {
	timeMu.Lock()
	defer timeMu.Unlock()

	nano := timeNow().UnixNano()
	milli = nano / nanoPerMilli
	// Sequence number is between 0 and 3906 (nanoPerMilli>>8)
	seq = (nano - milli*nanoPerMilli) >> 8
	now := milli<<12 + seq
	if now <= lastUUIDtime {
		now = lastUUIDtime + 1
		milli = now >> 12
		seq = now & 0xfff
	}
	lastUUIDtime = now
	return milli, seq
}

func newRandomFromReader(r io.Reader) (uuID, error) {
	var uuid uuID
	_, err := io.ReadFull(r, uuid[:])
	if err != nil {
		return Nil, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}

func (uuid uuID) string() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid uuID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

func must(uuid uuID, err error) uuID {
	if err != nil {
		panic(err)
	}
	return uuid
}

//func EnableRandPool() {
//	poolEnabled = true
//}
//
//func DisableRandPool() {
//	poolEnabled = false
//	defer poolMu.Unlock()
//	poolMu.Lock()
//	poolPos = randPoolSize
//}
