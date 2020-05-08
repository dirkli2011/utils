package dt

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 单元测试
func TestDate(t *testing.T) {
	now := Now()
	date := CurrentDate()
	datetime := CurrentDateTime()

	testData := []struct {
		input  string
		expect int64
	}{
		{"2019-01-01 00:00:00", 1546272000},
		{"2019-01-01 00:00:01", 1546272001},
		{"2019-01-01 00:01:40", 1546272100},
		{"2011-01-01 00:01:40", 1293811300},
		{input: datetime, expect: now},
		{input: date + " 00:00:00", expect: Time(date)},
	}

	for _, v := range testData {
		assert.Equal(t, Time(v.input), v.expect)
		assert.Equal(t, DateTime(v.expect), v.input)
		assert.Equal(t, Date(v.expect), strings.Split(v.input, " ")[0])
	}

}

// 性能测试，b.N会根据函数的运行时间自动取一个合适的值
func BenchmarkDate(b *testing.B) {
	t := Now()
	for i := 0; i < b.N; i++ {
		Date(t + rand.Int63() + int64(i))
	}
}

// 多线程安全测试
func BenchmarkDateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Date(Now())
		}
	})
}
