package keqing

import (
	"testing"
	"time"
)

func midnightReviewLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return loc
}

func midnightReviewTime(t *testing.T, value string, loc *time.Location) time.Time {
	t.Helper()
	date, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatal(err)
	}
	return date.In(loc)
}

func TestCalendarRangesAcrossMidnightTransitions(t *testing.T) {
	cases := []struct {
		name, zone, now, nextNow, start, next string
	}{
		{
			"before missing midnight", "America/Santiago",
			"2026-09-05T12:00:00-04:00", "2026-09-06T12:00:00-03:00",
			"2026-09-05T00:00:00-04:00", "2026-09-06T01:00:00-03:00",
		},
		{
			"missing midnight", "America/Santiago",
			"2026-09-06T12:00:00-03:00", "2026-09-07T12:00:00-03:00",
			"2026-09-06T01:00:00-03:00", "2026-09-07T00:00:00-03:00",
		},
		{
			"repeated midnight", "America/Havana",
			"2026-11-01T12:00:00-05:00", "2026-11-02T12:00:00-05:00",
			"2026-11-01T00:00:00-04:00", "2026-11-02T00:00:00-05:00",
		},
		{
			"midnight normalized to second occurrence", "Asia/Gaza",
			"1996-09-20T12:00:00+02:00", "1996-09-21T12:00:00+02:00",
			"1996-09-20T00:00:00+03:00", "1996-09-21T00:00:00+02:00",
		},
		{
			"fractional midnight gap", "Africa/Monrovia",
			"1972-01-07T12:00:00Z", "1972-01-08T12:00:00Z",
			"1972-01-07T00:44:30Z", "1972-01-08T00:00:00Z",
		},
		{
			"almost repeated date", "Pacific/Kwajalein",
			"1969-09-30T12:00:00-12:00", "1969-10-01T12:00:00-12:00",
			"1969-09-30T00:00:00+11:00", "1969-10-01T00:00:00-12:00",
		},
		{
			"before skipped date", "Pacific/Apia",
			"2011-12-29T12:00:00-10:00", "",
			"2011-12-29T00:00:00-10:00", "2011-12-31T00:00:00+14:00",
		},
		{
			"after skipped date", "Pacific/Apia",
			"2011-12-31T12:00:00+14:00", "2012-01-01T12:00:00+14:00",
			"2011-12-31T00:00:00+14:00", "2012-01-01T00:00:00+14:00",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			loc := midnightReviewLocation(t, tc.zone)
			start := midnightReviewTime(t, tc.start, loc)
			next := midnightReviewTime(t, tc.next, loc)
			wantEnd := next.Add(-time.Nanosecond)
			gotStart, gotEnd := startAndEnd4Today(midnightReviewTime(t, tc.now, loc))
			if !gotStart.Equal(start) || !gotEnd.Equal(wantEnd) {
				t.Errorf("today = %s .. %s, want %s .. %s", gotStart, gotEnd, start, wantEnd)
			}
			if tc.nextNow != "" {
				gotStart, gotEnd = startAndEnd4Yesterday(midnightReviewTime(t, tc.nextNow, loc))
				if !gotStart.Equal(start) || !gotEnd.Equal(wantEnd) {
					t.Errorf("yesterday = %s .. %s, want %s .. %s", gotStart, gotEnd, start, wantEnd)
				}
			}
		})
	}
}

func TestCalendarSkippedYesterdayIsEmpty(t *testing.T) {
	loc := midnightReviewLocation(t, "Pacific/Apia")
	now := midnightReviewTime(t, "2011-12-31T12:00:00+14:00", loc)
	wantStart := midnightReviewTime(t, "2011-12-31T00:00:00+14:00", loc)
	start, end := startAndEnd4Yesterday(now)
	if !start.Equal(wantStart) || !end.Equal(start.Add(-time.Nanosecond)) {
		t.Fatalf("skipped 2011-12-30 must be an empty closed range, got %s .. %s", start, end)
	}
}

func TestCalendarMonthRangesAcrossMidnightTransitions(t *testing.T) {
	for _, tc := range []struct{ zone, now, start, next string }{
		{"America/Asuncion", "2017-09-15T12:00:00-04:00", "2017-09-01T00:00:00-04:00", "2017-10-01T01:00:00-03:00"},
		{"America/Asuncion", "2017-10-15T12:00:00-03:00", "2017-10-01T01:00:00-03:00", "2017-11-01T00:00:00-03:00"},
		{"America/Havana", "2026-11-15T12:00:00-05:00", "2026-11-01T00:00:00-04:00", "2026-12-01T00:00:00-05:00"},
		{"Pacific/Apia", "2011-12-15T12:00:00-10:00", "2011-12-01T00:00:00-10:00", "2012-01-01T00:00:00+14:00"},
	} {
		t.Run(tc.zone+tc.now, func(t *testing.T) {
			loc := midnightReviewLocation(t, tc.zone)
			wantStart := midnightReviewTime(t, tc.start, loc)
			wantEnd := midnightReviewTime(t, tc.next, loc).Add(-time.Nanosecond)
			start, end := startAndEnd4CurrentMonth(midnightReviewTime(t, tc.now, loc))
			if !start.Equal(wantStart) || !end.Equal(wantEnd) {
				t.Fatalf("month = %s .. %s, want %s .. %s", start, end, wantStart, wantEnd)
			}
		})
	}
}

func TestMonthDayCountDoesNotDependOnLocalTransitions(t *testing.T) {
	oldLocal := time.Local
	t.Cleanup(func() { time.Local = oldLocal })
	time.Local = midnightReviewLocation(t, "America/Asuncion")
	if got := GetDaysInMonth(time.September, 2017); got != 30 {
		t.Fatalf("September 2017 has %d days, want 30", got)
	}
	for _, tc := range []struct{ year, want int }{{0, 29}, {-400, 29}, {-100, 28}, {2000, 29}, {2100, 28}} {
		if got := GetDaysInMonth(time.February, tc.year); got != tc.want {
			t.Errorf("February %d has %d days, want %d", tc.year, got, tc.want)
		}
	}
}

func TestCalendarRangesInFixedZones(t *testing.T) {
	for _, loc := range []*time.Location{time.UTC, time.FixedZone("UTC+05:45", 5*3600+45*60)} {
		now := time.Date(2026, 1, 1, 12, 0, 0, 0, loc)
		wantStart := time.Date(2026, 1, 1, 0, 0, 0, 0, loc)
		wantNext := time.Date(2026, 1, 2, 0, 0, 0, 0, loc)
		start, end := startAndEnd4Today(now)
		if !start.Equal(wantStart) || !end.Equal(wantNext.Add(-time.Nanosecond)) {
			t.Errorf("%s today = %s .. %s", loc, start, end)
		}
		wantYesterday := time.Date(2025, 12, 31, 0, 0, 0, 0, loc)
		start, end = startAndEnd4Yesterday(now)
		if !start.Equal(wantYesterday) || !end.Equal(wantStart.Add(-time.Nanosecond)) {
			t.Errorf("%s yesterday = %s .. %s", loc, start, end)
		}
	}
}

func TestCalendarRangesAcrossFutureYearEnd(t *testing.T) {
	for _, name := range []string{"America/Santiago", "America/Havana"} {
		loc := midnightReviewLocation(t, name)
		now := time.Date(2040, time.December, 31, 12, 0, 0, 0, loc)
		wantStart := time.Date(2040, time.December, 31, 0, 0, 0, 0, loc)
		wantNext := time.Date(2041, time.January, 1, 0, 0, 0, 0, loc)
		start, end := startAndEnd4Today(now)
		if !start.Equal(wantStart) || !end.Equal(wantNext.Add(-time.Nanosecond)) {
			t.Errorf("%s future today = %s .. %s", name, start, end)
		}
		start, end = startAndEnd4Yesterday(wantNext.Add(12 * time.Hour))
		if !start.Equal(wantStart) || !end.Equal(wantNext.Add(-time.Nanosecond)) {
			t.Errorf("%s future yesterday = %s .. %s", name, start, end)
		}
		start, end = startAndEnd4CurrentMonth(now)
		wantMonthStart := time.Date(2040, time.December, 1, 0, 0, 0, 0, loc)
		if !start.Equal(wantMonthStart) || !end.Equal(wantNext.Add(-time.Nanosecond)) {
			t.Errorf("%s future month = %s .. %s", name, start, end)
		}
	}
}
