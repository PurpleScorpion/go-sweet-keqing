package keqing

import (
	"fmt"
	"math"
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

// Default hashing follows Go value equality instead of formatted output.
// Pointer values are hashed by identity, and all signed zeros share a hash.
func defaultHash(obj any) uint64 {
	h := collectionHash(14695981039346656037)
	h.addValue(reflect.ValueOf(obj))
	return uint64(h)
}

type collectionHash uint64

func (h *collectionHash) addByte(v byte) {
	*h = (*h ^ collectionHash(v)) * 1099511628211
}

func (h *collectionHash) addUint(v uint64) {
	for i := 0; i < 8; i++ {
		h.addByte(byte(v))
		v >>= 8
	}
}

func (h *collectionHash) addFloat(v float64) {
	if v == 0 {
		v = 0 // Normalize negative zero, including nested fields and complex parts.
	}
	h.addUint(math.Float64bits(v))
}

func (h *collectionHash) addValue(v reflect.Value) {
	if !v.IsValid() {
		h.addByte(byte(reflect.Invalid))
		return
	}
	h.addByte(byte(v.Kind()))
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			h.addByte(1)
		} else {
			h.addByte(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		h.addUint(uint64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		h.addUint(v.Uint())
	case reflect.Float32, reflect.Float64:
		h.addFloat(v.Float())
	case reflect.Complex64, reflect.Complex128:
		h.addFloat(real(v.Complex()))
		h.addFloat(imag(v.Complex()))
	case reflect.String:
		text := v.String()
		h.addUint(uint64(len(text)))
		for i := 0; i < len(text); i++ {
			h.addByte(text[i])
		}
	case reflect.Ptr, reflect.Chan, reflect.UnsafePointer:
		h.addUint(uint64(v.Pointer()))
	case reflect.Array:
		h.addUint(uint64(v.Len()))
		for i := 0; i < v.Len(); i++ {
			h.addValue(v.Index(i))
		}
	case reflect.Struct:
		h.addUint(uint64(v.NumField()))
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).Name != "_" {
				h.addValue(v.Field(i))
			}
		}
	case reflect.Interface:
		h.addValue(v.Elem())
		// Dynamically non-comparable values use a shared bucket and DeepEqual.
		// Do not traverse maps/slices: they can contain cycles and mutable contents.
	}
}

// Arbitrary Equal methods and dynamic interface values need a value scan when
// their hash bucket misses. A custom hash may also inspect mutable references.
func collectionNeedsHashScan[T comparable](customHash bool) bool {
	valueType := reflect.TypeOf((*T)(nil)).Elem()
	equalType := reflect.TypeOf((*comparableEqual[T])(nil)).Elem()
	return valueType.Implements(equalType) || collectionTypeNeedsScan(valueType, customHash)
}

func collectionTypeNeedsScan(t reflect.Type, customHash bool) bool {
	switch t.Kind() {
	case reflect.Interface:
		return true
	case reflect.Ptr, reflect.Chan, reflect.UnsafePointer:
		return customHash
	case reflect.Array:
		return collectionTypeNeedsScan(t.Elem(), customHash)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Name != "_" && collectionTypeNeedsScan(t.Field(i).Type, customHash) {
				return true
			}
		}
	}
	return false
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
