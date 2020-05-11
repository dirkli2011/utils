package rand

import (
	"math/rand"
	"time"
)

type Rander struct {
	rand *rand.Rand
}

func New() *Rander {
	return &Rander{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// return a int value
func (r *Rander) Base(n int) int {
	return r.rand.Intn(n)
}

// return a int value between <min, max>
func (r *Rander) Rand(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return r.Base(max-min) + min
}

// return a int slice
func (r *Rander) Batch(n int) []int {
	return r.rand.Perm(n)
}

// 返回一个随机字符串
func (r *Rander) RandString(n int, chars ...byte) string {
	if len(chars) == 0 {
		chars = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	}

	var bytes = make([]byte, n)
	for i, _ := range bytes {
		bytes[i] = chars[r.Base(len(chars))]
	}
	return string(bytes)
}
