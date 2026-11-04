package keqing

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestIssue09StringFormatting(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"nil", nil, "null"},
		{"uint", []uint64{math.MaxUint64}, "[18446744073709551615]"},
		{"nil map", map[string]interface{}{"v": nil}, "{v: null}"},
		{"nil slice", []interface{}{nil}, "[null]"},
		{"nil field", struct{ Value interface{} }{}, "{Value: null}"},
		{"private", struct {
			Public  int
			private string
		}{1, "secret"}, "{Public: 1}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ToString(tc.value); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
	if got := ToString(struct{ At time.Time }{time.Unix(0, 0)}); !strings.HasPrefix(got, "{At:") {
		t.Fatal(got)
	}
	if got := ToString([]int{-1, 2}); got != "[-1, 2]" {
		t.Fatal(got)
	}
}

func TestIssue09ArrayContainsInterfaces(t *testing.T) {
	for _, tc := range []struct {
		arr   interface{}
		value interface{}
		want  bool
	}{
		{[]interface{}{1, "ok", nil}, 1, true}, {[]interface{}{nil}, nil, true},
		{[]interface{}{[]int{1}}, []int{1}, true}, {[]int{1}, nil, false},
		{nil, nil, false}, {[]interface{}{int64(1)}, 1, false},
	} {
		if got := ArrayContains(tc.arr, tc.value); got != tc.want {
			t.Fatalf("ArrayContains(%#v,%#v)=%v", tc.arr, tc.value, got)
		}
	}
}

func TestIssue09NamedStringSort(t *testing.T) {
	type label string
	got := []label{"z", "a", "m"}
	SortArray(got)
	if !reflect.DeepEqual(got, []label{"a", "m", "z"}) {
		t.Fatal(got)
	}
	SortArrayDesc(&got)
	if !reflect.DeepEqual(got, []label{"z", "m", "a"}) {
		t.Fatal(got)
	}
}

type Issue09Embedded struct{ Name string }
type issue09Outer struct{ *Issue09Embedded }

func TestIssue09CopyPropertiesValidation(t *testing.T) {
	source := struct{ Name string }{"ok"}
	var nilDestination *struct{ Name string }
	number := 0
	for _, dst := range []interface{}{nil, nilDestination, source, &number} {
		if err := CopyProperties(dst, source); err == nil {
			t.Errorf("expected invalid destination error for %T", dst)
		}
	}
	dst := issue09Outer{}
	if err := CopyProperties(&dst, source); err == nil {
		t.Error("nil embedded pointer should return error")
	}
	dst.Issue09Embedded = &Issue09Embedded{}
	if err := CopyProperties(&dst, source); err != nil || dst.Name != "ok" {
		t.Fatalf("promoted copy: %+v, %v", dst, err)
	}
	out := struct {
		Name   string
		Count  int
		hidden int
	}{"keep", 0, 7}
	in := struct {
		Name   string
		Count  int
		hidden int
	}{"replace", 2, 9}
	if err := CopyProperties(&out, &in, "Name"); err != nil {
		t.Fatal(err)
	}
	if out.Name != "keep" || out.Count != 2 || out.hidden != 7 {
		t.Fatalf("ignore/private fields: %+v", out)
	}
	if err := CopyProperties(&out, nil); err == nil {
		t.Error("expected nil source error")
	}
}

func TestIssue09ToArrayNilPointers(t *testing.T) {
	var ptr *[]int
	if got := ToArray[int](ptr); got != nil {
		t.Fatal(got)
	}
	n := 7
	if got := ToArray[int]([]*int{nil, &n}); !reflect.DeepEqual(got, []int{0, 7}) {
		t.Fatal(got)
	}
	values := []int{1, 2}
	if got := ToArray[int](&values); !reflect.DeepEqual(got, values) {
		t.Fatal(got)
	}
	if got := ToArray[int](values); !reflect.DeepEqual(got, values) {
		t.Fatal(got)
	}
}
