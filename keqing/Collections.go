package keqing

import (
	"fmt"
	"hash/fnv"
	"reflect"
)

type Cloneable[T any] interface {
	Clone() T
}

// comparableEqual 接口用于支持自定义相等判断
type comparableEqual[T any] interface {
	Equal(T) bool
}

// Hashable 接口用于自定义对象实现 hash 计算
type Hashable interface {
	Hash() uint64
}

// 默认使用 FNV-1a 算法计算 hash 值
func defaultHash(obj any) uint64 {
	s := fmt.Sprintf("%v", obj)
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

type BasicType interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string
}

func compareInt(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Int() < b.Int()
	}
	return a.Int() > b.Int()
}

func compareInt8(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Int() < b.Int()
	}
	return a.Int() > b.Int()
}

func compareInt16(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Int() < b.Int()
	}
	return a.Int() > b.Int()
}

func compareInt32(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Int() < b.Int()
	}
	return a.Int() > b.Int()
}

func compareInt64(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Int() < b.Int()
	}
	return a.Int() > b.Int()
}

func compareUint(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Uint() < b.Uint()
	}
	return a.Uint() > b.Uint()
}

func compareUint8(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Uint() < b.Uint()
	}
	return a.Uint() > b.Uint()
}

func compareUint16(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Uint() < b.Uint()
	}
	return a.Uint() > b.Uint()
}

func compareUint32(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Uint() < b.Uint()
	}
	return a.Uint() > b.Uint()
}

func compareUint64(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Uint() < b.Uint()
	}
	return a.Uint() > b.Uint()
}

func compareFloat32(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Float() < b.Float()
	}
	return a.Float() > b.Float()
}

func compareFloat64(a, b reflect.Value, asc bool) bool {
	if asc {
		return a.Float() < b.Float()
	}
	return a.Float() > b.Float()
}

func compareString(a, b string, asc bool) bool {
	if asc {
		return a < b
	}
	return a > b
}
