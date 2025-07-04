package keqing

import (
	"reflect"
	"sort"
	"sync"
)

// List 是一个线程安全、可自定义 hash 函数的有序列表结构
type List[T comparable] struct {
	data     []T
	hashMap  map[uint64][]int // 支持重复值，记录每个 hash 对应的所有索引位置
	hashFunc func(T) uint64
	mu       sync.RWMutex
}

// NewList 创建一个新的 List，可以传入自定义 hash 函数（可为 nil）
func NewList[T comparable](hashFunc func(T) uint64) *List[T] {
	if hashFunc == nil {
		hashFunc = func(t T) uint64 {
			return defaultHash(t)
		}
	}
	return &List[T]{
		data:     make([]T, 0),
		hashMap:  make(map[uint64][]int),
		hashFunc: hashFunc,
	}
}

// Add 添加元素到列表末尾
func (l *List[T]) Add(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	h := l.hashFunc(value)
	l.data = append(l.data, value)
	l.hashMap[h] = append(l.hashMap[h], len(l.data)-1)
}

// Set 替换指定索引位置的元素，并更新 hash 映射
func (l *List[T]) Set(index int, value T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index >= len(l.data) {
		panic("Index out of bounds: " + Num2Str(index))
	}

	oldValue := l.data[index]
	oldHash := l.hashFunc(oldValue)
	newHash := l.hashFunc(value)

	// 从 hashMap 中移除旧 hash 的当前索引
	if indices, exists := l.hashMap[oldHash]; exists {
		for i, idx := range indices {
			if idx == index {
				// 删除该索引记录
				l.hashMap[oldHash] = append(indices[:i], indices[i+1:]...)
				break
			}
		}
		// 如果该 hash 已无引用，则删除该键
		if len(l.hashMap[oldHash]) == 0 {
			delete(l.hashMap, oldHash)
		}
	}

	// 添加新的 hash 映射
	l.hashMap[newHash] = append(l.hashMap[newHash], index)

	// 替换数据
	l.data[index] = value
}

// Get 获取指定索引位置的元素
func (l *List[T]) Get(index int) T {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if index < 0 || index >= len(l.data) {
		panic("ArrayIndexOutOfBoundsException: " + Num2Str(index))
	}
	return l.data[index]
}

// Contains 判断是否包含某个元素
func (l *List[T]) Contains(value T) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	h := l.hashFunc(value)
	indices, exists := l.hashMap[h]
	if !exists {
		return false
	}
	for _, idx := range indices {
		if l.equal(l.data[idx], value) {
			return true
		}
	}
	return false
}

// IndexOf 返回第一个匹配项的索引位置，不存在则返回 -1
func (l *List[T]) IndexOf(value T) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	h := l.hashFunc(value)
	indices, exists := l.hashMap[h]
	if !exists {
		return -1
	}
	for _, idx := range indices {
		if l.equal(l.data[idx], value) {
			return idx
		}
	}
	return -1
}

// LastIndexOf 返回最后一个匹配项的索引位置，不存在则返回 -1
func (l *List[T]) LastIndexOf(value T) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	h := l.hashFunc(value)
	indices, exists := l.hashMap[h]
	if !exists {
		return -1
	}

	// 倒序遍历索引数组，找到第一个等于 value 的位置
	for i := len(indices) - 1; i >= 0; i-- {
		idx := indices[i]
		if l.equal(l.data[idx], value) {
			return idx
		}
	}

	return -1
}

// GetData 获取所有元素切片
func (l *List[T]) GetData() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return append([]T(nil), l.data...)
}

func (l *List[T]) ToArray() []T {
	return l.GetData()
}

// GetFirst 获取第一个元素，如果列表为空则 panic
func (l *List[T]) GetFirst() T {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.data) == 0 {
		panic("List is empty, cannot get first element")
	}
	return l.data[0]
}

// GetLast 获取最后一个元素，如果列表为空则 panic
func (l *List[T]) GetLast() T {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.data) == 0 {
		panic("List is empty, cannot get last element")
	}
	return l.data[len(l.data)-1]
}

// Size 返回列表大小
func (l *List[T]) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.data)
}

// IsEmpty 判断是否为空
func (l *List[T]) IsEmpty() bool {
	return l.Size() == 0
}

// Clear 清空列表
func (l *List[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.data = nil
	l.hashMap = make(map[uint64][]int)
}

// Remove 删除指定元素的第一个匹配项
func (l *List[T]) Remove(value T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	h := l.hashFunc(value)
	indices, exists := l.hashMap[h]
	if !exists {
		return false
	}

	for i, idx := range indices {
		if l.equal(l.data[idx], value) {
			// 删除 data 中的元素
			l.data = append(l.data[:idx], l.data[idx+1:]...)

			// 更新 hashMap 中所有后续索引
			deleteKey := false
			if len(indices) == 1 {
				deleteKey = true
			}
			if deleteKey {
				delete(l.hashMap, h)
			} else {
				// 移除当前索引并更新后续索引
				l.hashMap[h] = append(indices[:i], indices[i+1:]...)
			}

			// 调整后续元素的索引映射
			for key, list := range l.hashMap {
				for j := range list {
					if list[j] > idx {
						list[j]--
					}
				}
				l.hashMap[key] = list
			}

			return true
		}
	}
	return false
}

// RemoveAt 删除指定索引位置的元素
func (l *List[T]) RemoveAt(index int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index >= len(l.data) {
		return false
	}

	value := l.data[index]
	h := l.hashFunc(value)

	// 删除 hashMap 中该索引
	indices, exists := l.hashMap[h]
	if exists {
		newIndices := make([]int, 0)
		for _, idx := range indices {
			if idx != index {
				newIndices = append(newIndices, idx)
			}
		}
		if len(newIndices) == 0 {
			delete(l.hashMap, h)
		} else {
			l.hashMap[h] = newIndices
		}
	}

	// 删除 data 中的元素并调整后续索引
	l.data = append(l.data[:index], l.data[index+1:]...)

	// 更新所有后续索引
	for key, list := range l.hashMap {
		newList := make([]int, 0)
		for _, idx := range list {
			if idx > index {
				newList = append(newList, idx-1)
			} else if idx < index {
				newList = append(newList, idx)
			}
		}
		l.hashMap[key] = newList
	}

	return true
}

// Iterator 返回一个只读通道作为迭代器
func (l *List[T]) Iterator() <-chan T {
	l.mu.RLock()
	ch := make(chan T)
	go func() {
		for _, item := range l.data {
			ch <- item
		}
		close(ch)
		l.mu.RUnlock()
	}()
	return ch
}

// equal 比较两个元素是否相等（使用 hash 相等 + Go 的 == 运算符兜底）
func (l *List[T]) equal(a, b T) bool {
	if l.hashFunc(a) == l.hashFunc(b) {
		equal, ok := any(a).(comparableEqual[T])
		if ok && equal.Equal(b) {
			return true
		}
		return true
	}
	return false
}

// Union 返回两个 List 的并集（保留顺序）
func (l *List[T]) Union(other *List[T]) *List[T] {
	l.mu.RLock()
	other.mu.RLock()
	defer l.mu.RUnlock()
	defer other.mu.RUnlock()

	union := NewList[T](l.hashFunc)

	// 添加自身元素
	for _, item := range l.data {
		union.Add(item)
	}

	// 添加 other 中未包含的元素
	for _, item := range other.data {
		if !union.Contains(item) {
			union.Add(item)
		}
	}

	return union
}

// Intersect 返回两个 List 的交集（保留顺序）
func (l *List[T]) Intersect(other *List[T]) *List[T] {
	l.mu.RLock()
	other.mu.RLock()
	defer l.mu.RUnlock()
	defer other.mu.RUnlock()

	intersection := NewList[T](l.hashFunc)

	for _, item := range l.data {
		if other.Contains(item) {
			intersection.Add(item)
		}
	}

	return intersection
}

// Difference 返回两个 List 的差集（属于 l 不属于 other 的元素）
func (l *List[T]) Difference(other *List[T]) *List[T] {
	l.mu.RLock()
	other.mu.RLock()
	defer l.mu.RUnlock()
	defer other.mu.RUnlock()

	diff := NewList[T](l.hashFunc)

	for _, item := range l.data {
		if !other.Contains(item) {
			diff.Add(item)
		}
	}

	return diff
}

// RemoveAll 批量删除多个元素的所有匹配项
func (l *List[T]) RemoveAll(items ...T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, item := range items {
		h := l.hashFunc(item)
		indices, exists := l.hashMap[h]
		if !exists {
			continue
		}

		var newIndices []int
		for _, idx := range indices {
			if !l.equal(l.data[idx], item) {
				newIndices = append(newIndices, idx)
			} else {
				// 删除 data 中的元素
				l.data = append(l.data[:idx], l.data[idx+1:]...)

				// 更新后续索引
				for key, list := range l.hashMap {
					var newList []int
					for _, i := range list {
						if i > idx {
							newList = append(newList, i-1)
						} else if i < idx {
							newList = append(newList, i)
						}
					}
					l.hashMap[key] = newList
				}
			}
		}

		if len(newIndices) == 0 {
			delete(l.hashMap, h)
		} else {
			l.hashMap[h] = newIndices
		}
	}
}

// RetainAll 保留指定集合中的元素，其余全部删除
func (l *List[T]) RetainAll(items ...T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	itemSet := make(map[uint64][]T)
	for _, item := range items {
		h := l.hashFunc(item)
		itemSet[h] = append(itemSet[h], item)
	}

	newHashMap := make(map[uint64][]int)
	newData := make([]T, 0)

	var newIndex int
	for _, item := range l.data {
		h := l.hashFunc(item)
		candidates, ok := itemSet[h]
		if !ok {
			continue
		}

		found := false
		for _, candidate := range candidates {
			if l.equal(item, candidate) {
				found = true
				break
			}
		}

		if found {
			newData = append(newData, item)
			newHashMap[h] = append(newHashMap[h], newIndex)
			newIndex++
		}
	}

	l.data = newData
	l.hashMap = newHashMap
}

// Filter 返回一个满足条件的新 List，不会修改原集合
func (l *List[T]) Filter(predicate func(T) bool) *List[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()

	filtered := NewList[T](l.hashFunc)

	for _, item := range l.data {
		if predicate(item) {
			filtered.Add(item)
		}
	}

	return filtered
}

// SubList 返回一个新的 List，表示原列表的一个子集。
// 支持两种调用方式：
// - SubList(start int)：从 start 开始到末尾
// - SubList(start int, end int)：从 start（包含）到 end（不包含）
func (l *List[T]) SubList(indexs ...int) *List[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(indexs) == 0 || len(indexs) > 2 {
		panic("SubList requires 1 or 2 arguments")
	}

	start := indexs[0]
	end := len(l.data)

	if len(indexs) == 2 {
		end = indexs[1]
	}

	// 检查索引越界
	if start < 0 || start >= len(l.data) {
		panic("Start index out of bounds: " + Num2Str(start))
	}
	if end < start || end > len(l.data) {
		panic("End index out of bounds: " + Num2Str(end))
	}

	sub := NewList[T](l.hashFunc)
	for i := start; i < end; i++ {
		sub.Add(l.data[i])
	}

	return sub
}

// Sort 对列表进行排序，接受一个比较函数 less(a, b T) bool，表示 a 是否应排在 b 前面
func (l *List[T]) Sort(less func(a, b T) bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 使用 sort.SliceStable 进行稳定排序
	sort.SliceStable(l.data, func(i, j int) bool {
		return less(l.data[i], l.data[j])
	})

	// 重建 hashMap，因为索引已变化
	l.hashMap = make(map[uint64][]int)
	for idx, item := range l.data {
		h := l.hashFunc(item)
		l.hashMap[h] = append(l.hashMap[h], idx)
	}
}

// SortDESC 对基本类型列表进行降序排序（仅支持 string/int/float）
func (l *List[T]) SortDESC() {
	l.sortBasic(false)
}

// SortASC 对基本类型列表进行升序排序（仅支持 string/int/float）
func (l *List[T]) SortASC() {
	l.sortBasic(true)
}

func (l *List[T]) sortBasic(asc bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 获取当前列表元素的类型
	var example T
	elemType := reflect.TypeOf(example)

	// 检查是否为基础类型
	switch elemType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		// 合法类型，继续
	default:
		panic("SortASC/SortDESC 只允许用于基本类型（int/uint/float/string）")
	}

	// 排序逻辑
	sortFunc := func(i, j int) bool {
		a := l.data[i]
		b := l.data[j]

		switch elemType.Kind() {
		case reflect.Int:
			return compareInt(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Int8:
			return compareInt8(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Int16:
			return compareInt16(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Int32:
			return compareInt32(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Int64:
			return compareInt64(reflect.ValueOf(a), reflect.ValueOf(b), asc)

		case reflect.Uint:
			return compareUint(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Uint8:
			return compareUint8(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Uint16:
			return compareUint16(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Uint32:
			return compareUint32(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Uint64:
			return compareUint64(reflect.ValueOf(a), reflect.ValueOf(b), asc)

		case reflect.Float32:
			return compareFloat32(reflect.ValueOf(a), reflect.ValueOf(b), asc)
		case reflect.Float64:
			return compareFloat64(reflect.ValueOf(a), reflect.ValueOf(b), asc)

		case reflect.String:
			return compareString(reflect.ValueOf(a).String(), reflect.ValueOf(b).String(), asc)

		default:
			panic("不支持的排序类型")
		}
	}

	sort.SliceStable(l.data, sortFunc)

	// 重建 hashMap
	l.hashMap = make(map[uint64][]int)
	for idx, item := range l.data {
		h := l.hashFunc(item)
		l.hashMap[h] = append(l.hashMap[h], idx)
	}
}

func (l *List[T]) DeepCopy(cloneFunc func(T) T) *List[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()

	copyList := NewList[T](l.hashFunc)

	for _, item := range l.data {
		copyList.Add(cloneFunc(item))
	}

	return copyList
}

// Copy 返回一个新的 List，包含当前列表的所有元素（浅拷贝）
func (l *List[T]) Copy() *List[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()

	copyList := NewList[T](l.hashFunc)

	for _, item := range l.data {
		copyList.Add(item)
	}

	return copyList
}

func (l *List[T]) ToSet() *Set[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()

	set := NewSet[T](l.hashFunc)

	for _, item := range l.data {
		set.Add(item)
	}

	return set
}

func (l *List[T]) ToString() string {
	return ToJsonString(l.GetData())
}
