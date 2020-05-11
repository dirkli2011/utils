package cryptolib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAes(t *testing.T) {
	key := strings.Repeat("1", 1)

	testData := []string{
		"1", "22", "333",
		"abc",
		"ABCDEFG",
	}

	for _, v := range testData {
		output, err := AESEncrypt([]byte(v), []byte(key))
		assert.Nil(t, err)
		actual, err := AESDecrypt(output, []byte(key))
		assert.Nil(t, err)
		assert.Equal(t, v, string(actual))
	}
}
