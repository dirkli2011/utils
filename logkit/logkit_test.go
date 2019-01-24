package logkit

import (
	"testing"
)

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("info ....")
	}
}

func BenchmarkDateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Info("info ....")
		}
	})
}
