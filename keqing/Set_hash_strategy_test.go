package keqing

import (
	"math"
	"testing"
)

type hashStringer int

func (v hashStringer) String() string { panic("default hashing must not invoke Stringer") }

func TestDefaultHashFollowsValueEquality(t *testing.T) {
	z := math.Copysign(0, -1)
	tests := [][2]any{
		{float64(0), z},
		{float32(0), float32(z)},
		{complex(0.0, 0.0), complex(z, z)},
		{[2]float64{0, 1}, [2]float64{z, 1}},
		{struct {
			F float64
			A [1]float64
		}{0, [1]float64{0}}, struct {
			F float64
			A [1]float64
		}{z, [1]float64{z}}},
	}
	for _, pair := range tests {
		if defaultHash(pair[0]) != defaultHash(pair[1]) {
			t.Fatalf("equal values hash differently: %#v, %#v", pair[0], pair[1])
		}
	}
	s := NewSet[hashStringer](nil)
	if !s.Add(1) || s.Add(1) || !s.Contains(1) {
		t.Fatal("named primitive set failed")
	}
}

func TestDefaultHashPointerIdentityIsStable(t *testing.T) {
	value := &struct{ N int }{1}
	before := defaultHash(value)
	value.N = 2
	if defaultHash(value) != before {
		t.Fatal("pointer hash depends on mutable content")
	}
	s := NewSet[*struct{ N int }](nil)
	s.Add(value)
	if !s.Add(&struct{ N int }{2}) {
		t.Fatal("different pointers were merged")
	}
	value.N = 3
	if !s.Contains(value) || !s.Remove(value) {
		t.Fatal("pointer identity was lost")
	}
}

func TestSetHashLookupComparisonBudget(t *testing.T) {
	const count = 1024
	for _, custom := range []bool{false, true} {
		name := "default"
		var h func(int) uint64
		if custom {
			name = "custom"
			h = func(v int) uint64 { return uint64(v) }
		}
		t.Run(name, func(t *testing.T) {
			s := NewSet[int](h)
			comparisons := 0
			equal := func(a, b int) bool { comparisons++; return a == b }
			for i := 0; i < count; i++ {
				if s.indexOfByLocked(i, equal) != -1 {
					t.Fatal("new value already found")
				}
				s.Add(i)
			}
			for i := 0; i < count; i++ {
				if s.indexOfByLocked(i, equal) < 0 {
					t.Fatal("inserted value missing")
				}
				if s.indexOfByLocked(i+count, equal) >= 0 {
					t.Fatal("absent value found")
				}
			}
			// A constant work budget catches full scans without wall-clock thresholds.
			// The previous implementation performs more than one million comparisons.
			if comparisons > 4*count {
				t.Fatalf("%d comparisons for %d values; hash lookup regressed to scanning", comparisons, count)
			}
		})
	}
}

func TestSetHashCollisionBucketsRemainCorrect(t *testing.T) {
	s := NewSet[int](func(int) uint64 { return 1 })
	for i := 0; i < 100; i++ {
		if !s.Add(i) {
			t.Fatalf("lost colliding value %d", i)
		}
	}
	for i := 0; i < 100; i++ {
		if !s.Contains(i) || s.Add(i) {
			t.Fatalf("collision lookup failed for %d", i)
		}
	}
	if s.Contains(100) || s.Remove(100) {
		t.Fatal("nonexistent colliding value matched")
	}
	if !s.Remove(50) || s.Contains(50) || s.Size() != 99 {
		t.Fatal("collision removal failed")
	}
}

func TestSetNestedSignedZeroUsesOneEntry(t *testing.T) {
	type record struct {
		N [2]float64
		C complex128
	}
	negative := math.Copysign(0, -1)
	a := record{N: [2]float64{0, 1}, C: complex(0, 0)}
	b := record{N: [2]float64{negative, 1}, C: complex(negative, negative)}
	s := NewSet[record](nil)
	if !s.Add(a) || s.Add(b) || !s.Contains(b) || !s.Remove(b) || !s.IsEmpty() {
		t.Fatal("equal nested zeros were split into different entries")
	}
}
