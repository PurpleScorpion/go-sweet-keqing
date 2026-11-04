package keqing

import (
	"math"
	"reflect"
	"strconv"
	"testing"
)

func valuesMustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic for invalid input")
		}
	}()
	fn()
}

func TestIssue08IntegerParsing(t *testing.T) {
	if ParseInt8("127") != 127 || ParseInt8("-128") != -128 || ParseInt16("32767") != 32767 || ParseInt32("2147483647") != 2147483647 || ParseInt64("9223372036854775807") != math.MaxInt64 {
		t.Fatal("integer boundary value changed")
	}
	for name, fn := range map[string]func(){
		"int8 high": func() { ParseInt8("128") }, "int8 low": func() { ParseInt8("-129") },
		"int16": func() { ParseInt16("65536") }, "int32": func() { ParseInt32("2147483648") },
		"int64": func() { ParseInt64("9223372036854775808") }, "syntax": func() { ParseInt("bad") },
	} {
		t.Run(name, func(t *testing.T) { valuesMustPanic(t, fn) })
	}
	if strconv.IntSize == 32 {
		valuesMustPanic(t, func() { ParseInt("2147483648") })
	}
}

func TestIssue08FloatParsing(t *testing.T) {
	if ParseFloat32("0") != 0 || ParseFloat64("1.25") != 1.25 {
		t.Fatal("valid float parsing changed")
	}
	for name, fn := range map[string]func(){
		"float32 syntax": func() { ParseFloat32("bad") }, "float64 syntax": func() { ParseFloat64("") },
		"float32 overflow": func() { ParseFloat32("1e100") }, "float64 overflow": func() { ParseFloat64("1e999") },
	} {
		t.Run(name, func(t *testing.T) { valuesMustPanic(t, fn) })
	}
	// Standard explicit non-finite literals remain supported; syntax/range errors must not be discarded.
	if !math.IsNaN(ParseFloat64("NaN")) || !math.IsInf(ParseFloat64("+Inf"), 1) {
		t.Fatal("explicit special float literals changed")
	}
}

func TestIssue08YAMLIntegerRange(t *testing.T) {
	oldParent, oldChild := applicationYml, applicationChildYml
	t.Cleanup(func() { applicationYml, applicationChildYml = oldParent, oldChild })
	read := func(v interface{}) {
		t.Helper()
		if err := LoadYml(map[string]interface{}{"v": v}, map[string]interface{}{}); err != nil {
			t.Fatal(err)
		}
	}
	read(int64(3000000000))
	if got := ValueInt64("${v}"); got != 3000000000 {
		t.Fatalf("int64 narrowed to int: %d", got)
	}
	read(int64(math.MaxInt64))
	if got := ValueInt64("${v}"); got != math.MaxInt64 {
		t.Fatal(got)
	}
	read(int64(math.MinInt64))
	if got := ValueInt64("${v}"); got != math.MinInt64 {
		t.Fatal(got)
	}
	for _, v := range []interface{}{int64(2147483648), float64(2147483648), math.Inf(1), math.NaN()} {
		read(v)
		valuesMustPanic(t, func() { ValueInt32("${v}") })
	}
	for _, v := range []interface{}{math.Ldexp(1, 63), -math.Ldexp(1, 63) * 2, math.Inf(-1), math.NaN()} {
		read(v)
		valuesMustPanic(t, func() { ValueInt64("${v}") })
	}
	read(-1.9)
	if ValueInt64("${v}") != -1 {
		t.Fatal("finite truncation changed")
	}
	read(nil)
	if ValueInt64("${v}") != 0 {
		t.Fatal("missing values should remain zero")
	}
	if strconv.IntSize == 32 {
		read(int64(3000000000))
		valuesMustPanic(t, func() { ValueInt("${v}") })
	}
}

func TestIssue08Int64SortPreservesValues(t *testing.T) {
	values := []int64{math.MaxInt64, 3000000000, -3000000000, math.MinInt64, 0}
	SortArray(values)
	want := []int64{math.MinInt64, -3000000000, 0, 3000000000, math.MaxInt64}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("ascending: %v", values)
	}
	SortArrayDesc(&values)
	want = []int64{math.MaxInt64, 3000000000, 0, -3000000000, math.MinInt64}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("descending: %v", values)
	}
}

func TestIssue08RandomWideInterval(t *testing.T) {
	max := int(^uint(0) >> 1)
	min := -max - 1
	for i := 0; i < 256; i++ {
		got := RandomNextInt(min, max)
		if got < min || got >= max {
			t.Fatalf("%d outside [%d,%d)", got, min, max)
		}
	}
	valuesMustPanic(t, func() { RandomNextInt(1, 1) })
	valuesMustPanic(t, func() { RandomNextInt(2, 1) })
	if strconv.IntSize == 32 {
		valuesMustPanic(t, func() { RandomNextInt(int64(3000000000)) })
	}
}

func TestIssue08YAMLFloat32Range(t *testing.T) {
	oldParent, oldChild := applicationYml, applicationChildYml
	t.Cleanup(func() { applicationYml, applicationChildYml = oldParent, oldChild })
	read := func(value float64) float32 {
		t.Helper()
		if err := LoadYml(map[string]interface{}{"v": value}, map[string]interface{}{}); err != nil {
			t.Fatal(err)
		}
		return ValueFloat32("${v}")
	}
	for _, value := range []float64{0, 1.25, -1.25, math.MaxFloat32, -math.MaxFloat32} {
		if got := read(value); got != float32(value) {
			t.Errorf("ValueFloat32(%g) = %g, want %g", value, got, float32(value))
		}
	}
	// Values just above MaxFloat32 that still round to a finite float32 are valid.
	if got := read(math.Nextafter(float64(math.MaxFloat32), math.Inf(1))); got != math.MaxFloat32 {
		t.Errorf("finite boundary rounding = %g", got)
	}
	for _, value := range []float64{1e100, -1e100, math.MaxFloat64, -math.MaxFloat64, math.Ldexp(1, 128), -math.Ldexp(1, 128)} {
		valuesMustPanic(t, func() { read(value) })
	}
	// Explicit non-finite values retain the general ParseFloat contract.
	if !math.IsNaN(float64(read(math.NaN()))) {
		t.Error("explicit NaN changed")
	}
	for _, sign := range []int{1, -1} {
		if got := read(math.Inf(sign)); !math.IsInf(float64(got), sign) {
			t.Errorf("explicit Inf(%d) changed: %g", sign, got)
		}
	}
}
