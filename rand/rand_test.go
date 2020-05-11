package rand

import (
	"testing"
)

func TestRand(t *testing.T) {
	r := New()
	for i := 1; i < 10; i++ {
		t.Log(r.Base(i))
		t.Log(r.RandString(i))
	}
}
