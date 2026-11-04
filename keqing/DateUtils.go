package keqing

import (
	"fmt"
	"time"
)

var (
	DEFAULT_LOCAL_FORMAT = "2006-01-02 15:04:05"
	LOCAL_FORMAT_DATE    = "2006-01-02"
	LOCAL_FORMAT_TIME    = "15:04:05"
	DEFAULT_UTC_FORMAT   = "2006-01-02T15:04:05.999999Z"
	UTC_FORMAT_SHORT     = "2006-01-02T15:04:05Z"
	UTC                  = time.UTC
	LOCAL                = time.Local
)

type TimeVO struct {
	Year        int    //年
	Month       int    //月
	Day         int    //日
	Hour        int    //时
	Minute      int    //分
	Second      int    //秒
	Millisecond int    //毫秒
	YearStr     string //年
	// 以下不足10时会补0
	MonthStr       string //月
	DayStr         string //日
	HourStr        string //时
	MinuteStr      string //分
	SecondStr      string //秒
	MillisecondStr string //毫秒
}

/*
获取日期各个位置的详细信息
*/
func DateInfo(date time.Time) TimeVO {
	return TimeVO{
		Year:           date.Year(),
		Month:          int(date.Month()),
		Day:            date.Day(),
		Hour:           date.Hour(),
		Minute:         date.Minute(),
		Second:         date.Second(),
		Millisecond:    date.Nanosecond() / int(time.Millisecond),
		YearStr:        fmt.Sprintf("%d", date.Year()),
		MonthStr:       fmt.Sprintf("%02d", date.Month()),
		DayStr:         fmt.Sprintf("%02d", date.Day()),
		HourStr:        fmt.Sprintf("%02d", date.Hour()),
		MinuteStr:      fmt.Sprintf("%02d", date.Minute()),
		SecondStr:      fmt.Sprintf("%02d", date.Second()),
		MillisecondStr: fmt.Sprintf("%03d", date.Nanosecond()/int(time.Millisecond)),
	}
}

/*
获取当前时间
*/
func NowDate() time.Time {
	return time.Now()
}

/*
获取当前时间-UTC
*/
func NowUTCDate() time.Time {
	return time.Now().UTC()
}

/*
将当前时间转为 2006-01-02 15:04:05 格式的时间字符串
*/
func NowDateStr() string {
	return NowDate().Format(DEFAULT_LOCAL_FORMAT)
}

/*
将当前时间-UTC , 转为 2006-01-02T15:04:05.999999Z 格式的时间字符串
*/
func NowUTCDateStr() string {
	return NowUTCDate().Format(DEFAULT_UTC_FORMAT)
}

func FormatDate(date time.Time, format string) string {
	return date.Format(format)
}

func ParseDate(dateStr string, format string, location *time.Location) time.Time {
	if location == nil {
		location = time.UTC
	}
	date, err := time.ParseInLocation(format, dateStr, location)
	if err != nil {
		panic("time parse error")
	}
	return date
}

func DateAddSecond(date time.Time, seconds int64) time.Time {
	return date.Add(checkedDateDuration(seconds, time.Second))
}

func DateAddMinute(date time.Time, minute int64) time.Time {
	return date.Add(checkedDateDuration(minute, time.Minute))
}

func DateAddHour(date time.Time, hour int64) time.Time {
	return date.Add(checkedDateDuration(hour, time.Hour))
}

// DateAddDay 按日期所在时区增加日历日；不可表示的偏移会 panic。
func DateAddDay(date time.Time, day int64) time.Time {
	// AddDate accepts an int day offset. Check conversion and its day-of-month addition.
	const maxDayOffset = int64(^uint64(0)>>1) / (24 * 60 * 60)
	maxInt := int64(int(^uint(0) >> 1))
	if day < -maxDayOffset || day > maxDayOffset || int64(int(day)) != day || day > maxInt-int64(date.Day()) {
		panic("day offset is outside the supported calendar range")
	}
	result := date.AddDate(0, 0, int(day))
	if (day > 0 && result.Before(date)) || (day < 0 && result.After(date)) {
		panic("calendar date overflow")
	}
	return result
}

func DateSubSecond(date time.Time, seconds int64) time.Time {
	return date.Add(-checkedDateDuration(seconds, time.Second))
}
func DateSubMinute(date time.Time, minute int64) time.Time {
	return date.Add(-checkedDateDuration(minute, time.Minute))
}

func DateSubHour(date time.Time, hour int64) time.Time {
	return date.Add(-checkedDateDuration(hour, time.Hour))
}

func DateSubDay(date time.Time, day int64) time.Time {
	const minInt64 = -int64(^uint64(0)>>1) - 1
	if day == minInt64 {
		panic("day offset is outside the supported calendar range")
	}
	return DateAddDay(date, -day)
}

func ParseUTC(utc string) time.Time {
	return ParseDate(utc, DEFAULT_UTC_FORMAT, UTC)
}

func ParseLocal(local string) time.Time {
	return ParseDate(local, DEFAULT_LOCAL_FORMAT, LOCAL)
}

func UTC2Local(utc string) string {
	date := ParseDate(utc, DEFAULT_UTC_FORMAT, UTC)
	local := date.Local()
	return FormatDate(local, DEFAULT_LOCAL_FORMAT)
}

// 获取当前月份的天数
func GetDaysInMonth(data any, years ...int) int {
	year, month, err := getMonthFromData(data)
	if err != nil {
		panic(err)
	}
	if len(years) > 0 {
		year = years[0]
	} else if _, fromTime := data.(time.Time); !fromTime {
		panic("please provide the year")
	}

	// Gregorian month lengths do not depend on the machine's time zone.
	switch month {
	case time.February:
		if year%400 == 0 || year%4 == 0 && year%100 != 0 {
			return 29
		}
		return 28
	case time.April, time.June, time.September, time.November:
		return 30
	default:
		return 31
	}
}
func getMonthFromData(data any) (int, time.Month, error) {
	var month int
	switch t := data.(type) {
	case int:
		month = t
	case int16:
		month = int(t)
	case int32:
		month = int(t)
	case int64:
		if t < 1 || t > 12 {
			return 0, 0, fmt.Errorf("month must be between 1 and 12: %d", t)
		}
		month = int(t)
	case time.Month:
		month = int(t)
	case time.Time:
		return t.Year(), t.Month(), nil
	default:
		return 0, 0, fmt.Errorf("The value is of an unknown type: %T", t)
	}

	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("month must be between 1 and 12: %d", month)
	}

	return 0, time.Month(month), nil
}

func UTC2LocalCustom(utc string, utcFormat string, localFormat string) string {
	date := ParseDate(utc, utcFormat, UTC)
	local := date.Local()
	return FormatDate(local, localFormat)
}

func Local2UTC(local string) string {
	date := ParseDate(local, DEFAULT_LOCAL_FORMAT, LOCAL)
	utc := date.UTC()
	return FormatDate(utc, DEFAULT_UTC_FORMAT)
}

func Local2UTCCustom(local string, localFormat string, utcFormat string) string {
	date := ParseDate(local, localFormat, LOCAL)
	utc := date.UTC()
	return FormatDate(utc, utcFormat)
}

/*
获取当前系统的毫秒级时间戳
*/
func CurrentTimeMillis() int64 {
	// 获取当前时间的纳秒级时间戳
	nanosSinceEpoch := NowDate().UnixNano()

	// 将纳秒转换为毫秒
	millisecondsSinceEpoch := nanosSinceEpoch / int64(time.Millisecond)

	return millisecondsSinceEpoch
}

/*
获取当前系统的秒级时间戳
*/
func CurrentTimeSeconds() int64 {
	// 获取当前时间的秒级时间戳
	return NowDate().Unix()
}

/*
获取昨天的开始和结束时间（闭区间，结束时间为当天最后一纳秒）
*/
func GetStartAndEnd4Yesterday() (time.Time, time.Time) {
	return startAndEnd4Yesterday(time.Now())
}

func startAndEnd4Yesterday(now time.Time) (time.Time, time.Time) {
	today := calendarDate(now)
	todayStart := calendarDateStart(today, now.Location())
	yesterdayStart := calendarDateStart(today.AddDate(0, 0, -1), now.Location())
	// A wholly skipped civil date has no instants: end is one nanosecond before start.
	return yesterdayStart, todayStart.Add(-time.Nanosecond)
}

/*
获取昨天的开始和结束时间-UTC
*/
func GetStartAndEndUTC4Yesterday() (time.Time, time.Time) {
	yesterdayStart, yesterdayEnd := GetStartAndEnd4Yesterday()
	return yesterdayStart.UTC(), yesterdayEnd.UTC()
}

/*
获取今天的开始和结束时间（闭区间，结束时间为当天最后一纳秒）
*/
func GetStartAndEnd4Today() (time.Time, time.Time) {
	return startAndEnd4Today(time.Now())
}

func startAndEnd4Today(now time.Time) (time.Time, time.Time) {
	today := calendarDate(now)
	start := calendarDateStart(today, now.Location())
	next := calendarDateStart(today.AddDate(0, 0, 1), now.Location())
	return start, next.Add(-time.Nanosecond)
}

/*
获取今天的开始和结束时间-UTC
*/
func GetStartAndEndUTC4Today() (time.Time, time.Time) {
	start, now := GetStartAndEnd4Today()
	return start.UTC(), now.UTC()
}

/*
获取本月的开始和结束时间（闭区间，结束时间为当月最后一纳秒）
*/
func GetStartAndEnd4CurrentMonth() (time.Time, time.Time) {
	return startAndEnd4CurrentMonth(time.Now())
}

func startAndEnd4CurrentMonth(now time.Time) (time.Time, time.Time) {
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	start := calendarDateStart(month, now.Location())
	next := calendarDateStart(month.AddDate(0, 1, 0), now.Location())
	return start, next.Add(-time.Nanosecond)
}

// calendarDate represents a civil date in UTC for calendar arithmetic only.
// Arithmetic on an already normalized local midnight can change the intended date.
func calendarDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

// calendarDateStart returns the first valid instant at or after a local civil
// date's midnight. It chooses the first occurrence of a repeated midnight and
// advances past a missing midnight (or a wholly skipped date).
func calendarDateStart(date time.Time, loc *time.Location) time.Time {
	cursor := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	// time.Date may choose the second occurrence of an ambiguous midnight.
	// Walk earlier zone intervals while they still reach the requested civil date.
	for {
		start, _ := cursor.ZoneBounds()
		if start.IsZero() {
			break
		}
		previous := start.Add(-time.Nanosecond)
		if calendarDate(previous).Before(date) {
			break
		}
		cursor = previous
	}
	for {
		start, end := cursor.ZoneBounds()
		_, offset := cursor.Zone()
		// A fixed offset has no missing/ambiguous clock times. Avoid converting
		// offsets to time.Duration, which would unnecessarily limit large fixed zones.
		candidate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.FixedZone("", offset))
		if !start.IsZero() && candidate.Before(start) {
			candidate = start
		}
		actual := candidate.In(loc)
		_, actualOffset := actual.Zone()
		if actualOffset == offset {
			// Verify the candidate against the actual zone rather than relying on
			// ZoneBounds' synthetic POSIX year-end boundary (which may be stale
			// on the last day of a future leap year).
			return actual
		}
		if end.IsZero() || !end.After(cursor) {
			panic("time zone boundary does not advance")
		}
		// The requested midnight is beyond this interval. Crossing the transition
		// finds the first valid instant when a forward jump skipped midnight.
		cursor = end
	}
}

/*
获取本月的开始和结束时间-UTC
*/
func GetStartAndEndUTC4CurrentMonth() (time.Time, time.Time) {
	start, now := GetStartAndEnd4CurrentMonth()
	return start.UTC(), now.UTC()
}

// checkedDateDuration rejects multiplication overflow before converting a unit count.
func checkedDateDuration(amount int64, unit time.Duration) time.Duration {
	const maxInt64 = int64(^uint64(0) >> 1)
	const minInt64 = -maxInt64 - 1
	scale := int64(unit)
	if amount > maxInt64/scale || amount < minInt64/scale {
		panic("time duration overflow")
	}
	return time.Duration(amount * scale)
}
