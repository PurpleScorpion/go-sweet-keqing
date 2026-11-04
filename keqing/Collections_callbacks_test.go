package keqing

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// The callback performs a read and a write against the original collection.
// No scheduler timing is needed to expose a callback retained under its read lock.
func TestCollectionCallbacksUseSnapshots(t *testing.T) {
	tests := []struct {
		name string
		run  func() error
	}{
		{"List.Filter", func() error {
			source := AsList([]int{1, 2}, nil)
			result := source.Filter(func(v int) bool {
				if !source.Contains(v) {
					panic("snapshot value missing")
				}
				source.Add(v + 10)
				return v == 2
			})
			if !reflect.DeepEqual(result.GetData(), []int{2}) || !reflect.DeepEqual(source.GetData(), []int{1, 2, 11, 12}) {
				return fmt.Errorf("result=%v source=%v", result.GetData(), source.GetData())
			}
			return nil
		}},
		{"Set.Filter", func() error {
			source := AsSet([]int{1, 2}, nil)
			result := source.Filter(func(v int) bool {
				if !source.Contains(v) {
					panic("snapshot value missing")
				}
				source.Add(v + 10)
				return v == 2
			})
			if !reflect.DeepEqual(result.GetData(), []int{2}) || !reflect.DeepEqual(source.GetData(), []int{1, 2, 11, 12}) {
				return fmt.Errorf("result=%v source=%v", result.GetData(), source.GetData())
			}
			return nil
		}},
		{"List.DeepCopy", func() error {
			source := AsList([]int{1, 2}, nil)
			result := source.DeepCopy(func(v int) int {
				if !source.Contains(v) {
					panic("snapshot value missing")
				}
				source.Add(v + 10)
				return v * 2
			})
			if !reflect.DeepEqual(result.GetData(), []int{2, 4}) || !reflect.DeepEqual(source.GetData(), []int{1, 2, 11, 12}) {
				return fmt.Errorf("result=%v source=%v", result.GetData(), source.GetData())
			}
			return nil
		}},
		{"Set.DeepCopy", func() error {
			source := AsSet([]int{1, 2}, nil)
			result := source.DeepCopy(func(v int) int {
				if !source.Contains(v) {
					panic("snapshot value missing")
				}
				source.Add(v + 10)
				return v * 2
			})
			if !reflect.DeepEqual(result.GetData(), []int{2, 4}) || !reflect.DeepEqual(source.GetData(), []int{1, 2, 11, 12}) {
				return fmt.Errorf("result=%v source=%v", result.GetData(), source.GetData())
			}
			return nil
		}},
		{"List2Map", func() error {
			source := AsList([]int{1, 2}, nil)
			result := List2Map(source, func(v int) int {
				if source.IndexOf(v) < 0 {
					panic("snapshot value missing")
				}
				source.Add(v + 10)
				return v
			}, OverwriteNone)
			if !reflect.DeepEqual(result, map[int]int{1: 1, 2: 2}) || !reflect.DeepEqual(source.GetData(), []int{1, 2, 11, 12}) {
				return fmt.Errorf("result=%v source=%v", result, source.GetData())
			}
			return nil
		}},
		{"Set2Map", func() error {
			source := AsSet([]int{1, 2}, nil)
			result := Set2Map(source, func(v int) int {
				if !source.Contains(v) {
					panic("snapshot value missing")
				}
				source.Add(v + 10)
				return v
			}, OverwriteNone)
			if !reflect.DeepEqual(result, map[int]int{1: 1, 2: 2}) || !reflect.DeepEqual(source.GetData(), []int{1, 2, 11, 12}) {
				return fmt.Errorf("result=%v source=%v", result, source.GetData())
			}
			return nil
		}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			completed := make(chan error, 1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						completed <- fmt.Errorf("callback panic: %v", r)
					}
				}()
				completed <- tc.run()
			}()
			select {
			case err := <-completed:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(time.Second):
				t.Fatal("callback deadlocked when querying/mutating the original collection")
			}
		})
	}
}
