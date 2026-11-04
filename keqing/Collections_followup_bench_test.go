package keqing

import "testing"

func BenchmarkCollectionFollowup8000(b *testing.B) {
	const count = 8000
	values, otherValues := make([]int, count), make([]int, count)
	pointers := make([]*followupEntity, count)
	for i := range values {
		values[i], otherValues[i] = i, count+i
		pointers[i] = &followupEntity{ID: i}
	}
	left, right := AsList(values, nil), AsList(otherValues, nil)
	set := AsSet(pointers, nil)
	b.Run("ListUnion", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := left.Union(right).Size(); got != 2*count {
				b.Fatal(got)
			}
		}
	})
	b.Run("ListIntersectDisjoint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := left.Intersect(right).Size(); got != 0 {
				b.Fatal(got)
			}
		}
	})
	b.Run("SetRetainAllInts", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s := AsSet(values, nil)
			s.RetainAll(values...)
			if got := s.Size(); got != count {
				b.Fatal(got)
			}
		}
	})
	b.Run("ListRetainAllInts", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			list := AsList(values, nil)
			list.RetainAll(values...)
			if got := list.Size(); got != count {
				b.Fatal(got)
			}
		}
	})
	b.Run("SetCopyPointers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := set.Copy().Size(); got != count {
				b.Fatal(got)
			}
		}
	})
	b.Run("SetFilterPointers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := set.Filter(func(*followupEntity) bool { return true }).Size(); got != count {
				b.Fatal(got)
			}
		}
	})
	b.Run("SetUnionPointers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := set.Union(set).Size(); got != count {
				b.Fatal(got)
			}
		}
	})
}
