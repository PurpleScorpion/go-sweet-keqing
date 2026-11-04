package keqing

import (
	"math"
	"testing"
	"time"
	_ "time/tzdata"
)

func TestIssue13CalendarDays(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, date := range []time.Time{time.Date(2026, 3, 8, 0, 0, 0, 0, loc), time.Date(2026, 11, 1, 0, 0, 0, 0, loc)} {
		want := date.AddDate(0, 0, 1)
		if got := DateAddDay(date, 1); !got.Equal(want) {
			t.Errorf("add: got %s want %s", got, want)
		}
		if got := DateSubDay(want, 1); !got.Equal(date) {
			t.Errorf("sub: got %s want %s", got, date)
		}
	}
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if got, want := DateAddDay(date, 106752), date.AddDate(0, 0, 106752); !got.Equal(want) {
		t.Errorf("large day count: %s != %s", got, want)
	}
	if got := DateAddDay(date, -1); !got.Equal(DateSubDay(date, 1)) {
		t.Fatal("negative day addition differs")
	}
	for _, fn := range []func(){func() { DateAddDay(date, math.MaxInt64) }, func() { DateAddDay(date, math.MinInt64) }, func() { DateSubDay(date, math.MinInt64) }} {
		valuesMustPanic(t, fn)
	}
}

func TestIssue13DurationRange(t *testing.T) {
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		unit     time.Duration
		add, sub func(time.Time, int64) time.Time
	}{
		{"second", time.Second, DateAddSecond, DateSubSecond},
		{"minute", time.Minute, DateAddMinute, DateSubMinute},
		{"hour", time.Hour, DateAddHour, DateSubHour},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			limit := int64(math.MaxInt64) / int64(tc.unit)
			if got, want := tc.add(date, limit), date.Add(time.Duration(limit)*tc.unit); !got.Equal(want) {
				t.Fatal(got, want)
			}
			if got, want := tc.sub(date, 2), date.Add(-2*tc.unit); !got.Equal(want) {
				t.Fatal(got, want)
			}
			if got, want := tc.add(date, -2), tc.sub(date, 2); !got.Equal(want) {
				t.Fatal(got, want)
			}
			valuesMustPanic(t, func() { tc.add(date, limit+1) })
			valuesMustPanic(t, func() { tc.add(date, -limit-1) })
			valuesMustPanic(t, func() { tc.sub(date, math.MinInt64) })
		})
	}
}

func TestIssue13CalendarRanges(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		month time.Month
		day   int
		hours int
	}{{3, 8, 23}, {11, 1, 25}} {
		start := time.Date(2026, tc.month, tc.day, 0, 0, 0, 0, loc)
		next := start.AddDate(0, 0, 1)
		wantEnd := next.Add(-time.Nanosecond)
		// Deliberately avoid changing time.Local: helpers must use the supplied time's location.
		todayStart, todayEnd := startAndEnd4Today(start.Add(12 * time.Hour))
		if !todayStart.Equal(start) || !todayEnd.Equal(wantEnd) {
			t.Errorf("today: %s .. %s want %s .. %s", todayStart, todayEnd, start, wantEnd)
		}
		yesterdayStart, yesterdayEnd := startAndEnd4Yesterday(next.Add(30 * time.Minute))
		if !yesterdayStart.Equal(start) || !yesterdayEnd.Equal(wantEnd) {
			t.Errorf("yesterday: %s .. %s want %s .. %s", yesterdayStart, yesterdayEnd, start, wantEnd)
		}
		if got := todayEnd.Sub(todayStart) + time.Nanosecond; got != time.Duration(tc.hours)*time.Hour {
			t.Errorf("day length: %s", got)
		}
		monthStart, monthEnd := startAndEnd4CurrentMonth(start)
		expectedMonthStart := time.Date(2026, tc.month, 1, 0, 0, 0, 0, loc)
		expectedMonthEnd := expectedMonthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
		if !monthStart.Equal(expectedMonthStart) || !monthEnd.Equal(expectedMonthEnd) {
			t.Errorf("month: %s .. %s", monthStart, monthEnd)
		}
	}
}
