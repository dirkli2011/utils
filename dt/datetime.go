// 日期转换
package dt

import (
	"time"
)

// 指数退避算法，产生一个需求范围内的随机数，用于重试
func ExpBackOff(prev, max time.Duration) time.Duration {
	if prev == 0 {
		return time.Second
	}

	if prev > max/2 {
		return max
	}
	return 2 * prev
}

// 当前时间戳
func Now() int64 {
	return time.Now().Unix()
}

// 当前日期
func CurrentDate() string {
	return Date(Now())
}

// 当前日期时间
func CurrentDateTime() string {
	return DateTime(Now())
}

// 日期转时间戳
func Time(date string) int64 {
	format := "2006-01-02 15:04:05"
	nFormat := format[0:len(date)]

	loc, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation(nFormat, date, loc)
	if err != nil {
		return 0
	}

	return stime.Unix()
}

// 时间戳转日期
func Date(timestamp int64, f ...string) string {
	return FormatDate("2006-01-02", timestamp, f...)
}

// 时间戳转日期时间
func DateTime(timestamp int64, f ...string) string {
	return FormatDate("2006-01-02 15:04:05", timestamp, f...)
}

// 自定义格式输出日期字符串
func FormatDate(format string, timestamp int64, f ...string) string {
	if timestamp <= 0 {
		return ""
	}
	tp := time.Unix(timestamp, 0)
	if len(f) > 0 {
		format = f[0]
	}
	return tp.Format(format)
}
