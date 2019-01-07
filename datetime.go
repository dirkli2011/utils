// 日期转换
package utils

import (
	"time"
)

func DateToTime(date string) int64 {
	format := "2006-01-02 15:04:05"
	nFormat := format[0:len(date)]

	loc, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation(nFormat, date, loc)
	if err != nil {
		return 0
	}

	return stime.Unix()
}

func TimeToDate(timestamp int64, f ...string) string {
	if timestamp <= 0 {
		return ""
	}
	tp := time.Unix(timestamp, 0)
	format := "2006-01-02"
	if len(f) > 0 {
		format = f[0]
	}
	return tp.Format(format)
}
