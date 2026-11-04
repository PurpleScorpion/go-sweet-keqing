package keqing

import (
	"sync"
)

// Set 是一个线程安全、可自定义 hash 函数的集合结构
type Set[T comparable] struct {
	data     []T
	hashMap  map[uint64][]int
	hashFunc func(T) uint64
	mu       sync.RWMutex
}

// NewSet 创建一个新的 Set，可以传入自定义 hash 函数（可为 nil）
func NewSet[T comparable](hashFunc func(T) uint64) *Set[T] {
	if hashFunc == nil {
		hashFunc = func(t T) uint64 {
			return defaultHash(t)
		}
	}
	return &Set[T]{
		data:     make([]T, 0),
		hashMap:  make(map[uint64][]int),
		hashFunc: hashFunc,
	}
}

// Add 添加元素到 Set，自动去重
func (s *Set[T]) Add(value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.indexOfLocked(value) >= 0 {
		return false
	}
	h := s.hashFunc(value)
	s.data = append(s.data, value)
	s.hashMap[h] = append(s.hashMap[h], len(s.data)-1)
	return true
}

// Contains 判断是否包含某个元素
func (s *Set[T]) Contains(value T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.indexOfLocked(value) >= 0
}

func (s *Set[T]) indexOfLocked(value T) int {
	for _, idx := range s.hashMap[s.hashFunc(value)] {
		if s.equal(s.data[idx], value) {
			return idx
		}
	}
	// Equal may match across hashes, and a stored pointer's fields may change.
	for idx, item := range s.data {
		if s.equal(item, value) {
			return idx
		}
	}
	return -1
}

func (s *Set[T]) rebuildIndexLocked() {
	s.hashMap = make(map[uint64][]int)
	for idx, item := range s.data {
		h := s.hashFunc(item)
		s.hashMap[h] = append(s.hashMap[h], idx)
	}
}

// GetAll 获取所有元素切片
// Deprecated: Use NewMethod instead.
func (s *Set[T]) GetAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]T(nil), s.data...)
}

func (s *Set[T]) GetData() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]T(nil), s.data...)
}

func (s *Set[T]) ToArray() []T {
	return s.GetData()
}

// Size 返回集合大小
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear 清空集合
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = nil
	s.hashMap = make(map[uint64][]int)
}

// Iterator 返回一个只读通道作为迭代器
func (s *Set[T]) Iterator() <-chan T {
	s.mu.RLock()
	ch := make(chan T)
	go func() {
		for _, item := range s.data {
			ch <- item
		}
		close(ch)
		s.mu.RUnlock()
	}()
	return ch
}

// Union 求两个集合的并集
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	union := NewSet[T](s.hashFunc)
	union.addAll(s.data...)
	union.addAll(other.data...)
	return union
}

// Intersect 求两个集合的交集
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	intersection := NewSet[T](s.hashFunc)
	for _, item := range s.data {
		if other.indexOfLocked(item) >= 0 {
			intersection.Add(item)
		}
	}
	return intersection
}

// Difference 求两个集合的差集（属于 s 不属于 other 的元素）
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	diff := NewSet[T](s.hashFunc)
	for _, item := range s.data {
		if other.indexOfLocked(item) < 0 {
			diff.Add(item)
		}
	}
	return diff
}

// addAll 批量添加元素
func (s *Set[T]) addAll(items ...T) {
	for _, item := range items {
		s.Add(item)
	}
}

func (s *Set[T]) ToString() string {
	return ToJsonString(s.GetData())
}

// Remove 删除指定元素，返回是否删除成功
func (s *Set[T]) Remove(value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := s.indexOfLocked(value)
	if idx < 0 {
		return false
	}
	s.data = append(s.data[:idx], s.data[idx+1:]...)
	s.rebuildIndexLocked()
	return true
}

// equal 比较实际元素，hash 仅用于定位候选。
func (s *Set[T]) equal(a, b T) bool {
	return collectionEqual(a, b)
}

// RemoveAll 批量删除多个元素
func (s *Set[T]) RemoveAll(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	kept := make([]T, 0, len(s.data))
	for _, item := range s.data {
		remove := false
		for _, candidate := range items {
			if s.equal(item, candidate) {
				remove = true
				break
			}
		}
		if !remove {
			kept = append(kept, item)
		}
	}
	s.data = kept
	s.rebuildIndexLocked()
}

// RetainAll 保留指定的元素，删除其他所有
func (s *Set[T]) RetainAll(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	kept := make([]T, 0, len(s.data))
	for _, item := range s.data {
		for _, candidate := range items {
			if s.equal(item, candidate) {
				kept = append(kept, item)
				break
			}
		}
	}
	s.data = kept
	s.rebuildIndexLocked()
}

// Filter 返回一个满足条件的新 Set，不会修改原集合
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newSet := NewSet[T](s.hashFunc)

	for _, item := range s.data {
		if predicate(item) {
			newSet.Add(item)
		}
	}

	return newSet
}

// Copy 返回一个新的 Set，包含当前 Set 的所有元素（浅拷贝）
func (s *Set[T]) Copy() *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copySet := NewSet[T](s.hashFunc)

	for _, item := range s.data {
		copySet.Add(item)
	}

	return copySet
}

// DeepCopy 返回一个新的 Set，使用 cloneFunc 对每个元素进行深拷贝
func (s *Set[T]) DeepCopy(cloneFunc func(T) T) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copySet := NewSet[T](s.hashFunc)

	for _, item := range s.data {
		copySet.Add(cloneFunc(item))
	}

	return copySet
}

// ToList 将 Set 转换为 List
func (s *Set[T]) ToList() *List[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := NewList[T](s.hashFunc)
	for _, item := range s.data {
		list.Add(item)
	}
	return list
}
