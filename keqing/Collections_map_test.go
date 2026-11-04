package keqing

import (
	"reflect"
	"testing"
)

func TestCollectionMapConflictPolicies(t *testing.T) {
	converters := map[string]func([]int, int) map[int]int{
		"List": func(values []int, policy int) map[int]int {
			return List2Map(AsList(values, nil), func(v int) int { return v % 2 }, policy)
		},
		"Set": func(values []int, policy int) map[int]int {
			return Set2Map(AsSet(values, nil), func(v int) int { return v % 2 }, policy)
		},
	}
	for name, convert := range converters {
		t.Run(name, func(t *testing.T) {
			for _, values := range [][]int{nil, {1}, {1, 2}} {
				got := convert(values, OverwriteNone)
				want := make(map[int]int)
				for _, v := range values {
					want[v%2] = v
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("got %v want %v", got, want)
				}
			}
			requirePanic(t, func() { convert([]int{1, 3}, OverwriteNone) })
			if got := convert([]int{1, 3}, OverwriteFirst); got[1] != 1 {
				t.Fatalf("first: %v", got)
			}
			if got := convert([]int{1, 3}, OverwriteLast); got[1] != 3 {
				t.Fatalf("last: %v", got)
			}
		})
	}
}
