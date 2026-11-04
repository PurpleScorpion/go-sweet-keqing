package keqing

import (
	"reflect"
	"testing"
)

func assertListIndexState(t *testing.T, l *List[int], want []int) {
	t.Helper()
	got := l.GetData()
	if !reflect.DeepEqual(got, want) && !(len(got) == 0 && len(want) == 0) {
		t.Fatalf("data %v, want %v", got, want)
	}
	for _, v := range []int{0, 1, 2, 3, 4, 5} {
		first, last := -1, -1
		for i, x := range want {
			if x == v {
				if first < 0 {
					first = i
				}
				last = i
			}
		}
		if l.Contains(v) != (first >= 0) || l.IndexOf(v) != first || l.LastIndexOf(v) != last {
			t.Fatalf("lookup for %v inconsistent with %v", v, want)
		}
	}
}

func TestListIndexMaintenanceSet(t *testing.T) {
	l := AsList([]int{1, 2, 1}, nil)
	l.Set(0, 1)
	assertListIndexState(t, l, []int{1, 2, 1})
	l.Set(1, 1)
	assertListIndexState(t, l, []int{1, 1, 1})
	l.Remove(1)
	assertListIndexState(t, l, []int{1, 1})
	l.Set(1, 2)
	assertListIndexState(t, l, []int{1, 2})
	l.RemoveAt(0)
	assertListIndexState(t, l, []int{2})
}

func TestListIndexMaintenanceRemoveAll(t *testing.T) {
	tests := []struct {
		name               string
		data, remove, want []int
	}{
		{"adjacent all", []int{1, 1}, []int{1}, nil},
		{"adjacent partial", []int{1, 1, 2}, []int{1}, []int{2}},
		{"separated", []int{1, 2, 1, 3, 1}, []int{1}, []int{2, 3}},
		{"multiple targets", []int{1, 2, 1, 3, 2, 4}, []int{1, 2}, []int{3, 4}},
		{"duplicate targets", []int{1, 2, 1, 3}, []int{1, 1}, []int{2, 3}},
		{"absent", []int{1, 2}, []int{5}, []int{1, 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := AsList(tc.data, nil)
			l.RemoveAll(tc.remove...)
			assertListIndexState(t, l, tc.want)
		})
	}
}

func TestListIndexMaintenanceBulkRemovalAfterSetAndCollision(t *testing.T) {
	l := AsList([]int{1, 2, 1, 3, 2}, func(int) uint64 { return 0 })
	l.Set(0, 2)
	l.RemoveAll(2)
	assertListIndexState(t, l, []int{1, 3})
	l.RemoveAll(1, 3)
	assertListIndexState(t, l, nil)
}
