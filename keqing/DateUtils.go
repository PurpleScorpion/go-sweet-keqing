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

func ParseDate(dateStr string, format string) time.Time {
	date, err := time.Parse(format, dateStr)
	if err != nil {
		panic("utc time parse error")
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
	return ParseDate(utc, DEFAULT_UTC_FORMAT)
}

func ParseLocal(local string) time.Time {
	return ParseDate(local, DEFAULT_LOCAL_FORMAT)
}

func UTC2Local(utc string) string {
	date := ParseDate(utc, DEFAULT_UTC_FORMAT)
	local := date.Local()
	return FormatDate(local, DEFAULT_LOCAL_FORMAT)
}

func UTC2LocalCustom(utc string, utcFormat string, localFormat string) string {
	date := ParseDate(utc, utcFormat)
	local := date.Local()
	return FormatDate(local, localFormat)
}

func Local2UTC(local string) string {
	date := ParseDate(local, DEFAULT_LOCAL_FORMAT)
	utc := date.UTC()
	return FormatDate(utc, DEFAULT_UTC_FORMAT)
}

func Local2UTCCustom(local string, localFormat string, utcFormat string) string {
	date := ParseDate(local, localFormat)
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
获取当前系统的毫秒级时间戳
*/
func CurrentTimeSeconds() int64 {
	// 获取当前时间的纳秒级时间戳
	return NowDate().Unix()
}
