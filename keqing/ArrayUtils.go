package keqing

import (
	"fmt"
	"reflect"
	"sort"
)

/*
判断数组中是否包含某个元素
*/
func ArrayContains(obj interface{}, value interface{}) bool {
	if !isSlice(obj) {
		return false
	}
	// 获取 obj 的反射值
	objVal := reflect.ValueOf(obj)
	// 获取 value 的反射值
	valueVal := reflect.ValueOf(value)

	// 遍历切片
	for i := 0; i < objVal.Len(); i++ {
		elem := objVal.Index(i)
		// 检查元素类型是否与 value 相同
		if elem.Type() == valueVal.Type() {
			// 使用 reflect.DeepEqual 比较值
			if reflect.DeepEqual(elem.Interface(), value) {
				return true
			}
		}
	}
	return false
}

/*
获取数组第0个元素
*/
func GetOne[T comparable](arr []T) (T, bool) {
	if len(arr) > 0 {
		return arr[0], true
	}
	var zero T
	return zero, false
}

/*
两个数组取差集
*/
func Difference[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	var result []T

	// Add elements from the first slice to the map
	for _, item := range a {
		m[item] = true
	}

	// Remove elements found in the second slice
	for _, item := range b {
		delete(m, item)
	}

	// Build the result slice from the remaining keys in the map
	for k := range m {
		result = append(result, k)
	}

	return result
}

/*
两个数组取并集
*/
func Union[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	var result []T

	// Add elements from both slices to the map
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		m[item] = true
	}

	// Build the result slice from the map keys
	for k := range m {
		result = append(result, k)
	}

	return result
}

/*
两个数组取交集
*/
func Intersection[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	var result []T

	// Add elements from the first slice to the map
	for _, item := range a {
		m[item] = true
	}

	// Check elements from the second slice against the map
	for _, item := range b {
		if m[item] {
			result = append(result, item)
			m[item] = false // Mark as used
		}
	}

	return result
}

func SortArrayDesc(obj interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // 获取指针指向的值
		if val.Kind() != reflect.Slice {
			panic("object is not a slice")
		}
	} else if val.Kind() != reflect.Slice {
		panic("object is not a slice")
	}

	// 检查切片元素的类型
	elemType := val.Type().Elem()

	switch elemType.Kind() {
	case reflect.Int:
		// 如果是 int 切片，则进行倒序排序
		slice := make([]int, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = int(val.Index(i).Int())
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetInt(int64(slice[i]))
		}
		return
	case reflect.Int32:
		// 如果是 int32 切片，则进行倒序排序
		slice := make([]int, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = int(val.Index(i).Int())
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetInt(int64(slice[i]))
		}
		return
	case reflect.Int64:
		// 如果是 int64 切片，则进行倒序排序
		slice := make([]int, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = int(val.Index(i).Int())
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetInt(int64(slice[i]))
		}
		return
	case reflect.Float32:
		// 如果是 float32 切片，则进行倒序排序
		slice := make([]float64, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Float()
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetFloat(slice[i])
		}
		return

	case reflect.Float64:
		// 如果是 float64 切片，则进行倒序排序
		slice := make([]float64, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Float()
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetFloat(slice[i])
		}
		return

	case reflect.String:
		// 如果是 string 切片，则进行倒序排序
		slice := make([]string, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Interface().(string)
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetString(slice[i])
		}
		return

	default:
		// 如果是其他类型，则返回错误
		panic(fmt.Sprintf("unsupported slice element type: %s", elemType))
	}
}

func SortArray(obj interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // 获取指针指向的值
		if val.Kind() != reflect.Slice {
			panic("object is not a slice")
		}
	} else if val.Kind() != reflect.Slice {
		panic("object is not a slice")
	}

	// 检查切片元素的类型
	elemType := val.Type().Elem()

	switch elemType.Kind() {
	case reflect.Int:
		// 如果是 int 切片，则进行排序
		slice := make([]int, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = int(val.Index(i).Int())
		}
		sort.Ints(slice)
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetInt(int64(slice[i]))
		}
		return
	case reflect.Int32:
		// 如果是 int 切片，则进行排序
		slice := make([]int, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = int(val.Index(i).Int())
		}
		sort.Ints(slice)
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetInt(int64(slice[i]))
		}
		return
	case reflect.Int64:
		// 如果是 int 切片，则进行排序
		slice := make([]int, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = int(val.Index(i).Int())
		}
		sort.Ints(slice)
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetInt(int64(slice[i]))
		}
		return
	case reflect.Float32:
		// 如果是 int 切片，则进行排序
		slice := make([]float64, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Float()
		}
		sort.Float64s(slice)
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetFloat(slice[i])
		}
		return
	case reflect.Float64:
		// 如果是 int 切片，则进行排序
		slice := make([]float64, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Float()
		}
		sort.Float64s(slice)
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetFloat(slice[i])
		}
		return
	case reflect.String:
		// 如果是 string 切片，则进行排序
		slice := make([]string, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Interface().(string)
		}
		sort.Strings(slice)
		// 将排序后的切片替换到原始切片中
		for i := 0; i < val.Len(); i++ {
			val.Index(i).SetString(slice[i])
		}
		return
	default:
		// 如果是其他类型，则返回错误
		panic(fmt.Sprintf("unsupported slice element type: %s", elemType))
	}
}
