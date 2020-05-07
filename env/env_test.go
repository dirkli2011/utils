package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	testData := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	for key, value := range testData {
		Set(key, value)
		assert.True(t, Has(key))
		assert.Equal(t, Get(key), value)
	}

	assert.Equal(t, Get("k4", "v4"), "v4")

	all := GetAll()
	assert.True(t, len(all) > 0)
	for key, value := range testData {
		assert.Equal(t, all[key], value)
	}

}
