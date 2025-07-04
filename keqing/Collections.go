package keqing

import (
	"fmt"
	"hash/fnv"
	"reflect"
)

const (
	OverwriteNone = iota
	OverwriteFirst
	OverwriteLast
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

// AsSet 创建一个新的 Set，可以传入自定义 hash 函数（可为 nil）
func AsSet[T comparable](arr []T, hashFunc func(T) uint64) *Set[T] {
	set := NewSet[T](hashFunc)
	for _, item := range arr {
		set.Add(item)
	}
	return set
}

// AsList 创建 List 并使用指定 hash 函数
func AsList[T comparable](arr []T, hashFunc func(T) uint64) *List[T] {
	list := NewList[T](hashFunc)
	for _, item := range arr {
		list.Add(item)
	}
	return list
}

// List2Map List集合转为Map集合
// List2Map 支持 key 冲突策略
func List2Map[T comparable, K comparable](
	list *List[T],
	keyFunc func(T) K,
	onConflict int,
) map[K]T {
	list.mu.RLock()
	defer list.mu.RUnlock()

	result := make(map[K]T, len(list.data))

	for _, item := range list.data {
		key := keyFunc(item)
		_, exists := result[key]

		switch onConflict {
		case OverwriteNone:
			panic("Duplicate key detected: " + fmt.Sprintf("%v", key))
		case OverwriteFirst:
			if !exists {
				result[key] = item
			}
		case OverwriteLast:
			fallthrough
		default:
			result[key] = item
		}
	}

	return result
}

// Set2Map Set集合转为Map集合
// Set2Map 支持 key 冲突策略
func Set2Map[T comparable, K comparable](
	set *Set[T],
	keyFunc func(T) K,
	onConflict int,
) map[K]T {
	set.mu.RLock()
	defer set.mu.RUnlock()

	result := make(map[K]T)

	for _, item := range set.data {
		key := keyFunc(item)
		_, exists := result[key]

		switch onConflict {
		case OverwriteNone:
			panic("Duplicate key detected: " + fmt.Sprintf("%v", key))
		case OverwriteFirst:
			if !exists {
				result[key] = item
			}
		case OverwriteLast:
			fallthrough
		default:
			result[key] = item
		}
	}

	return result
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
