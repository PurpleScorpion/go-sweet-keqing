package keqing

import (
	"sync"
)

// Set 是一个线程安全、可自定义 hash 函数的集合结构
type Set[T comparable] struct {
	data           []T
	hashMap        map[uint64][]int
	hashFunc       func(T) uint64
	scanOnHashMiss bool
	mu             sync.RWMutex
}

// NewSet 创建一个新的 Set，可以传入自定义 hash 函数（可为 nil）
func NewSet[T comparable](hashFunc func(T) uint64) *Set[T] {
	scanOnHashMiss := collectionNeedsHashScan[T](hashFunc != nil)
	if hashFunc == nil {
		hashFunc = func(t T) uint64 {
			return defaultHash(t)
		}
	}
	return newSetWithHash(hashFunc, scanOnHashMiss)
}

// Preserve strategy metadata when copying, filtering, or converting a collection.
func newSetWithHash[T comparable](hashFunc func(T) uint64, scanOnHashMiss bool) *Set[T] {
	return &Set[T]{
		data:           make([]T, 0),
		hashMap:        make(map[uint64][]int),
		hashFunc:       hashFunc,
		scanOnHashMiss: scanOnHashMiss,
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
	return s.indexOfByLocked(value, s.equal)
}

// Keep the lookup algorithm separate from equality so its comparison budget can
// be verified without substituting an Equal method that changes the key contract.
func (s *Set[T]) indexOfByLocked(value T, equal func(T, T) bool) int {
	for _, idx := range s.hashMap[s.hashFunc(value)] {
		if equal(s.data[idx], value) {
			return idx
		}
	}
	if !s.scanOnHashMiss {
		return -1
	}
	// Equal may match across hashes, and a custom hash may inspect changed fields.
	for idx, item := range s.data {
		if equal(item, value) {
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

// Iterator 返回创建时的浅快照；可提前停止消费，也可在遍历中修改原集合。
func (s *Set[T]) Iterator() <-chan T {
	// A buffered snapshot needs no producer goroutine or lock after return.
	snapshot := s.GetData()
	ch := make(chan T, len(snapshot))
	for _, item := range snapshot {
		ch <- item
	}
	close(ch)
	return ch
}

// Union 求两个集合的并集
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	left, right := s.GetData(), other.GetData()
	union := newSetWithHash(s.hashFunc, s.scanOnHashMiss)
	union.addAll(left...)
	union.addAll(right...)
	return union
}

// Intersect 求两个集合的交集
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	left, right := s.GetData(), other.GetData()
	otherSnapshot := newSetWithHash(other.hashFunc, other.scanOnHashMiss)
	otherSnapshot.data = right
	otherSnapshot.rebuildIndexLocked()
	intersection := newSetWithHash(s.hashFunc, s.scanOnHashMiss)
	for _, item := range left {
		if otherSnapshot.Contains(item) {
			intersection.Add(item)
		}
	}
	return intersection
}

// Difference 求两个集合的差集（属于 s 不属于 other 的元素）
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	left, right := s.GetData(), other.GetData()
	otherSnapshot := newSetWithHash(other.hashFunc, other.scanOnHashMiss)
	otherSnapshot.data = right
	otherSnapshot.rebuildIndexLocked()
	diff := newSetWithHash(s.hashFunc, s.scanOnHashMiss)
	for _, item := range left {
		if !otherSnapshot.Contains(item) {
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
	s.data = filterCollection(s.data, items, s.hashFunc, s.scanOnHashMiss, false, s.equal)
	s.rebuildIndexLocked()
}

// RetainAll 保留指定的元素，删除其他所有
func (s *Set[T]) RetainAll(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = filterCollection(s.data, items, s.hashFunc, s.scanOnHashMiss, true, s.equal)
	s.rebuildIndexLocked()
}

// Filter 返回一个满足条件的新 Set，不会修改原集合
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T] {
	snapshot := s.GetData()

	newSet := newSetWithHash(s.hashFunc, s.scanOnHashMiss)

	for _, item := range snapshot {
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

	copySet := newSetWithHash(s.hashFunc, s.scanOnHashMiss)

	for _, item := range s.data {
		copySet.Add(item)
	}

	return copySet
}

// DeepCopy 返回一个新的 Set，使用 cloneFunc 对每个元素进行深拷贝
func (s *Set[T]) DeepCopy(cloneFunc func(T) T) *Set[T] {
	snapshot := s.GetData()

	copySet := newSetWithHash(s.hashFunc, s.scanOnHashMiss)

	for _, item := range snapshot {
		copySet.Add(cloneFunc(item))
	}

	return copySet
}

// ToList 将 Set 转换为 List
func (s *Set[T]) ToList() *List[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := newListWithHash(s.hashFunc, s.scanOnHashMiss)
	for _, item := range s.data {
		list.Add(item)
	}
	return list
}
