package keqing

import (
	"testing"
	"time"
)

func TestIssue12MonthYearOverride(t *testing.T) {
	for _, tc := range []struct{ baseYear, override, want int }{{2023, 2024, 29}, {2024, 2023, 28}, {2024, 2100, 28}, {2023, 2000, 29}} {
		date := time.Date(tc.baseYear, 2, 1, 0, 0, 0, 0, time.UTC)
		if got := GetDaysInMonth(date, tc.override); got != tc.want {
			t.Errorf("%d -> %d: got %d want %d", tc.baseYear, tc.override, got, tc.want)
		}
	}
	date := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	if got := GetDaysInMonth(date); got != 29 {
		t.Fatal(got)
	}
	date = time.Date(0, 2, 1, 0, 0, 0, 0, time.UTC)
	if got := GetDaysInMonth(date); got != 29 {
		t.Fatalf("year zero: %d", got)
	}
}

func TestIssue12MonthValidation(t *testing.T) {
	for _, month := range []interface{}{0, 13, int16(13), int32(0), int64(13), int64(1<<32) + 2, time.Month(0), time.Month(13)} {
		t.Run("invalid", func(t *testing.T) { valuesMustPanic(t, func() { GetDaysInMonth(month, 2024) }) })
	}
	valuesMustPanic(t, func() { GetDaysInMonth(time.February) })
	for _, month := range []interface{}{2, int16(2), int32(2), int64(2), time.February} {
		if got := GetDaysInMonth(month, 2024); got != 29 {
			t.Errorf("%T: %d", month, got)
		}
	}
	if GetDaysInMonth(time.January, 2024) != 31 || GetDaysInMonth(time.December, 2024) != 31 {
		t.Fatal("month edge values changed")
	}
}
