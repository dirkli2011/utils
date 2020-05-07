package cmap

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCMap(t *testing.T) {
	m := New()
	assert.NotNil(t, m)
	assert.Equal(t, m.Len(), 0)
	assert.True(t, m.Empty())

	testData := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	m.MultiSet(testData)
	assert.False(t, m.Empty())
	assert.Equal(t, m.Len(), len(testData))

	assert.True(t, m.Exist("k1"))
	v, ok := m.Get("k1")
	assert.True(t, ok)
	assert.Equal(t, v, "v1")

	m.Delete("k1")

	v, ok = m.Get("k1")
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok = m.Pop("k2")
	assert.True(t, ok)
	assert.Equal(t, v, "v2")

	assert.False(t, m.Exist("k2"))
	ok = m.SetIfNot("k2", "v2")
	assert.True(t, ok)
	assert.True(t, m.Exist("k2"))
}

func TestConcurrent(t *testing.T) {
	count := 1000
	goNums := 10
	m := New()
	var wg sync.WaitGroup

	wg.Add(goNums * 2)
	for i := 0; i < goNums; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				m.Set(strconv.Itoa(j), j)
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				m.SetIfNot(strconv.Itoa(j), j)
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, m.Len(), count)

	wg.Add(goNums)
	for i := 0; i < goNums; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				v, ok := m.Get(strconv.Itoa(j))
				assert.True(t, ok)
				assert.Equal(t, v, j)
			}
		}()
	}
	wg.Wait()

	total := 0
	for item := range m.Iterator() {
		total++
		v, ok := m.Get(item.Key)
		assert.True(t, ok)
		assert.Equal(t, v, item.Value)
	}
	assert.True(t, total > 0)
	assert.Equal(t, m.Len(), total)

	assert.Equal(t, len(m.Keys()), m.Len())
	assert.Equal(t, len(m.Items()), m.Len())

	fn := func(key string, value interface{}) bool {
		k, _ := strconv.Atoi(key)
		if k >= 100 {
			return false
		}
		return true
	}
	assert.Equal(t, len(m.FilterBy(fn)), 100)

	gt200 := 0
	f := func(key string, value interface{}) {
		k, _ := strconv.Atoi(key)
		if k >= 200 {
			gt200++
		}
	}
	m.IteratorBy(f)
	assert.Equal(t, gt200, count-200)
}

func BenchmarkCMap(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), i)
	}
}
