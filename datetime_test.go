package utils

import (
	"math/rand"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 单元测试
func TestDate(t *testing.T) {
	testData := []struct {
		input  string
		expect int64
	}{
		{"2019-01-01 00:00:00", 1546272000},
		{"2019-01-01 00:00:01", 1546272001},
		{"2019-01-01 00:01:40", 1546272100},
		{"2011-01-01 00:01:40", 1293811300},
	}

	for _, v := range testData {

		Convey("日期单元测试", t, func() {
			So(DateToTime(v.input), ShouldEqual, v.expect)
			So(TimeToDate(v.expect), ShouldEqual, strings.Split(v.input, " ")[0])
		})
	}

}

// 性能测试，b.N会根据函数的运行时间自动取一个合适的值
func BenchmarkDate(b *testing.B) {
	t := 1293811300
	for i := 0; i < b.N; i++ {
		TimeToDate(int64(t+rand.Intn(i+1)), "2006-01-02 15:04:05")
	}
}

// 多线程安全测试
func BenchmarkDateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			TimeToDate(1293811300)
		}
	})
}
