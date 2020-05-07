package logkit

import (
	"strings"
	"testing"
)

var s = strings.Repeat("abcd", 100)

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info(s)
	}
}

func BenchmarkDateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Info(s)
		}
	})
}
