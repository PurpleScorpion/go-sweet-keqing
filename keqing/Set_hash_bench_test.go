package keqing

import "testing"

func BenchmarkSetBuild8000(b *testing.B) {
	for _, custom := range []bool{false, true} {
		name := "default"
		var h func(int) uint64
		if custom {
			name = "custom"
			h = func(v int) uint64 { return uint64(v) }
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				s := NewSet[int](h)
				for i := 0; i < 8000; i++ {
					s.Add(i)
				}
				if s.Size() != 8000 {
					b.Fatal(s.Size())
				}
			}
		})
	}
}
