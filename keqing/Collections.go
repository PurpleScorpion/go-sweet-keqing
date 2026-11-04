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

// collectionEqual keeps equality independent of the hash function. Interface
// values can hold dynamically non-comparable values even when T is comparable.
func collectionEqual[T comparable](a, b T) bool {
	av, bv := reflect.ValueOf(a), reflect.ValueOf(b)
	if !av.IsValid() || !bv.IsValid() {
		return !av.IsValid() && !bv.IsValid()
	}
	// Do not invoke a method on a nil pointer receiver.
	if collectionNil(av) || collectionNil(bv) {
		return reflect.DeepEqual(a, b)
	}
	if equal, ok := any(a).(comparableEqual[T]); ok {
		return equal.Equal(b)
	}
	if av.Comparable() && bv.Comparable() {
		return a == b
	}
	return reflect.DeepEqual(a, b)
}

func collectionNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
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
	snapshot := list.GetData()

	result := make(map[K]T, len(snapshot))

	for _, item := range snapshot {
		key := keyFunc(item)
		_, exists := result[key]

		switch onConflict {
		case OverwriteNone:
			if exists {
				panic("Duplicate key detected: " + fmt.Sprintf("%v", key))
			}
			result[key] = item
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
	snapshot := set.GetData()

	result := make(map[K]T)

	for _, item := range snapshot {
		key := keyFunc(item)
		_, exists := result[key]

		switch onConflict {
		case OverwriteNone:
			if exists {
				panic("Duplicate key detected: " + fmt.Sprintf("%v", key))
			}
			result[key] = item
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
