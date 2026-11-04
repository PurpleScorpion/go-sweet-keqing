package keqing

import (
	"math"
	"reflect"
	"testing"
)

type collisionPair struct{ A, B string }
type semanticKey struct {
	ID      int
	Version int
}

func (v semanticKey) Equal(other semanticKey) bool { return v.ID == other.ID }

type exactKey int

func (v exactKey) Equal(other exactKey) bool { return v == other }

type mutableEntity struct {
	ID   int
	Tags []string
}

func TestCollectionHashCollisions(t *testing.T) {
	a, b := collisionPair{"a b", "c"}, collisionPair{"a", "b c"}
	s := NewSet[collisionPair](nil)
	if !s.Add(a) || !s.Add(b) || s.Size() != 2 {
		t.Fatalf("distinct values were merged: %v", s.GetData())
	}
	l := AsList([]collisionPair{a}, nil)
	if l.Contains(b) || l.IndexOf(b) != -1 || l.Remove(b) {
		t.Fatalf("absent colliding value matched: %v", l.GetData())
	}
	s.Remove(b)
	if !reflect.DeepEqual(s.GetData(), []collisionPair{a}) {
		t.Fatalf("removed wrong collision: %v", s.GetData())
	}
	l = AsList([]collisionPair{a, b}, nil)
	l.RetainAll(b)
	if !reflect.DeepEqual(l.GetData(), []collisionPair{b}) {
		t.Fatalf("RetainAll retained collision: %v", l.GetData())
	}
	s = AsSet([]collisionPair{a, b}, nil)
	s.RemoveAll(b)
	if !reflect.DeepEqual(s.GetData(), []collisionPair{a}) {
		t.Fatalf("RemoveAll removed collision: %v", s.GetData())
	}
	l = AsList([]collisionPair{a, b}, nil)
	l.RemoveAll(b)
	if !reflect.DeepEqual(l.GetData(), []collisionPair{a}) {
		t.Fatalf("List.RemoveAll removed collision: %v", l.GetData())
	}
}

func TestCollectionEqualOverridesHash(t *testing.T) {
	h := func(v semanticKey) uint64 { return uint64(v.Version) }
	a, b := semanticKey{ID: 1, Version: 1}, semanticKey{ID: 1, Version: 2}
	l := AsList([]semanticKey{a}, h)
	if !l.Contains(b) || l.IndexOf(b) != 0 || l.LastIndexOf(b) != 0 {
		t.Fatal("Equal values with distinct hashes were not found")
	}
	if !l.Remove(b) || !l.IsEmpty() {
		t.Fatal("Equal value was not removed")
	}
	l = AsList([]semanticKey{a}, h)
	l.RetainAll(b)
	if l.Size() != 1 {
		t.Fatal("Equal value was not retained")
	}
	l.RemoveAll(b)
	if !l.IsEmpty() {
		t.Fatal("Equal value was not bulk-removed")
	}
	s := AsSet([]semanticKey{a}, h)
	if s.Add(b) || !s.Contains(b) || s.Size() != 1 {
		t.Fatal("Equal values with distinct hashes were duplicated")
	}
	s.RetainAll(b)
	if s.Size() != 1 {
		t.Fatal("Set did not retain Equal value")
	}
	s.RemoveAll(b)
	if !s.IsEmpty() {
		t.Fatal("Set did not bulk-remove Equal value")
	}
	exactHash := func(exactKey) uint64 { return 1 }
	el := AsList([]exactKey{1}, exactHash)
	if el.Contains(2) {
		t.Fatal("Equal(false) was ignored")
	}
	es := AsSet([]exactKey{1}, exactHash)
	if !es.Add(2) || es.Size() != 2 {
		t.Fatal("Set ignored Equal(false)")
	}
}

func TestCollectionGoEqualityAndDynamicValues(t *testing.T) {
	s := NewSet[any](nil)
	for _, v := range []any{1, "1", []int{1, 2}, map[string]int{"x": 1}, nil, (*mutableEntity)(nil)} {
		if !s.Add(v) {
			t.Fatalf("lost distinct value %#v", v)
		}
	}
	if s.Size() != 6 || !s.Contains([]int{1, 2}) || !s.Contains(map[string]int{"x": 1}) {
		t.Fatalf("dynamic values were not preserved: %#v", s.GetData())
	}
	if s.Add([]int{1, 2}) {
		t.Fatal("DeepEqual dynamic slice was duplicated")
	}
	negativeZero := math.Copysign(0, -1)
	fs := NewSet[float64](nil)
	fs.Add(negativeZero)
	if fs.Add(0) || !fs.Contains(0) || fs.Size() != 1 {
		t.Fatal("Go-equal signed zeros were split by hash")
	}
}

func TestCollectionMutablePointers(t *testing.T) {
	for _, custom := range []bool{false, true} {
		var h func(*mutableEntity) uint64
		if custom {
			h = func(v *mutableEntity) uint64 {
				if v == nil {
					return 0
				}
				return uint64(v.ID)
			}
		}
		a := &mutableEntity{ID: 1, Tags: []string{"a"}}
		b := &mutableEntity{ID: 1, Tags: []string{"a"}}
		s := AsSet([]*mutableEntity{a, b}, h)
		l := AsList([]*mutableEntity{a, b}, h)
		if s.Size() != 2 {
			t.Fatal("distinct pointer identities merged")
		}
		a.ID = 9
		a.Tags[0] = "changed"
		if !s.Contains(a) || s.Add(a) || !l.Contains(a) || l.IndexOf(a) != 0 {
			t.Fatal("pointer mutation invalidated membership")
		}
		if !s.Remove(a) || !l.Remove(a) || s.Size() != 1 || l.Size() != 1 {
			t.Fatal("mutated pointer could not be removed")
		}
		if s.GetData()[0] != b || l.Get(0) != b {
			t.Fatal("removed wrong pointer")
		}
	}
}

func TestCollectionSetOperationsWithDifferentHashes(t *testing.T) {
	a := AsSet([]int{1, 2}, func(x int) uint64 { return uint64(x) })
	b := AsSet([]int{2, 3}, func(x int) uint64 { return uint64(x + 100) })
	if got := a.Intersect(b).GetData(); !reflect.DeepEqual(got, []int{2}) {
		t.Fatalf("intersection %v", got)
	}
	if got := a.Difference(b).GetData(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("difference %v", got)
	}
	if got := a.Union(b).GetData(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("union %v", got)
	}
}
