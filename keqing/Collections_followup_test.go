package keqing

import (
	"reflect"
	"testing"
)

func TestCollectionFollowupListComparisonBudget(t *testing.T) {
	const count = 1024
	for _, custom := range []bool{false, true} {
		var hash func(int) uint64
		if custom {
			hash = func(v int) uint64 { return uint64(v) }
		}
		list := NewList[int](hash)
		for repeat := 0; repeat < 2; repeat++ {
			for i := 0; i < count; i++ {
				list.Add(i)
			}
		}
		comparisons := 0
		equal := func(a, b int) bool { comparisons++; return a == b }
		for i := 0; i < count; i++ {
			if got := list.indexOfByLocked(i, equal, false); got != i {
				t.Fatalf("first index %d = %d", i, got)
			}
			if got := list.indexOfByLocked(i, equal, true); got != i+count {
				t.Fatalf("last index %d = %d", i, got)
			}
			if list.indexOfByLocked(i+count, equal, false) != -1 || list.indexOfByLocked(i+count, equal, true) != -1 {
				t.Fatal("absent value found")
			}
		}
		if comparisons > 4*count {
			t.Errorf("custom=%v: %d comparisons for %d values; List lookup scans unrelated values", custom, comparisons, count)
		}
	}
}

func TestCollectionFollowupBulkComparisonBudget(t *testing.T) {
	const count = 1024
	values := make([]int, count)
	items := make([]int, count)
	for i := range values {
		values[i], items[i] = i, count-1-i
	}
	for _, keep := range []bool{false, true} {
		comparisons := 0
		equal := func(a, b int) bool { comparisons++; return a == b }
		result := filterCollection(values, items, func(v int) uint64 { return defaultHash(v) }, false, keep, equal)
		if keep && !reflect.DeepEqual(result, values) || !keep && len(result) != 0 {
			t.Fatalf("keep=%v: wrong filtered data", keep)
		}
		if comparisons > 2*count {
			t.Errorf("keep=%v: %d comparisons for %d values; bulk filter scans unrelated candidates", keep, comparisons, count)
		}
	}
}

type followupEntity struct{ ID int }

func TestCollectionFollowupDerivedPointerStrategy(t *testing.T) {
	a, b, c := &followupEntity{1}, &followupEntity{2}, &followupEntity{3}
	source := AsSet([]*followupEntity{a, b}, nil)
	other := AsSet([]*followupEntity{b, c}, nil)
	derived := map[string]*Set[*followupEntity]{
		"copy":              source.Copy(),
		"filter":            source.Filter(func(*followupEntity) bool { return true }),
		"deep copy":         source.DeepCopy(func(v *followupEntity) *followupEntity { return v }),
		"union":             source.Union(other),
		"intersection":      source.Intersect(other),
		"difference":        source.Difference(other),
		"set-list-set":      source.ToList().ToSet(),
		"list-set":          AsList([]*followupEntity{a, b}, nil).ToSet(),
		"sublist-set":       AsList([]*followupEntity{a, b}, nil).SubList(0, 1).ToSet(),
		"list-copy-set":     AsList([]*followupEntity{a, b}, nil).Copy().ToSet(),
		"list-filter-set":   AsList([]*followupEntity{a, b}, nil).Filter(func(*followupEntity) bool { return true }).ToSet(),
		"list-deepcopy-set": AsList([]*followupEntity{a, b}, nil).DeepCopy(func(v *followupEntity) *followupEntity { return v }).ToSet(),
		"list-union-set":    AsList([]*followupEntity{a}, nil).Union(AsList([]*followupEntity{b}, nil)).ToSet(),
	}
	for name, got := range derived {
		t.Run(name, func(t *testing.T) {
			if got.scanOnHashMiss {
				t.Fatal("default pointer hash was reclassified as a custom mutable hash")
			}
			comparisons := 0
			index := got.indexOfByLocked(c, func(a, b *followupEntity) bool { comparisons++; return a == b })
			if (index >= 0) != got.Contains(c) {
				t.Fatal("counted lookup differs from public Contains")
			}
			if comparisons > 1 {
				t.Fatalf("default pointer lookup made %d comparisons", comparisons)
			}
		})
	}
}

func TestCollectionFollowupDerivedCustomReferenceStrategy(t *testing.T) {
	a := &followupEntity{1}
	hash := func(v *followupEntity) uint64 { return uint64(v.ID) }
	set := AsSet([]*followupEntity{a}, hash)
	list := AsList([]*followupEntity{a}, hash)
	derivedSets := []*Set[*followupEntity]{set.Copy(), set.ToList().ToSet(), list.Copy().ToSet(), list.SubList(0).ToSet()}
	derivedLists := []*List[*followupEntity]{list.Copy(), list.SubList(0), set.Copy().ToList()}
	a.ID = 99
	for _, got := range derivedSets {
		if !got.scanOnHashMiss || !got.Contains(a) || got.Add(a) || !got.Remove(a) {
			t.Fatal("derived Set lost mutable custom-hash fallback")
		}
	}
	for _, got := range derivedLists {
		if !got.Contains(a) || got.IndexOf(a) != 0 || got.LastIndexOf(a) != 0 || !got.Remove(a) {
			t.Fatal("derived List lost mutable custom-hash fallback")
		}
	}
}

type directionalEqualKey int

func (v directionalEqualKey) Equal(other directionalEqualKey) bool { return v < other }

func TestCollectionFollowupBulkEqualDirection(t *testing.T) {
	// Preserve item.Equal(candidate), even if a caller supplies a directional
	// predicate rather than an equivalence relation.
	list := AsList([]directionalEqualKey{3, 2, 1}, nil)
	list.RetainAll(2)
	if !reflect.DeepEqual(list.GetData(), []directionalEqualKey{1}) {
		t.Fatalf("RetainAll reversed Equal receiver: %v", list.GetData())
	}
	list = AsList([]directionalEqualKey{3, 2, 1}, nil)
	list.RemoveAll(2)
	if !reflect.DeepEqual(list.GetData(), []directionalEqualKey{3, 2}) {
		t.Fatalf("RemoveAll reversed Equal receiver: %v", list.GetData())
	}
	set := AsSet([]directionalEqualKey{3, 2, 1}, nil)
	set.RetainAll(2)
	if !reflect.DeepEqual(set.GetData(), []directionalEqualKey{1}) {
		t.Fatalf("Set.RetainAll reversed Equal receiver: %v", set.GetData())
	}
}

func TestCollectionFollowupCrossHashFirstLast(t *testing.T) {
	values := []semanticKey{{1, 1}, {2, 2}, {1, 2}, {1, 3}}
	list := AsList(values, func(v semanticKey) uint64 { return uint64(v.Version) })
	query := semanticKey{1, 2}
	if !list.Contains(query) || list.IndexOf(query) != 0 || list.LastIndexOf(query) != 3 {
		t.Fatal("cross-hash Equal lost first/last list order")
	}
	if !list.Remove(query) || !reflect.DeepEqual(list.GetData(), values[1:]) {
		t.Fatal("Remove did not remove the first Equal value")
	}
}

func TestCollectionFollowupBulkOrderAndCollisions(t *testing.T) {
	for _, hash := range []func(int) uint64{nil, func(int) uint64 { return 1 }} {
		list := AsList([]int{2, 1, 2, 3, 4}, hash)
		list.RetainAll(3, 2, 2)
		if !reflect.DeepEqual(list.GetData(), []int{2, 2, 3}) {
			t.Fatal("RetainAll changed duplicate/source order")
		}
		list.RemoveAll(2)
		if !reflect.DeepEqual(list.GetData(), []int{3}) {
			t.Fatal("RemoveAll lost collision equality")
		}
		set := AsSet([]int{2, 1, 3, 4}, hash)
		set.RetainAll(3, 2, 2)
		if !reflect.DeepEqual(set.GetData(), []int{2, 3}) {
			t.Fatal("Set.RetainAll changed source order")
		}
		set.RemoveAll(2)
		if !reflect.DeepEqual(set.GetData(), []int{3}) {
			t.Fatal("Set.RemoveAll lost collision equality")
		}
	}
}
