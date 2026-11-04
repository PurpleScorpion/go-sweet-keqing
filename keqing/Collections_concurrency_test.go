package keqing

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type iteratorCollection interface {
	Iterator() <-chan int
	AddForTest(int)
	Remove(int) bool
	GetData() []int
}
type listIteratorAdapter struct{ *List[int] }

func (l listIteratorAdapter) AddForTest(v int) { l.Add(v) }

type setIteratorAdapter struct{ *Set[int] }

func (s setIteratorAdapter) AddForTest(v int) { s.Add(v) }

func TestCollectionIteratorCanBeAbandoned(t *testing.T) {
	factories := []struct {
		name string
		make func() iteratorCollection
	}{
		{"List", func() iteratorCollection { return listIteratorAdapter{AsList([]int{1, 2, 3}, nil)} }},
		{"Set", func() iteratorCollection { return setIteratorAdapter{AsSet([]int{1, 2, 3}, nil)} }},
	}
	for _, factory := range factories {
		t.Run(factory.name, func(t *testing.T) {
			for _, readFirst := range []bool{false, true} {
				c := factory.make()
				ch := c.Iterator()
				var got []int
				if readFirst {
					got = append(got, <-ch)
				}
				done := make(chan struct{})
				go func() { c.AddForTest(4); close(done) }()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Error("abandoned iterator blocks a subsequent write")
				}
				// Always drain for cleanup, including the failing original implementation.
				for v := range ch {
					got = append(got, v)
				}
				<-done
				if !reflect.DeepEqual(got, []int{1, 2, 3}) {
					t.Fatalf("iterator snapshot changed: %v", got)
				}
			}
		})
	}
}

func TestCollectionIteratorAllowsMutationInLoop(t *testing.T) {
	l := AsList([]int{1, 2, 3}, nil)
	s := AsSet([]int{1, 2, 3}, nil)
	for _, c := range []iteratorCollection{listIteratorAdapter{l}, setIteratorAdapter{s}} {
		done := make(chan struct{})
		go func() {
			for v := range c.Iterator() {
				c.Remove(v)
			}
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("mutation inside iterator loop deadlocked")
		}
		if len(c.GetData()) != 0 {
			t.Fatalf("loop did not remove original elements: %v", c.GetData())
		}
	}
}

func collectionPendingWriter(l *List[int], done chan<- struct{}) { l.Add(9); close(done) }
func waitForCollectionWriter(done <-chan struct{}) error {
	deadline := time.Now().Add(time.Second)
	for {
		select {
		case <-done:
			return nil
		default:
		}
		stack := make([]byte, 65536)
		n := runtime.Stack(stack, true)
		for _, g := range strings.Split(string(stack[:n]), "\n\n") {
			if strings.Contains(g, "collectionPendingWriter") && strings.Contains(g, "sync.(*RWMutex).Lock") {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("writer neither completed nor reached its lock")
		}
		runtime.Gosched()
	}
}

func TestCollectionOperationsDoNotRecursivelyReadLock(t *testing.T) {
	for _, operation := range []string{"Intersect", "Difference"} {
		var armed atomic.Bool
		entered, release := make(chan struct{}), make(chan struct{})
		other := NewList[int](func(v int) uint64 {
			if armed.CompareAndSwap(true, false) {
				close(entered)
				<-release
			}
			return uint64(v)
		})
		other.Add(1)
		other.Add(2)
		left := AsList([]int{1, 2}, nil)
		armed.Store(true)
		finished := make(chan struct{})
		go func() {
			if operation == "Intersect" {
				left.Intersect(other)
			} else {
				left.Difference(other)
			}
			close(finished)
		}()
		select {
		case <-entered:
		case <-time.After(time.Second):
			t.Fatal("operation did not reach hash callback")
		}
		written := make(chan struct{})
		go collectionPendingWriter(other, written)
		if err := waitForCollectionWriter(written); err != nil {
			close(release)
			t.Fatal(err)
		}
		close(release)
		select {
		case <-finished:
		case <-time.After(time.Second):
			t.Fatalf("%s deadlocked with a pending writer", operation)
		}
		select {
		case <-written:
		case <-time.After(time.Second):
			t.Fatalf("%s left writer blocked", operation)
		}
	}
}

func TestCollectionConcurrentOperationsAndSelf(t *testing.T) {
	l := AsList([]int{1, 2}, nil)
	s := AsSet([]int{1, 2}, nil)
	if !reflect.DeepEqual(l.Intersect(l).GetData(), []int{1, 2}) || !l.Difference(l).IsEmpty() || !reflect.DeepEqual(l.Union(l).GetData(), []int{1, 2}) {
		t.Fatal("List self operation mismatch")
	}
	if !reflect.DeepEqual(s.Intersect(s).GetData(), []int{1, 2}) || !s.Difference(s).IsEmpty() || !reflect.DeepEqual(s.Union(s).GetData(), []int{1, 2}) {
		t.Fatal("Set self operation mismatch")
	}
	var workers sync.WaitGroup
	workers.Add(3)
	go func() {
		defer workers.Done()
		for i := 0; i < 100; i++ {
			l.Add(i)
			l.Remove(i)
			s.Add(i)
			s.Remove(i)
		}
	}()
	go func() {
		defer workers.Done()
		for i := 0; i < 100; i++ {
			l.Intersect(l)
			l.Difference(l)
			l.Union(l)
		}
	}()
	go func() {
		defer workers.Done()
		for i := 0; i < 100; i++ {
			s.Intersect(s)
			s.Difference(s)
			s.Union(s)
		}
	}()
	done := make(chan struct{})
	go func() { workers.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("concurrent collection operations did not complete")
	}
}
