package keqing

import (
	"math"
	"testing"
)

func requirePanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	fn()
}

func TestMoneyExactConversion(t *testing.T) {
	for _, tc := range []struct {
		yuan string
		fen  int64
	}{
		{"0.29", 29}, {"-0.29", -29}, {"1.13", 113}, {"1.239", 123}, {"-1.239", -123},
		{"2.9e-1", 29}, {"90071992547409.93", 9007199254740993},
		{"92233720368547758.07", math.MaxInt64}, {"-92233720368547758.08", math.MinInt64},
	} {
		if got := Yuan2Fen(tc.yuan); got != tc.fen {
			t.Errorf("Yuan2Fen(%q)=%d want %d", tc.yuan, got, tc.fen)
		}
	}
	for _, fen := range []int64{0, 1, -1, 29, -29, 9007199254740993, math.MaxInt64, math.MinInt64} {
		if got := Yuan2Fen(Fen2Yuan(fen)); got != fen {
			t.Errorf("round trip %d=%d", fen, got)
		}
	}
}

func TestMoneyArithmetic(t *testing.T) {
	for _, tc := range []struct{ got, want string }{
		{YuanAdd("0.29", "0.01"), "0.30"}, {YuanSub("0.29", "0.01"), "0.28"},
		{YuanMul("0.29", 1), "0.29"}, {YuanDiv("0.29", 1), "0.29"},
		{YuanMul("0.03", 0.5), "0.01"}, {YuanDiv("-0.03", 2), "-0.01"},
	} {
		if tc.got != tc.want {
			t.Errorf("got %s want %s", tc.got, tc.want)
		}
	}
}

func TestMoneyRejectsInvalidAndOverflow(t *testing.T) {
	for _, input := range []string{"NaN", "+Inf", "-Inf", "bad", "1/2", "92233720368547758.08", "-92233720368547758.09"} {
		t.Run(input, func(t *testing.T) { requirePanic(t, func() { Yuan2Fen(input) }) })
	}
	for _, fn := range []func(){
		func() { YuanAdd("92233720368547758.07", "0.01") },
		func() { YuanSub("-92233720368547758.08", "0.01") },
		func() { YuanMul("1", math.NaN()) }, func() { YuanDiv("1", math.Inf(1)) },
		func() { YuanMul("92233720368547758.07", 2) }, func() { YuanDiv("1", 0) },
	} {
		requirePanic(t, fn)
	}
}
