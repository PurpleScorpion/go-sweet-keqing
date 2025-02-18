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
	return date.Add(time.Duration(seconds) * time.Second)
}

func DateAddMinute(date time.Time, minute int64) time.Time {
	return date.Add(time.Duration(minute) * time.Minute)
}

func DateAddHour(date time.Time, hour int64) time.Time {
	return date.Add(time.Duration(hour) * time.Hour)
}

func DateAddDay(date time.Time, day int64) time.Time {
	return date.Add(time.Duration(day) * time.Hour * 24)
}

func DateSubSecond(date time.Time, seconds int64) time.Time {
	return DateAddSecond(date, seconds*-1)
}
func DateSubMinute(date time.Time, minute int64) time.Time {
	return DateAddMinute(date, minute*-1)
}

func DateSubHour(date time.Time, hour int64) time.Time {
	return DateAddHour(date, hour*-1)
}

func DateSubDay(date time.Time, day int64) time.Time {
	return DateAddDay(date, day*-1)
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
获取昨天的开始和结束时间
*/
func GetStartAndEnd4Yesterday() (time.Time, time.Time) {
	now := time.Now()
	// 减去24小时得到昨天同一时间
	yesterdaySameTime := now.Add(-24 * time.Hour)
	// 设置时间为昨天0点
	yesterdayStart := time.Date(yesterdaySameTime.Year(), yesterdaySameTime.Month(), yesterdaySameTime.Day(), 0, 0, 0, 0, time.Local)
	// 计算昨天24点的时间
	yesterdayEnd := yesterdayStart.Add(24 * time.Hour)
	return yesterdayStart, yesterdayEnd
}

/*
获取昨天的开始和结束时间-UTC
*/
func GetStartAndEndUTC4Yesterday() (time.Time, time.Time) {
	yesterdayStart, yesterdayEnd := GetStartAndEnd4Yesterday()
	return yesterdayStart.UTC(), yesterdayEnd.UTC()
}

/*
获取今天的开始和结束时间
*/
func GetStartAndEnd4Today() (time.Time, time.Time) {
	now := time.Now()
	// 设置时间为昨天0点
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// 获取明天的开始时间
	nextDay := start.AddDate(0, 0, 1)
	// 获取今天的结束时间
	end := nextDay.Add(-time.Nanosecond)
	return start, end
}

/*
获取今天的开始和结束时间-UTC
*/
func GetStartAndEndUTC4Today() (time.Time, time.Time) {
	start, now := GetStartAndEnd4Today()
	return start.UTC(), now.UTC()
}

/*
获取本月的开始和结束时间
*/
func GetStartAndEnd4CurrentMonth() (time.Time, time.Time) {
	// 获取当前时间
	now := time.Now()

	// 获取本月的开始时间
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	// 获取下个月的开始时间
	firstDayOfNextMonth := firstDayOfMonth.AddDate(0, 1, 0)

	// 获取本月的结束时间
	lastDayOfMonth := firstDayOfNextMonth.Add(-time.Nanosecond)
	return firstDayOfMonth, lastDayOfMonth
}

/*
获取本月的开始和结束时间-UTC
*/
func GetStartAndEndUTC4CurrentMonth() (time.Time, time.Time) {
	start, now := GetStartAndEnd4CurrentMonth()
	return start.UTC(), now.UTC()
}
