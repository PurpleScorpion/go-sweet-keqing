package keqing

import (
	"reflect"
	"testing"
)

func TestIssue15LiteralNullStringOperations(t *testing.T) {
	for _, tc := range []struct {
		text, separator string
		want            []string
	}{
		{"NULL", "L", []string{"NU", "", ""}}, {"null", "l", []string{"nu", "", ""}},
		{"Null", "l", []string{"Nu", "", ""}}, {"null", "", []string{"n", "u", "l", "l"}},
		{"", "x", []string{""}}, {"a,b,c", ",", []string{"a", "b", "c"}}, {"null", "x", []string{"null"}},
	} {
		if got := Split(tc.text, tc.separator); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Split(%q,%q)=%#v want %#v", tc.text, tc.separator, got, tc.want)
		}
	}
	for _, tc := range []struct{ text, want string }{{"Null", "null"}, {"NULL", "n_u_l_l"}, {"null", "null"}, {"", ""}, {"HelloWorld", "hello_world"}} {
		if got := Camel2Snake(tc.text); got != tc.want {
			t.Errorf("Camel2Snake(%q)=%q want %q", tc.text, got, tc.want)
		}
	}
	if Camel2Hyphen("Null") != "null" {
		t.Fatal("Camel2Hyphen still inherits null sentinel")
	}
	if !IsEmpty("null") || !IsEmpty("NULL") {
		t.Fatal("IsEmpty's existing sentinel contract must remain unchanged")
	}
}
